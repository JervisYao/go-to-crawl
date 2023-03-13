package crawltask

import (
	"github.com/gogf/gf/v2/text/gstr"
	"go-to-crawl-common/internal/logic/task/crawltask/bilibilitask"
	"go-to-crawl-common/internal/logic/task/crawltask/nunuyytask"
	"go-to-crawl-common/internal/logic/task/taskdto"
	entity2 "go-to-crawl-video/internal/model/entity"
	"go-to-crawl-video/internal/service/crawlservice"
)

const (
	Nunuyy    = "nunuyy"
	Bilibbili = "bilibili"
	QQ        = "v.qq.com"
)

func getCrawlVodFlowStrategy(seed *entity2.CrawlQueue) taskdto.CrawlVodFlowInterface {

	url := seed.CrawlSeedUrl

	if gstr.Contains(url, Nunuyy) {
		return new(nunuyytask.NunuyyCrawl)
	} else if gstr.Contains(url, Bilibbili) {
		return new(bilibilitask.BilibiliCrawl)
	}

	return nil
}

func getCrawlVodTVStrategy(seed *entity2.CrawlVodConfig) taskdto.CrawlVodTVInterface {
	return nil
}

func getCrawlVodPadInfoStrategy(seed *entity2.CrawlVod) taskdto.CrawlVodTVInterface {
	return nil
}

func GetHostType(crawlSeedUrl string) int {
	if gstr.Contains(crawlSeedUrl, QQ) {
		return crawlservice.BusinessTypeCrawlLogin
	} else {
		return crawlservice.BusinessTypeNormal
	}
}
