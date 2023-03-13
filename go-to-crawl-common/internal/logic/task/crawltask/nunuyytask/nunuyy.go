package nunuyytask

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/tebeka/selenium"
	"go-to-crawl-common/internal/logic/task/taskdto"
	"go-to-crawl-common/internal/model/entity"
	"go-to-crawl-common/utility/browsermobutil"
	"go-to-crawl-common/utility/ffmpegutil"
	"go-to-crawl-common/utility/httputil"
	"go-to-crawl-common/utility/selectorutil"
	"time"
)

var (
	videoXpath        = "//*[@id='video'][@src]"
	sliderXpath       = "//*[@id='slider']"
	resourceNameXpath = "//*[@id='slider']//dt"
)

type NunuyyCrawl struct {
	*taskdto.AbstractCrawlVodFlow
}

func (c *NunuyyCrawl) OpenBrowser(ctx *taskdto.BrowserContext) {
	_ = ctx.Wd.Get(ctx.CrawlQueueSeed.CrawlSeedUrl)
	_ = ctx.Wd.WaitWithTimeout(selectorutil.GetXpathCondition(videoXpath), gtime.S*30)
}

func (c *NunuyyCrawl) OpenBrowserWithParams(ctx *taskdto.BrowserContext, json *gjson.Json) {
	_ = ctx.Wd.Get(ctx.CrawlQueueSeed.CrawlSeedUrl)
	_ = ctx.Wd.WaitWithTimeout(selectorutil.GetXpathCondition(sliderXpath), gtime.S*30)

	resource := json.Get("resource").String()
	resource = "量子资源"
	// 网站支持的资源列表
	resElements, _ := ctx.Wd.FindElements(selenium.ByXPATH, "//*[@id='slider']//dt")
	if len(resElements) == 0 {
		ctx.Log.Error(gctx.GetInitCtx(), "不存在该资源：", resourceNameXpath)
		return
	}

	// 量子资源出现在资源列表的位置idx
	idx := 0
	for i, resEle := range resElements {
		resText, _ := resEle.Text()
		if resource == resText {
			_ = resEle.Click()
			idx = i
			break
		}
	}

	// 所有资源节目单列表
	resProgramElements, _ := ctx.Wd.FindElements(selenium.ByXPATH, "//*[@class='tempWrap']//*[@class='playlist clearfix']")
	// 量子资源节目单
	resProgramElement := resProgramElements[idx+1]
	videoItemElements, _ := resProgramElement.FindElements(selenium.ByXPATH, "ul/li")

	// 通过节目名找到节目并点击
	videoItemName := json.Get("videoItem").String()
	for _, videoItemElement := range videoItemElements {
		videoItemText, _ := videoItemElement.Text()
		if videoItemName == videoItemText {
			_ = videoItemElement.Click()
			browsermobutil.NewHarWait(ctx.Wd, ctx.XClient)
			time.Sleep(time.Second)
		}
	}

	// 等待资源加载完成
	_ = ctx.Wd.WaitWithTimeout(selectorutil.GetXpathCondition(videoXpath), gtime.S*30)
}

func (c *NunuyyCrawl) ConvertM3U8(seed *entity.CrawlQueue, filePath string) (*ffmpegutil.M3u8DO, error) {
	baseUrl := c.ConvertM3U8GetBaseUrl(seed.CrawlM3U8Url)
	return ffmpegutil.ConvertM3U8(seed.CrawlSeedUrl, baseUrl, filePath)
}

func (c *NunuyyCrawl) ConvertM3U8GetBaseUrl(m3u8Url string) string {
	return httputil.GetBaseUrlByBackslash(m3u8Url)
}

func (c *NunuyyCrawl) FillTargetRequest(ctx *taskdto.BrowserContext) {
	request := browsermobutil.GetHarRequestLocalRetry(ctx.XClient, ".m3u8", "")
	if request != nil {
		ctx.CrawlQueueSeed.CrawlM3U8Url = request.Get("url").String()
	}
}
