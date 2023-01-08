package crawl

import (
	"github.com/gogf/gf/v2/text/gstr"
	"go-to-crawl-vod/internal/logic/task/crawl/bilibili"
	"go-to-crawl-vod/internal/logic/task/crawl/nunuyy"
	"go-to-crawl-vod/internal/logic/task/dto"
)

const (
	Nunuyy     = "nunuyy"
	Bilibbili  = "bilibili"
	Ole        = "ole"
	TangRenJie = "tangrenjie.tv"
	QQ         = "v.qq.com"
	NiVod      = "nivod.tv"
	MudVod     = "mudvod.tv"
	Banan      = "banan.tv"
	Iqiyi      = "iqiyi.com"
)

func getCrawlVodFlowStrategy(seed *model.CmsCrawlQueue) dto.CrawlVodFlowInterface {

	url := seed.CrawlSeedUrl

	if gstr.Contains(url, Nunuyy) {
		return new(nunuyy.NunuyyCrawl)
	} else if gstr.Contains(url, Bilibbili) {
		return new(bilibili.BilibiliCrawl)
	} else if gstr.Contains(url, Ole) {
		return new(olevod.OleVodCrawl)
	} else if gstr.Contains(url, TangRenJie) {
		return new(tangrenjie.TangRenJieCrawl)
	} else if gstr.Contains(url, QQ) {
		return new(qq.QQCrawl)
	} else if gstr.Contains(url, NiVod) || gstr.Contains(url, MudVod) {
		return new(nivod.NiVodCrawl)
	} else if gstr.Contains(url, Banan) {
		return new(banan.BananCrawl)
	} else if gstr.Contains(url, Iqiyi) {
		return new(iqiyi.IqiyiCrawl)
	}

	return nil
}

func getCrawlVodTVStrategy(seed *model.CmsCrawlVodConfig) dto.CrawlVodTVInterface {

	url := seed.SeedUrl

	if gstr.Contains(url, NiVod) || gstr.ContainsI(url, MudVod) {
		return new(nivod.NiVodTVTask)
	} else if gstr.Contains(url, Banan) {
		return new(banan.BananTvCrawl)
	}

	return nil
}

func getCrawlVodPadInfoStrategy(seed *model.CmsCrawlVodTv) dto.CrawlVodTVInterface {

	url := seed.SeedUrl
	if gstr.Contains(url, NiVod) || gstr.ContainsI(url, MudVod) {
		return new(nivod.NiVodTVTask)
	} else if gstr.Contains(url, Banan) {
		return new(banan.BananTvCrawl)
	}

	return nil
}

func GetHostType(crawlSeedUrl string) int {
	if gstr.Contains(crawlSeedUrl, QQ) {
		return crawl.HostTypeCrawlLogin
	} else {
		return crawl.HostTypeNormal
	}
}
