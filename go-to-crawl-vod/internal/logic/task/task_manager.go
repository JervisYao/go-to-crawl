package task

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"go-to-crawl-vod/internal/logic/task/crawl"
)

func StartAllTask() {
	log := g.Log().Line()

	taskMap := make(map[string]*taskBO)
	registryAllTask(taskMap)

	vodTaskNameListVar, _ := g.Cfg().Get(gctx.GetInitCtx(), "server.openVodTaskList")

	vodTaskNameList := vodTaskNameListVar.Slice()
	doStartAllTask(log, vodTaskNameList, taskMap)
}

func registryAllTask(taskMap map[string]*taskBO) {
	// QQ
	registryTask(taskMap, "@every 1h", "checkQQLoginTask", crawl.CrawlCheckTask.CheckQQLoginTask)
	registryTask(taskMap, "@every 5s", "crawlUrlType1Task", crawl.CrawlTask.CrawlUrlType1Task)
	registryTask(taskMap, "@every 5s", "crawlUrlType2Task", crawl.CrawlTask.CrawlUrlType2Task)
	registryTask(taskMap, "@every 5s", "crawlUrlType3Task", crawl.CrawlTask.CrawlUrlType3Task)
	registryTask(taskMap, "@every 20s", "downloadMp4Type1Task", crawl.DownloadMp4Type1Task)
	registryTask(taskMap, "@every 20s", "downloadMp4Type2Task", crawl.DownloadMp4Type2Task)
	registryTask(taskMap, "@every 20s", "downloadMp4Type3Task", crawl.DownloadMp4Type3Task)

	// 【点播】定时生成视频列表抓取任务实例，根据实例去真正启动视频列表抓取任务
	registryTask(taskMap, "@every 10s", "genVodConfigTask", crawl.CrawlVodTVTask.GenVodConfigTask)
	// 【点播】根据点播配置和支持的策略自动获取所有视频节目单
	registryTask(taskMap, "@every 10s", "vodTVTask", crawl.CrawlVodTVTask.VodTVTask)
	// 【点播】填充剧相关信息
	registryTask(taskMap, "@every 5s", "vodTVPadInfoTask", crawl.CrawlVodTVTask.VodTVPadInfoTask)
	// 【点播】填充剧ID
	registryTask(taskMap, "@every 5s", "vodTVPadIdTask", crawl.CrawlVodTVTask.VodTVPadIdTask)
	// 【点播】填充集数ID
	registryTask(taskMap, "@every 5s", "vodTVItemPadIdTask", crawl.CrawlVodTVTask.VodTVItemPadIdTask)

	// 点播默认爬虫处理
	registryTask(taskMap, "@every 5s", "crawlUrlTask", crawl.CrawlTask.CrawlUrlTask)
	// 点播默认下载处理
	registryTask(taskMap, "@every 20s", "downloadMp4Task", crawl.DownloadMp4Task)
	// 点播默认转码处理
	registryTask(taskMap, "@every 10s", "transformTask", crawl.TransformTask)

	// 重置hosttype=2 解析失败（2）和下载失败（5）的记录
	registryTask(taskMap, "@every 1h", "resetHostType2Task", crawlService.ResetHostType2)

	// 重置等待下载（3）状态的记录
	registryTask(taskMap, "@every 10s", "resetProcessingTask", crawlService.ResetProcessingTooLong)
	// 重置正在爬取 (1)状态的记录
	registryTask(taskMap, "@every 10s", "resetCrawlingTask", crawlService.ResetCrawlingTooLong)
	// 查找转码完成没有post到cms成功的数据重新post
	registryTask(taskMap, "@every 10s", "cheeckPostcms", crawl.PostCmsTask)
	// 重置重试次数cnt>=3的所有记录
	registryTask(taskMap, "@every 10m", "crawlUrlFailNotifyTask", crawl.CrawlCheckTask.CrawlUrlFailNotifyTask)
}

func doStartAllTask(log *glog.Logger, taskNameList []string, taskMap map[string]*taskBO) {
	//g.Dump(taskNameList)
	if taskNameList == nil {
		return
	}

	for _, taskName := range taskNameList {
		taskItem := taskMap[taskName]
		//g.Dump(taskMap)
		if taskItem == nil {
			continue
		}

		var task *gcron.Entry
		if taskItem.once {
			task, _ = gcron.AddOnce(gctx.GetInitCtx(), taskItem.pattern, taskItem.job, taskItem.name)
		} else {
			task, _ = gcron.Add(gctx.GetInitCtx(), taskItem.pattern, taskItem.job, taskItem.name)
		}
		log.Info(gctx.GetInitCtx(), "新增Task: ", task.Name)
		gcron.Start(taskItem.name)
	}
}

func registryTask(taskMap map[string]*taskBO, pattern string, taskName string, job func(ctx gctx.Ctx)) {
	taskMap[taskName] = getTask(taskName, pattern, job)
}

func registryOnceTask(taskMap map[string]*taskBO, pattern string, taskName string, job func(ctx gctx.Ctx)) {
	t := getTask(taskName, pattern, job)
	t.once = true
	taskMap[taskName] = t
}

func getTask(taskName string, pattern string, job func(ctx gctx.Ctx)) *taskBO {
	t := new(taskBO)
	t.name = taskName
	t.pattern = pattern
	t.job = job
	return t
}

type taskBO struct {
	pattern string
	name    string
	once    bool
	job     func(ctx gctx.Ctx)
}