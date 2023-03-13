package crawlservice

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"go-to-crawl-common/consts"
	"go-to-crawl-common/service/infra/configservice"
	"go-to-crawl-video/internal/dao"
	"go-to-crawl-video/internal/model/entity"
	"time"
)

var (
	c = dao.CrawlQueue.Columns()
)

// hostLabel: 配置文件里能标识出当前节点唯一就行，eg：prod-1, dev-1, cluster-1等
func GetSeed(status int, hostLabel string, businessType int) *entity.CrawlQueue {
	where := dao.CrawlQueue.Ctx(gctx.GetInitCtx()).Where(c.CrawlStatus, status).Where(c.BusinessType, businessType)
	if hostLabel != "" {
		where = where.Where(c.HostLabel, hostLabel)
	}

	var seed *entity.CrawlQueue
	_ = where.Scan(&seed)
	return seed
}

func GetNeedNotifySeedList() []*entity.CrawlQueue {
	where := dao.CrawlQueue.Ctx(gctx.GetInitCtx()).Where(c.CrawlM3U8Notify, CrawlM3U8NotifyNo).WhereGTE(c.CrawlM3U8Url, consts.ServerMaxRetry)

	var array []*entity.CrawlQueue
	where.ScanList(&array, "CrawlQueue")

	return array
}

func ExistCrawling(hostType int) bool {
	//爬虫状态为爬取中的数量大于0时 不继续
	//需要新增条件限制host_type不同类型限制 避免互相冲突
	count, _ := dao.CrawlQueue.Ctx(gctx.GetInitCtx()).Where(c.BusinessType, hostType).Count(c.CrawlStatus, Crawling)
	return count > 0
}

func UpdateStatus(seed *entity.CrawlQueue, status int) {
	if configservice.GetCrawlDebugBool("disableDB") {
		return
	}
	seed.CrawlStatus = status
	seed.UpdateTime = gtime.Now()
	dao.CrawlQueue.Ctx(gctx.GetInitCtx()).Data(seed).Where(c.Id, seed.Id).Update()
}

func UpdateUrlAndStatus(seed *entity.CrawlQueue) {
	if configservice.GetCrawlBool("disableDB") {
		return
	}
	seed.CrawlM3U8Cnt = seed.CrawlM3U8Cnt + 1
	seed.UpdateTime = gtime.Now()
	if seed.CrawlM3U8Url == "" && seed.CrawlM3U8Text == "" {
		if seed.CrawlM3U8Cnt >= consts.ServerMaxRetry {
			// 超过允许重试的最大次数
			seed.HostLabel = configservice.GetCrawlHostLabel()
			if seed.ErrorMsg == "" {
				seed.ErrorMsg = "M3U8 Empty"
			}
			UpdateStatus(seed, CrawlErr)
		} else {
			UpdateStatus(seed, Init)
		}
	} else {
		UpdateStatus(seed, CrawlFinish)
	}
}

// 重置处理中状态太久的
func ResetProcessingTooLong(ctx gctx.Ctx) {
	ResetHangingStatus(M3U8Parsing, CrawlFinish, 6*60)
}

// 重置抓取中状态太久的
func ResetCrawlingTooLong(ctx gctx.Ctx) {
	ResetHangingStatus(Crawling, Init, 1)
}

// 重置挂起中的状态
func ResetHangingStatus(fromStatus, toStatus, hangingMinutes int) {
	waterMark := gtime.Now().Add(time.Duration(hangingMinutes) * -time.Minute)

	var seed *entity.CrawlQueue
	dao.CrawlQueue.Ctx(gctx.GetInitCtx()).WhereLT(c.UpdateTime, waterMark).Scan(&seed, c.CrawlStatus, fromStatus)
	if seed == nil {
		return
	}

	UpdateStatus(seed, toStatus)
}

func ResetHostType2(ctx gctx.Ctx) {
	//waterMark := gtime.Now().Add(time.Duration(hangingMinutes) * -time.Minute)
	dao.CrawlQueue.Ctx(gctx.GetInitCtx()).Where(g.Map{
		"host_type = ?":                       2,
		"(crawl_status =2 or crawl_status=5)": "",
	}).Update(g.Map{
		"crawl_status": 0,
	})

	//UpdateStatus(seed, toStatus)
}
