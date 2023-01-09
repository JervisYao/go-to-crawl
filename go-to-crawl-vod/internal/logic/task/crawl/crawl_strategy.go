package crawl

import (
	"github.com/gogf/gf/v2/text/gstr"
	"go-to-crawl-vod/internal/logic/task/crawl/bilibili"
	"go-to-crawl-vod/internal/logic/task/crawl/nunuyy"
	"go-to-crawl-vod/internal/logic/task/dto"
	"go-to-crawl-vod/internal/model/entity"
	"go-to-crawl-vod/internal/service/crawl"
)

const (
	Nunuyy    = "nunuyy"
	Bilibbili = "bilibili"
	QQ        = "v.qq.com"
)

func getCrawlVodFlowStrategy(seed *entity.CrawlQueue) dto.CrawlVodFlowInterface {

	url := seed.CrawlSeedUrl

	if gstr.Contains(url, Nunuyy) {
		return new(nunuyy.NunuyyCrawl)
	} else if gstr.Contains(url, Bilibbili) {
		return new(bilibili.BilibiliCrawl)
	}

	return nil
}

func getCrawlVodTVStrategy(seed *entity.CrawlVodConfig) dto.CrawlVodTVInterface {
	return nil
}

func getCrawlVodPadInfoStrategy(seed *entity.CrawlVod) dto.CrawlVodTVInterface {
	return nil
}

func GetHostType(crawlSeedUrl string) int {
	if gstr.Contains(crawlSeedUrl, QQ) {
		return crawl.BusinessTypeCrawlLogin
	} else {
		return crawl.BusinessTypeNormal
	}
}
