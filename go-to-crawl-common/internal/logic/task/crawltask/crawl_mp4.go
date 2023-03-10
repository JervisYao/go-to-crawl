package crawltask

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"go-to-crawl-common/utility/ffmpegutil"
	"go-to-crawl-common/utility/fileutil"
	crawlservice2 "go-to-crawl-video/internal/service/crawlservice"
	"go-to-crawl-video/internal/service/infra/configservice"
	lockservice2 "go-to-crawl-video/internal/service/infra/lockservice"
	"go-to-crawl-video/internal/service/videoservice"
)

func DownloadMp4Type1Task(ctx gctx.Ctx) {
	doDownloadMp4(crawlservice2.BusinessTypeCrawlLogin)
}

func DownloadMp4Type2Task(ctx gctx.Ctx) {
	doDownloadMp4(crawlservice2.BusinessTypeNiVod)
}

func DownloadMp4Type3Task(ctx gctx.Ctx) {
	if lockservice2.IncreaseValue(lockservice2.DownloadMp4Type3) {
		defer lockservice2.DecreaseValue(lockservice2.DownloadMp4Type3)
		doDownloadMp4(crawlservice2.BusinessTypeBananTV)
	}
}

func DownloadMp4Task(ctx gctx.Ctx) {
	doDownloadMp4(crawlservice2.BusinessTypeNormal)
}

func doDownloadMp4(hostType int) {
	seed := crawlservice2.GetSeed(crawlservice2.CrawlFinish, configservice.GetCrawlHostLabel(), hostType)

	if seed == nil {
		return
	}
	log := g.Log().Line()

	crawlservice2.UpdateStatus(seed, crawlservice2.M3U8Parsing)
	// 创建最终目录
	videoDir := videoservice.GetVideoDir(seed.CountryCode, seed.VideoYear, seed.VideoCollId, seed.VideoItemId)
	_ = gfile.Mkdir(videoDir)

	// 下载M3U8文件
	orgM3U8File := videoDir + ffmpegutil.OrgM3U8Name
	proxyUrl := crawlservice2.GetProxyByUrl(seed.CrawlM3U8Url)

	err := fileutil.DownloadM3U8File(seed.CrawlM3U8Url, proxyUrl, orgM3U8File, fileutil.Retry, seed.CrawlM3U8Text)
	if err != nil {
		log.Info(gctx.GetInitCtx(), err)
		seed.ErrorMsg = "Download M3U8 Error"
		crawlservice2.UpdateStatus(seed, crawlservice2.M3U8Err)
		return
	}
	// 下载完M3U8后，后续操作都只能当前主机处理
	seed.HostLabel = configservice.GetCrawlHostLabel()

	if crawlservice2.TypeMP4Url == seed.CrawlType {
		// 直接下载MP4
		builder := fileutil.CreateBuilder()
		builder.Url(seed.CrawlSeedUrl)
		builder.SaveFile(fmt.Sprintf("%s%s", videoDir, ffmpegutil.OrgMp4Name))
		err2 := fileutil.DownloadFileByBuilder(builder)
		if err2 != nil {
			videoservice.UpdateDownloadStatus(seed, errors.New("MP4下载失败"))
			return
		} else {
			videoservice.UpdateDownloadStatus(seed, nil)
		}
	} else {
		strategy := getCrawlVodFlowStrategy(seed)
		m3u8DO, err2 := strategy.ConvertM3U8(seed, orgM3U8File)
		if err2 != nil {
			log.Info(gctx.GetInitCtx(), err2)
			seed.ErrorMsg = "标准化M3U8文件出错"
			crawlservice2.UpdateStatus(seed, crawlservice2.M3U8Err)
			return
		}

		err2 = strategy.DownLoadToMp4(m3u8DO)

		if err2 != nil {
			videoservice.UpdateDownloadStatus(seed, errors.New("M3U8转MP4出错"))
			return
		} else {
			videoservice.UpdateDownloadStatus(seed, nil)
		}
		//更新成功後刪除原m3u8文件
		gfile.Remove(orgM3U8File)

		if crawlservice2.BusinessTypeCrawlLogin == hostType {
			// 国内指定机器下载的，需要上传到国外点播服务器
			//file.UpLoadToFastDFS(m3u8DO.MP4File, seed)
		}
	}

	// 添加到转换队列
	// TODO
	/*upQueue := new(model.CmsUploadQueue)
	gconv.Struct(seed, upQueue)
	upQueue.Id = 0
	upQueue.FileName = ffmpegutil.OrgMp4Name
	upQueue.UploadStatus = upload.Uploaded
	upQueue.CreateTime = gtime.Now()
	dao.UploadQueue.Insert(upQueue)*/

}
