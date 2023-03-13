package taskdto

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/tebeka/selenium"
	"go-to-crawl-common/utility/ffmpegutil"
	"go-to-crawl-common/utility/httputil"
	entity2 "go-to-crawl-video/internal/model/entity"
	webproxyservice2 "go-to-crawl-video/internal/service/infra/webproxyservice"
	"go-to-crawl-video/internal/service/serviceentity"
)

type BrowserContext struct {
	Log             *glog.Logger
	CrawlQueueSeed  *entity2.CrawlQueue
	VodConfigTaskDO *serviceentity.CrawlVodConfigTaskDO
	VodTV           *entity2.CrawlVod
	VodTVItem       *entity2.CrawlVodItem
	Service         *selenium.Service
	XServer         *webproxyservice2.Server
	XClient         *webproxyservice2.Client
	Wd              selenium.WebDriver
}

// 抓取点播接口集合
type CrawlVodFlowInterface interface {
	CrawlByBrowserInterface

	// 下载视频接口集合
	ConvertM3U8(seed *entity2.CrawlQueue, filePath string) (*ffmpegutil.M3u8DO, error)
	ConvertM3U8GetBaseUrl(m3u8Url string) string
	DownLoadToMp4(m3u8DO *ffmpegutil.M3u8DO) error
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

func (r *AbstractCrawlVodFlow) ConvertM3U8(seed *entity2.CrawlQueue, filePath string) (*ffmpegutil.M3u8DO, error) {
	baseUrl := r.ConvertM3U8GetBaseUrl(seed.CrawlM3U8Url)
	return ffmpegutil.ConvertM3U8(seed.CrawlM3U8Url, baseUrl, filePath)
}

func (r *AbstractCrawlVodFlow) ConvertM3U8GetBaseUrl(m3u8Url string) string {
	return httputil.GetBaseUrlBySchema(m3u8Url)
}

func (r *AbstractCrawlVodFlow) DownLoadToMp4(m3u8DO *ffmpegutil.M3u8DO) error {
	return ffmpegutil.DownLoadToMp4(m3u8DO)
}
