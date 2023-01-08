package crawl

import (
	"go-to-crawl-vod/internal/model/do"
	"testing"
)

func TestBananTV(t *testing.T) {
	doStartVodTV(Banan, "https://banan.tv/vodtype/22-1.html")
}

// 偏单元测试
func TestNiVodTV(t *testing.T) {
	doStartVodTV(NiVod, "https://www.mudvod.tv/filter.html?x=1&channelId=7&showTypeId=145")
}

func TestNiVodTVPadInfo(t *testing.T) {
	doStartVodPadInfo("https://www.mudvod.tv/l01tKF0bT0wpcHiJVV2CsJpLH3gDDP8T-0-0-0-0-detail.html?x=1")
}

func doStartVodTV(domain, url string) {
	seed := new(do.CrawlVodConfigTaskDO)
	seed.CmsCrawlVodConfig = new(model.CmsCrawlVodConfig)
	seed.SeedUrl = url
	seed.PageSize = 6
	seed.DomainKeyPart = domain
	DoStartCrawlVodTV(seed)
}

func doStartVodPadInfo(seedUrl string) {
	vodTvItem := new(model.CmsCrawlVodTv)
	vodTvItem.SeedUrl = seedUrl
	DoStartCrawlVodPadInfo(vodTvItem)
}

func TestGenVodConfigTask(t *testing.T) {
	CrawlVodTVTask.GenVodConfigTask()
}

func TestVodTVTask(t *testing.T) {
	CrawlVodTVTask.VodTVTask()
}

func TestVodTVPadInfoTask(t *testing.T) {
	CrawlVodTVTask.VodTVPadInfoTask()
}

func TestVodTVPadIdTask(t *testing.T) {
	CrawlVodTVTask.VodTVPadIdTask()
}

func TestVodTVItemPadIdTask(t *testing.T) {
	CrawlVodTVTask.VodTVItemPadIdTask()
}
