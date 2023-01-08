package crawl

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"go-to-crawl-vod/internal/dao"
	"go-to-crawl-vod/internal/service/infra/config"
)

// 把任务队列里的视频资源转换成M3U8格式视频资源
func PostCmsTask(ctx gctx.Ctx) {
	log := g.Log().Line()

	//os.Exit(1)
	//查找数据库上传完毕状态的数据进行处理
	queue, err := dao.CmsUploadQueue.FindOne(
		g.Map{
			columns.UploadStatus: upload.Transformed,
			columns.HostIp:       config.GetCrawlCfg("hostIp"),
		})
	if err != nil {
		//log.Infof("PostCmsTask Not data")
		//log.Info("gettableErr:%v", err)
		return
	}
	if queue == nil {
		//log.Info("no trans task")
		return
	}
	//视频处理完毕通知CMS处理
	info, _ := g.NewVar(queue).MarshalJSON()

	log.Infof(gctx.GetInitCtx(), "DoCallBack:%v->Param:%v", config.GetCrawlCfg("callback_url"), string(info))
	postRes := g.Client().PostContent(config.GetCrawlCfg("callback_url"), g.Map{"info": info})
	if !g.NewVar(postRes).IsEmpty() {
		if !g.NewVar(g.NewVar(postRes).Map()["code"]).IsEmpty() {
			if g.NewVar(g.NewVar(postRes).Map()["code"]).Int() == 0 {
				//通知cms返回成功则更新状态为CmsPostSuccess（方便后续做新定时任务检测状态为upload.Transformed但是请求cms不成功的 重新发送请求）
				queue.UploadStatus = upload.CmsPostSuccess
				queue.UpdateTime = gtime.Now()
				_, _ = dao.CmsUploadQueue.Data(queue).Where(columns.Id, queue.Id).Update()
			}
		}
	}

}
