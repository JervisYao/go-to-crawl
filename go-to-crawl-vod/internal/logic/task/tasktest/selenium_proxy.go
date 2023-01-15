package tasktest

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/tebeka/selenium"
	"go-to-crawl-vod/internal/logic/task/taskdto"
	"go-to-crawl-vod/internal/service/infra/configservice"
	"go-to-crawl-vod/internal/service/infra/webproxyservice"
	"go-to-crawl-vod/utility/browsermobutil"
	"go-to-crawl-vod/utility/chromeutil"
	"go-to-crawl-vod/utility/processutil"
)

var BrowserTask = new(browser)

type browser struct {
}

// 测试web容器环境下多次打开和关闭浏览器代理是否会暂停应用
func (crawlUrl *browser) Start() {

	pid, _ := processutil.CheckRunning(webproxyservice.PORT)
	if pid != "" {
		return
	}

	ctx := new(taskdto.BrowserContext)
	ctx.Log = g.Log().Line()

	service, _ := chromeutil.GetChromeDriverService(chromeutil.DriverServicePort)
	ctx.Service = service
	defer ctx.Service.Stop()

	xServer := webproxyservice.NewServer(configservice.GetCrawlCfg("browserProxyPath"))
	xServer.Start()
	ctx.XServer = xServer
	defer ctx.XServer.Stop()
	proxy := xServer.CreateProxy(nil)
	ctx.XClient = proxy
	defer ctx.XClient.Close()

	// BrowserMobProxy抓包方式
	caps := chromeutil.GetAllCaps(proxy)

	webDriver, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", chromeutil.DriverServicePort))
	ctx.Wd = webDriver
	if ctx.Wd == nil {
		ctx.CrawlQueueSeed.ErrorMsg = "WebDriver Init Fail"
		return
	}
	defer ctx.Wd.Quit()

	browsermobutil.NewHarWait(ctx.Wd, ctx.XClient)
	_ = ctx.Wd.Get("https://www.tangrenjie.tv/vod/play/id/214425/sid/1/nid/1.html")
	browsermobutil.GetHarRequest(ctx.XClient, ".m3u8", "", 5)
}
