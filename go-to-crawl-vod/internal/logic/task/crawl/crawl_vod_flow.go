package crawl

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/tebeka/selenium"
	"go-to-crawl-vod/internal/logic/task/dto"
	"go-to-crawl-vod/internal/model/entity"
	"go-to-crawl-vod/internal/service/crawl"
	proxyServer "go-to-crawl-vod/internal/service/infra/browsermobproxy"
	"go-to-crawl-vod/internal/service/infra/config"
	"go-to-crawl-vod/internal/service/infra/lock"
	"go-to-crawl-vod/utility/browsermob"
	"go-to-crawl-vod/utility/chrome"
)

var CrawlTask = new(CrawlUrl)

type CrawlUrl struct {
}

func (crawlUrl *CrawlUrl) CrawlUrlTask(ctx gctx.Ctx) {
	if !lock.TryLockSelenium() {
		return
	}
	defer lock.ReleaseLockSelenium()

	seed := getEnvPreparedSeed("", crawl.BusinessTypeNormal)
	if seed != nil {
		DoStartCrawlVodFlow(seed)
	}

}

func (crawlUrl *CrawlUrl) CrawlUrlType1Task(ctx gctx.Ctx) {
	if !lock.TryLockSelenium() {
		return
	}
	defer lock.ReleaseLockSelenium()

	seed := getEnvPreparedSeed("", crawl.BusinessTypeCrawlLogin)
	if seed != nil {
		DoStartCrawlVodFlow(seed)
	}
}

func getEnvPreparedSeed(hostname string, hostType int) *entity.CrawlQueue {
	seed := crawl.GetSeed(crawl.Init, hostname, hostType)
	if seed == nil {
		return nil
	}
	crawl.UpdateStatus(seed, crawl.Crawling)

	return seed
}

func DoStartCrawlVodFlow(seed *entity.CrawlQueue) {
	ctx := new(dto.BrowserContext)
	ctx.Log = g.Log().Line()
	ctx.CrawlQueueSeed = seed
	strategy := getCrawlVodFlowStrategy(ctx.CrawlQueueSeed)
	if strategy.UseBrowser() {
		//g.Dump("使用浏览器")
		service, _ := chrome.GetChromeDriverService(chrome.DriverServicePort)
		ctx.Service = service
		defer ctx.Service.Stop()
		caps := chrome.GetAllCaps(nil)

		if strategy.UseBrowserMobProxy() {
			xServer := proxyServer.NewServer(config.GetCrawlCfg("browserProxyPath"))
			xServer.Start()
			ctx.XServer = xServer
			defer ctx.XServer.Stop()
			proxy := xServer.CreateProxy(nil)
			ctx.XClient = proxy
			defer ctx.XClient.Close()

			// BrowserMobProxy抓包方式
			caps = chrome.GetAllCaps(proxy)
		}

		webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", chrome.DriverServicePort))
		ctx.Wd = webDriver
		if ctx.Wd == nil {
			ctx.CrawlQueueSeed.ErrorMsg = "WebDriver Init Fail"
			ctx.Log.Error(gctx.GetInitCtx(), err)
			crawl.UpdateStatus(ctx.CrawlQueueSeed, crawl.CrawlErr)
			return
		}
		defer ctx.Wd.Quit()

		// 业务处理-start
		if ctx.CrawlQueueSeed.CrawlSeedParams != "" && ctx.CrawlQueueSeed.CrawlSeedParams != `{"videoitem":""}` {
			json, _ := gjson.LoadJson(ctx.CrawlQueueSeed.CrawlSeedParams)
			strategy.OpenBrowserWithParams(ctx, json)
		} else {
			if strategy.UseBrowserMobProxy() {
				browsermob.NewHarWait(ctx.Wd, ctx.XClient)
			}
			//g.Dump("打开浏览器")
			strategy.OpenBrowser(ctx)
		}
		// 业务处理-end
	}

	// 把URL,Headers信息保存起来
	strategy.FillTargetRequest(ctx)
	crawl.UpdateUrlAndStatus(ctx.CrawlQueueSeed)
}
