package crawl

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"go-to-crawl-vod/internal/dao"
	"go-to-crawl-vod/internal/service/infra/config"
	"time"
)

var (
	c = dao.CrawlQueue.Columns()
)

// hostName: 配置文件里能标识出当前节点唯一就行，eg：prod-1, dev-1, cluster-1等
func GetSeed(status int, hostname string, hostType int) *model.CmsCrawlQueue {
	where := dao.CrawlQueue.Where(c.CrawlStatus, status).And(c.HostType, hostType)
	if hostname != "" {
		where = where.And(c.HostIp, hostname)
	}
	seed, _ := where.FindOne()
	return seed
}

func GetNeedNotifySeedList() []*model.CmsCrawlQueue {
	where := dao.CrawlQueue.Where(c.CrawlM3U8Notify, CrawlM3U8NotifyNo).And("crawl_m3u8_cnt >= ?", constant.ServerMaxRetry)
	all, _ := where.FindAll()
	return all
}

func ExistCrawling(hostType int) bool {
	//爬虫状态为爬取中的数量大于0时 不继续
	//需要新增条件限制host_type不同类型限制 避免互相冲突
	count, _ := dao.CrawlQueue.Where(c.HostType, hostType).Count(c.CrawlStatus, Crawling)
	return count > 0
}

func UpdateStatus(seed *model.CmsCrawlQueue, status int) {
	if config.GetCrawlDebugBool("disableDB") {
		return
	}
	seed.CrawlStatus = status
	seed.UpdateTime = gtime.Now()
	dao.CrawlQueue.Data(seed).Where(c.Id, seed.Id).Update()
}

func UpdateUrlAndStatus(seed *model.CmsCrawlQueue) {
	if config.GetCrawlBool("disableDB") {
		return
	}
	seed.CrawlM3U8Cnt = seed.CrawlM3U8Cnt + 1
	seed.UpdateTime = gtime.Now()
	if seed.CrawlM3U8Url == "" && seed.CrawlM3U8Text == "" {
		if seed.CrawlM3U8Cnt >= constant.ServerMaxRetry {
			// 超过允许重试的最大次数
			seed.HostIp = config.GetCrawlHostIp()
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
func ResetProcessingTooLong() {
	ResetHangingStatus(M3U8Parsing, CrawlFinish, 6*60)
}

// 重置抓取中状态太久的
func ResetCrawlingTooLong() {
	ResetHangingStatus(Crawling, Init, 1)
}

// 重置挂起中的状态
func ResetHangingStatus(fromStatus, toStatus, hangingMinutes int) {
	waterMark := gtime.Now().Add(time.Duration(hangingMinutes) * -time.Minute)
	seed, _ := dao.CrawlQueue.Where("update_time < ", waterMark).FindOne(c.CrawlStatus, fromStatus)
	if seed == nil {
		return
	}

	UpdateStatus(seed, toStatus)
}

func ResetHostType2() {
	//waterMark := gtime.Now().Add(time.Duration(hangingMinutes) * -time.Minute)
	dao.CrawlQueue.Where(g.Map{
		"host_type = ?":                       2,
		"(crawl_status =2 or crawl_status=5)": "",
	}).Update(g.Map{
		"crawl_status": 0,
	})

	//UpdateStatus(seed, toStatus)
}
