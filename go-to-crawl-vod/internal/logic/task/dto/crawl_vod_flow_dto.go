package dto

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/tebeka/selenium"
	"go-to-crawl-vod/utility/ffmpeg"
	"go-to-crawl-vod/utility/http"
)

type BrowserContext struct {
	Log             *glog.Logger
	CrawlQueueSeed  *model.CmsCrawlQueue
	VodConfigTaskDO *do.CrawlVodConfigTaskDO
	VodTV           *model.CmsCrawlVodTv
	VodTVItem       *model.CmsCrawlVodTvItem
	Service         *selenium.Service
	XServer         *proxyServer.Server
	XClient         *proxyServer.Client
	Wd              selenium.WebDriver
}

// 抓取点播接口集合
type CrawlVodFlowInterface interface {
	CrawlByBrowserInterface

	// 下载视频接口集合
	ConvertM3U8(seed *model.CmsCrawlQueue, filePath string) (*ffmpeg.M3u8DO, error)
	ConvertM3U8GetBaseUrl(m3u8Url string) string
	DownLoadToMp4(m3u8DO *ffmpeg.M3u8DO) error
}

type AbstractCrawlVodFlow struct {
	CrawlByBrowserInterface
	*AbstractCrawlByBrowser
}

func (r *AbstractCrawlVodFlow) UseBrowser() bool {
	return true
}

func (r *AbstractCrawlVodFlow) UseCrawlerProxy() bool {
	return false
}

func (r *AbstractCrawlVodFlow) UseBrowserMobProxy() bool {
	return true
}

func (r *AbstractCrawlVodFlow) OpenBrowser(ctx *BrowserContext) {
}

func (r *AbstractCrawlVodFlow) OpenBrowserWithParams(ctx *BrowserContext, json *gjson.Json) {
}

func (r *AbstractCrawlVodFlow) FillTargetRequest(ctx *BrowserContext) {
}

func (r *AbstractCrawlVodFlow) ConvertM3U8(seed *model.CmsCrawlQueue, filePath string) (*ffmpeg.M3u8DO, error) {
	baseUrl := r.ConvertM3U8GetBaseUrl(seed.CrawlM3U8Url)
	return ffmpeg.ConvertM3U8(seed.CrawlM3U8Url, baseUrl, filePath)
}

func (r *AbstractCrawlVodFlow) ConvertM3U8GetBaseUrl(m3u8Url string) string {
	return http.GetBaseUrlBySchema(m3u8Url)
}

func (r *AbstractCrawlVodFlow) DownLoadToMp4(m3u8DO *ffmpeg.M3u8DO) error {
	return ffmpeg.DownLoadToMp4(m3u8DO)
}
