package crawltask

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"go-to-crawl-video/internal/dao"
	"go-to-crawl-video/internal/model/entity"
	"go-to-crawl-video/internal/service/infra/configservice"
	"go-to-crawl-video/internal/service/uploadservice"
	"go-to-crawl-video/internal/service/videoservice"
	"go-to-crawl-video/utility/ffmpegutil"
	"path/filepath"
)

var (
	columns = dao.CrawlUploadQueue.Columns()
)

// 把任务队列里的视频资源转换成M3U8格式视频资源
func TransformTask(ctx gctx.Ctx) {
	//g.Dump("==============TransformTask================")
	log := g.Log().Line()
	//查找配置文件IP下正在转码的数据
	tans, err := dao.CrawlUploadQueue.Ctx(gctx.GetInitCtx()).Count(g.Map{
		columns.UploadStatus: uploadservice.Transforming,
		columns.HostLabel:    configservice.GetCrawlCfg("hostLabel"),
	})

	if err != nil {
		log.Infof(gctx.GetInitCtx(), "countErr:%v", err)
		return
	}
	if tans >= g.NewVar(configservice.GetCrawlCfg("maxTrans")).Int() {
		//	同时转码数量超过配置文件数量则不继续
		log.Infof(gctx.GetInitCtx(), "Trans Count Over Config.maxTrans")
		return
	}

	var queue *entity.CrawlUploadQueue
	//查找数据库上传完毕状态的数据进行处理
	dao.CrawlUploadQueue.Ctx(gctx.GetInitCtx()).Scan(&queue,
		g.Map{
			columns.UploadStatus: uploadservice.Uploaded,
			columns.HostLabel:    configservice.GetCrawlCfg("hostLabel"),
		})
	if err != nil {
		//log.Info("gettableErr:%v", err)
		return
	}
	if queue == nil {
		//log.Info("no trans task")
		return
	}
	queue.UploadStatus = uploadservice.Transforming
	_, err = dao.CrawlUploadQueue.Ctx(gctx.GetInitCtx()).Data(queue).Where(columns.Id, queue.Id).Update()
	if err != nil {
		g.Log().Infof(gctx.GetInitCtx(), "UploadStatusErr:%v,row:%v", err, queue)
		return
	}
	finalFileDir := videoservice.GetVideoDir(queue.CountryCode, queue.VideoYear, queue.VideoCollId, queue.VideoItemId)
	if !gfile.Exists(finalFileDir) {
		errMk := gfile.Mkdir(finalFileDir)
		if errMk != nil {
			// 创建目录是为了本地调试知道目录对应项目的位置
			log.Info(gctx.GetInitCtx(), "创建目录失败")
		}
	}

	//视频文件处理开始
	//finalFilePath 需要处理的文件
	finalFilePath := finalFileDir + queue.FileName
	// mp4file 转换成MP4后的文件
	mp4file := filepath.Join(finalFileDir, "segment.mp4")
	//ffmpeg对象 转码切片在里面完成
	ffm := ffmpegutil.FmpegTrans("ffmpeg")
	err = ffm.CheckFile(finalFileDir, finalFilePath, mp4file)
	if err != nil {
		//视频文件处理错误更新状态
		log.Infof(gctx.GetInitCtx(), "err:%v", err)
		queue.UploadStatus = uploadservice.TransformErr
		queue.Msg = err.Error()
		_, _ = dao.CrawlUploadQueue.Ctx(gctx.GetInitCtx()).Data(queue).Where(columns.Id, queue.Id).Update()
		return
	}
	//视频文件处理完毕

	//视频处理完毕通知CMS处理
	info, _ := g.NewVar(queue).MarshalJSON()
	//	通知CMS这部剧已经切片完成
	log.Infof(gctx.GetInitCtx(), "DoCallBack:%v->Param:%v", configservice.GetCrawlCfg("callback_url"), string(info))
	postRes := g.Client().PostContent(gctx.GetInitCtx(), configservice.GetCrawlCfg("callback_url"), g.Map{"info": info})
	//状态默认为转码完成
	queue.UploadStatus = uploadservice.Transformed
	if !g.NewVar(postRes).IsEmpty() {
		if !g.NewVar(g.NewVar(postRes).Map()["code"]).IsEmpty() {
			if g.NewVar(g.NewVar(postRes).Map()["code"]).Int() == 0 {
				//通知cms返回成功则更新状态为CmsPostSuccess（方便后续做新定时任务检测状态为upload.Transformed但是请求cms不成功的 重新发送请求）
				queue.UploadStatus = uploadservice.CmsPostSuccess
			}
		}
	}
	queue.UpdateTime = gtime.Now()
	_, _ = dao.CrawlUploadQueue.Ctx(gctx.GetInitCtx()).Data(queue).Where(columns.Id, queue.Id).Update()

}
