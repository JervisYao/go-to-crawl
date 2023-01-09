package nunuyy

import (
	"easygoadmin/appnewcms/model"
	"easygoadmin/appnewcms/task/dto"
	"easygoadmin/appnewcms/utils/browsermob"
	"easygoadmin/appnewcms/utils/ffmpeg"
	"easygoadmin/appnewcms/utils/http"
	"easygoadmin/appnewcms/utils/selector"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/gtime"
	"github.com/tebeka/selenium"
	"time"
)

var (
	videoXpath        = "//*[@id='video'][@src]"
	sliderXpath       = "//*[@id='slider']"
	resourceNameXpath = "//*[@id='slider']//dt"
)

type NunuyyCrawl struct {
	*dto.AbstractCrawlVodFlow
}

func (c *NunuyyCrawl) OpenBrowser(ctx *dto.BrowserContext) {
	_ = ctx.Wd.WaitWithTimeout(selector.GetXpathCondition(videoXpath), gtime.S*30)
	_ = ctx.Wd.Get(ctx.CrawlQueueSeed.CrawlSeedUrl)
}

func (c *NunuyyCrawl) OpenBrowserWithParams(ctx *dto.BrowserContext, json *gjson.Json) {
	_ = ctx.Wd.WaitWithTimeout(selector.GetXpathCondition(sliderXpath), gtime.S*30)
	_ = ctx.Wd.Get(ctx.CrawlQueueSeed.CrawlSeedUrl)

	resource := json.GetString("resource")
	resource = "量子资源"
	// 网站支持的资源列表
	resElements, _ := ctx.Wd.FindElements(selenium.ByXPATH, "//*[@id='slider']//dt")
	if len(resElements) == 0 {
		ctx.Log.Error("不存在该资源：", resourceNameXpath)
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
	videoItemName := json.GetString("videoItem")
	for _, videoItemElement := range videoItemElements {
		videoItemText, _ := videoItemElement.Text()
		if videoItemName == videoItemText {
			_ = videoItemElement.Click()
			browsermob.NewHarWait(ctx.Wd, ctx.XClient)
			time.Sleep(time.Second)
		}
	}

	// 等待资源加载完成
	_ = ctx.Wd.WaitWithTimeout(selector.GetXpathCondition(videoXpath), gtime.S*30)
}

func (c *NunuyyCrawl) ConvertM3U8(seed *model.CmsCrawlQueue, filePath string) (*ffmpeg.M3u8DO, error) {
	baseUrl := c.ConvertM3U8GetBaseUrl(seed.CrawlM3U8Url)
	return ffmpeg.ConvertM3U8(seed.CrawlSeedUrl, baseUrl, filePath)
}

func (c *NunuyyCrawl) ConvertM3U8GetBaseUrl(m3u8Url string) string {
	return http.GetBaseUrlByBackslash(m3u8Url)
}

func (c *NunuyyCrawl) FillTargetRequest(ctx *dto.BrowserContext) {
	request := browsermob.GetHarRequestLocalRetry(ctx.XClient, ".m3u8", "")
	if request != nil {
		ctx.CrawlQueueSeed.CrawlM3U8Url = request.GetString("url")
	}
}