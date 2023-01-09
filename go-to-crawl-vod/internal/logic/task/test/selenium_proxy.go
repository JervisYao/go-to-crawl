package test

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/tebeka/selenium"
	"go-to-crawl-vod/internal/logic/task/dto"
	"go-to-crawl-vod/internal/service/infra/config"
	"go-to-crawl-vod/utility/browsermob"
	"go-to-crawl-vod/utility/chrome"
	"go-to-crawl-vod/utility/process"
)

var BrowserTask = new(browser)

type browser struct {
}

// 测试web容器环境下多次打开和关闭浏览器代理是否会暂停应用
func (crawlUrl *browser) Start() {

	pid, _ := process.CheckRunning(proxyServer.PORT)
	if pid != "" {
		return
	}

	ctx := new(dto.BrowserContext)
	ctx.Log = g.Log().Line()

	service, _ := chrome.GetChromeDriverService(chrome.DriverServicePort)
	ctx.Service = service
	defer ctx.Service.Stop()

	xServer := proxyServer.NewServer(config.GetCrawlCfg("browserProxyPath"))
	xServer.Start()
	ctx.XServer = xServer
	defer ctx.XServer.Stop()
	proxy := xServer.CreateProxy(nil)
	ctx.XClient = proxy
	defer ctx.XClient.Close()

	// BrowserMobProxy抓包方式
	caps := chrome.GetAllCaps(proxy)

	webDriver, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", chrome.DriverServicePort))
	ctx.Wd = webDriver
	if ctx.Wd == nil {
		ctx.CrawlQueueSeed.ErrorMsg = "WebDriver Init Fail"
		return
	}
	defer ctx.Wd.Quit()

	browsermob.NewHarWait(ctx.Wd, ctx.XClient)
	_ = ctx.Wd.Get("https://www.tangrenjie.tv/vod/play/id/214425/sid/1/nid/1.html")
	browsermob.GetHarRequest(ctx.XClient, ".m3u8", "", 5)
}