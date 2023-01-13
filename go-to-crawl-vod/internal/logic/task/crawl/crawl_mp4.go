package crawl

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"go-to-crawl-vod/internal/dao"
	"go-to-crawl-vod/internal/service/crawl"
	"go-to-crawl-vod/internal/service/infra/config"
	"go-to-crawl-vod/internal/service/infra/lock"
	"go-to-crawl-vod/utility/ffmpegutil"
	"go-to-crawl-vod/utility/fileutil"
)

func DownloadMp4Type1Task(ctx gctx.Ctx) {
	doDownloadMp4(crawl.BusinessTypeCrawlLogin)
}

func DownloadMp4Type2Task(ctx gctx.Ctx) {
	doDownloadMp4(crawl.BusinessTypeNiVod)
}

func DownloadMp4Type3Task(ctx gctx.Ctx) {
	if lock.IncreaseValue(lock.DownloadMp4Type3) {
		defer lock.DecreaseValue(lock.DownloadMp4Type3)
		doDownloadMp4(crawl.BusinessTypeBananTV)
	}
}

func DownloadMp4Task(ctx gctx.Ctx) {
	doDownloadMp4(crawl.BusinessTypeNormal)
}

func doDownloadMp4(hostType int) {
	seed := crawl.GetSeed(crawl.CrawlFinish, config.GetCrawlHostIp(), hostType)

	if seed == nil {
		return
	}
	log := g.Log().Line()

	crawl.UpdateStatus(seed, crawl.M3U8Parsing)
	// 创建最终目录
	videoDir := video.GetVideoDir(seed.CountryCode, seed.VideoYear, seed.VideoCollId, seed.VideoItemId)
	_ = gfile.Mkdir(videoDir)

	// 下载M3U8文件
	orgM3U8File := videoDir + ffmpegutil.OrgM3U8Name
	proxyUrl := crawl.GetProxyByUrl(seed.CrawlM3U8Url)

	err := fileutil.DownloadM3U8File(seed.CrawlM3U8Url, proxyUrl, orgM3U8File, fileutil.Retry, seed.CrawlM3U8Text)
	if err != nil {
		log.Info(err)
		seed.ErrorMsg = "Download M3U8 Error"
		crawl.UpdateStatus(seed, crawl.M3U8Err)
		return
	}
	// 下载完M3U8后，后续操作都只能当前主机处理
	seed.HostIp = config.GetCrawlHostIp()

	if crawl.TypeMP4Url == seed.CrawlType {
		// 直接下载MP4
		builder := fileutil.CreateBuilder()
		builder.Url(seed.CrawlSeedUrl)
		builder.SaveFile(fmt.Sprintf("%s%s", videoDir, ffmpegutil.OrgMp4Name))
		err2 := fileutil.DownloadFileByBuilder(builder)
		if err2 != nil {
			video.UpdateDownloadStatus(seed, errors.New("MP4下载失败"))
			return
		} else {
			video.UpdateDownloadStatus(seed, nil)
		}
	} else {
		strategy := getCrawlVodFlowStrategy(seed)
		m3u8DO, err2 := strategy.ConvertM3U8(seed, orgM3U8File)
		if err2 != nil {
			log.Info(err2)
			seed.ErrorMsg = "标准化M3U8文件出错"
			crawl.UpdateStatus(seed, crawl.M3U8Err)
			return
		}

		err2 = strategy.DownLoadToMp4(m3u8DO)

		if err2 != nil {
			video.UpdateDownloadStatus(seed, errors.New("M3U8转MP4出错"))
			return
		} else {
			video.UpdateDownloadStatus(seed, nil)
		}
		//更新成功後刪除原m3u8文件
		gfile.Remove(orgM3U8File)

		if crawl.BusinessTypeCrawlLogin == hostType {
			// 国内指定机器下载的，需要上传到国外点播服务器
			//file.UpLoadToFastDFS(m3u8DO.MP4File, seed)
		}
	}

	// 添加到转换队列
	upQueue := new(model.CmsUploadQueue)
	gconv.Struct(seed, upQueue)
	upQueue.Id = 0
	upQueue.FileName = ffmpegutil.OrgMp4Name
	upQueue.UploadStatus = upload.Uploaded
	upQueue.CreateTime = gtime.Now()
	dao.UploadQueue.Insert(upQueue)

}
