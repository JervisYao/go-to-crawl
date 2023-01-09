package crawl

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"go-to-crawl-vod/internal/dao"
	"go-to-crawl-vod/internal/service/infra/config"
	"go-to-crawl-vod/utility/ffmpeg"
	"path/filepath"
)

var (
	columns = dao.UploadQueue.Columns
)

// 把任务队列里的视频资源转换成M3U8格式视频资源
func TransformTask(ctx gctx.Ctx) {
	//g.Dump("==============TransformTask================")
	log := g.Log().Line()
	//查找配置文件IP下正在转码的数据
	tans, err := dao.UploadQueue.Count(g.Map{
		columns.UploadStatus: upload.Transforming,
		columns.HostIp:       config.GetCrawlCfg("hostIp"),
	})
	//g.Dump(tans)
	if err != nil {
		log.Infof("countErr:%v", err)
		return
	}
	if tans >= g.NewVar(config.GetCrawlCfg("maxTrans")).Int() {
		//	同时转码数量超过配置文件数量则不继续
		log.Infof("Trans Count Over Config.maxTrans")
		return
	}

	//查找数据库上传完毕状态的数据进行处理
	queue, err := dao.UploadQueue.FindOne(
		g.Map{
			columns.UploadStatus: upload.Uploaded,
			columns.HostIp:       config.GetCrawlCfg("hostIp"),
		})
	if err != nil {
		//log.Info("gettableErr:%v", err)
		return
	}
	if queue == nil {
		//log.Info("no trans task")
		return
	}
	queue.UploadStatus = upload.Transforming
	_, err = dao.UploadQueue.Data(queue).Where(columns.Id, queue.Id).Update()
	if err != nil {
		g.Log().Infof(gctx.GetInitCtx(), "UploadStatusErr:%v,row:%v", err, queue)
		return
	}
	finalFileDir := video.GetVideoDir(queue.CountryCode, queue.VideoYear, queue.VideoCollId, queue.VideoItemId)
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
	ffm := ffmpeg.FmpegTrans("ffmpeg")
	err = ffm.CheckFile(finalFileDir, finalFilePath, mp4file)
	if err != nil {
		//视频文件处理错误更新状态
		log.Infof(gctx.GetInitCtx(), "err:%v", err)
		queue.UploadStatus = upload.TransformErr
		queue.Msg = err.Error()
		_, _ = dao.UploadQueue.Data(queue).Where(columns.Id, queue.Id).Update()
		return
	}
	//视频文件处理完毕

	//视频处理完毕通知CMS处理
	info, _ := g.NewVar(queue).MarshalJSON()
	//	通知CMS这部剧已经切片完成
	log.Infof(gctx.GetInitCtx(), "DoCallBack:%v->Param:%v", config.GetCrawlCfg("callback_url"), string(info))
	postRes := g.Client().PostContent(config.GetCrawlCfg("callback_url"), g.Map{"info": info})
	//状态默认为转码完成
	queue.UploadStatus = upload.Transformed
	if !g.NewVar(postRes).IsEmpty() {
		if !g.NewVar(g.NewVar(postRes).Map()["code"]).IsEmpty() {
			if g.NewVar(g.NewVar(postRes).Map()["code"]).Int() == 0 {
				//通知cms返回成功则更新状态为CmsPostSuccess（方便后续做新定时任务检测状态为upload.Transformed但是请求cms不成功的 重新发送请求）
				queue.UploadStatus = upload.CmsPostSuccess
			}
		}
	}
	queue.UpdateTime = gtime.Now()
	_, _ = dao.UploadQueue.Data(queue).Where(columns.Id, queue.Id).Update()

}