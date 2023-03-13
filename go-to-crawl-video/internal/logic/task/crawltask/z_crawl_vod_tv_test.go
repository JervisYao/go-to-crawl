package crawltask

import (
	"go-to-crawl-video/internal/model/entity"
	"testing"
)

func doStartVodTV(domain, url string) {
	seed := new(domain.CrawlVodConfigTaskDO)
	seed.CrawlVodConfig = new(entity.CrawlVodConfig)
	seed.SeedUrl = url
	seed.PageSize = 6
	seed.DomainKeyPart = domain
	DoStartCrawlVodTV(seed)
}

func doStartVodPadInfo(seedUrl string) {
	vodTvItem := new(entity.CrawlVod)
	vodTvItem.SeedUrl = seedUrl
	DoStartCrawlVodPadInfo(vodTvItem)
}

func TestGenVodConfigTask(t *testing.T) {
	CrawlVodTVTask.GenVodConfigTask(nil)
}

func TestVodTVTask(t *testing.T) {
	CrawlVodTVTask.VodTVTask(nil)
}

func TestVodTVPadInfoTask(t *testing.T) {
	CrawlVodTVTask.VodTVPadInfoTask(nil)
}

func TestVodTVPadIdTask(t *testing.T) {
	CrawlVodTVTask.VodTVPadIdTask(nil)
}

func TestVodTVItemPadIdTask(t *testing.T) {
	CrawlVodTVTask.VodTVItemPadIdTask(nil)
}
