package crawltask

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/tebeka/selenium"
	"go-to-crawl-common/consts"
	"go-to-crawl-common/internal/logic/task/taskdto"
	"go-to-crawl-common/utility/browsermobutil"
	"go-to-crawl-common/utility/chromeutil"
	"go-to-crawl-video/internal/model/entity"
	crawlservice2 "go-to-crawl-video/internal/service/crawlservice"
	"go-to-crawl-video/internal/service/infra/configservice"
	"go-to-crawl-video/internal/service/infra/lockservice"
	"go-to-crawl-video/internal/service/infra/webproxyservice"
)

var CrawlTask = new(CrawlUrl)

type CrawlUrl struct {
}

func (crawlUrl *CrawlUrl) CrawlUrlTask(ctx gctx.Ctx) {
	if !lockservice.TryLockSelenium() {
		return
	}
	defer lockservice.ReleaseLockSelenium()

	seed := getEnvPreparedSeed("", crawlservice2.BusinessTypeNormal)
	if seed != nil {
		DoStartCrawlVodFlow(seed)
	}

}

func (crawlUrl *CrawlUrl) CrawlUrlType1Task(ctx gctx.Ctx) {
	if !lockservice.TryLockSelenium() {
		return
	}
	defer lockservice.ReleaseLockSelenium()

	seed := getEnvPreparedSeed("", crawlservice2.BusinessTypeCrawlLogin)
	if seed != nil {
		DoStartCrawlVodFlow(seed)
	}
}

func getEnvPreparedSeed(hostname string, hostType int) *entity.CrawlQueue {
	seed := crawlservice2.GetSeed(crawlservice2.Init, hostname, hostType)
	if seed == nil {
		return nil
	}
	crawlservice2.UpdateStatus(seed, crawlservice2.Crawling)

	return seed
}

func DoStartCrawlVodFlow(seed *entity.CrawlQueue) {
	ctx := new(taskdto.BrowserContext)
	ctx.Log = g.Log().Line()
	ctx.CrawlQueueSeed = seed
	strategy := getCrawlVodFlowStrategy(ctx.CrawlQueueSeed)
	if strategy.UseBrowser() {
		//g.Dump("使用浏览器")
		service, _ := chromeutil.GetChromeDriverService(chromeutil.DriverServicePort)
		ctx.Service = service
		defer ctx.Service.Stop()
		caps := chromeutil.GetAllCaps(nil)

		if strategy.UseBrowserMobProxy() {
			xServer := webproxyservice.NewServer(configservice.GetCrawlCfg(consts.CrawlBrowserProxyPath))
			xServer.Start()
			ctx.XServer = xServer
			defer ctx.XServer.Stop()
			proxy := xServer.CreateProxy(nil)
			ctx.XClient = proxy
			defer ctx.XClient.Close()

			// BrowserMobProxy抓包方式
			caps = chromeutil.GetAllCaps(proxy)
		}

		webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", chromeutil.DriverServicePort))
		ctx.Wd = webDriver
		if ctx.Wd == nil {
			ctx.CrawlQueueSeed.ErrorMsg = "WebDriver Init Fail"
			ctx.Log.Error(gctx.GetInitCtx(), err)
			crawlservice2.UpdateStatus(ctx.CrawlQueueSeed, crawlservice2.CrawlErr)
			return
		}
		defer ctx.Wd.Quit()

		// 业务处理-start
		if ctx.CrawlQueueSeed.CrawlSeedParams != "" && ctx.CrawlQueueSeed.CrawlSeedParams != `{"videoitem":""}` {
			json, _ := gjson.LoadJson(ctx.CrawlQueueSeed.CrawlSeedParams)
			strategy.OpenBrowserWithParams(ctx, json)
		} else {
			if strategy.UseBrowserMobProxy() {
				browsermobutil.NewHarWait(ctx.Wd, ctx.XClient)
			}
			//g.Dump("打开浏览器")
			strategy.OpenBrowser(ctx)
		}
		// 业务处理-end
	}

	// 把URL,Headers信息保存起来
	strategy.FillTargetRequest(ctx)
	crawlservice2.UpdateUrlAndStatus(ctx.CrawlQueueSeed)
}
