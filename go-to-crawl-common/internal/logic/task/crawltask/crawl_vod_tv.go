package crawltask

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/tebeka/selenium"
	"go-to-crawl-common/consts"
	"go-to-crawl-common/internal/logic/task/taskdto"
	"go-to-crawl-common/utility/browsermobutil"
	"go-to-crawl-common/utility/chromeutil"
	dao2 "go-to-crawl-video/internal/dao"
	entity2 "go-to-crawl-video/internal/model/entity"
	crawlservice2 "go-to-crawl-video/internal/service/crawlservice"
	"go-to-crawl-video/internal/service/infra/configservice"
	"go-to-crawl-video/internal/service/infra/lockservice"
	"go-to-crawl-video/internal/service/infra/webproxyservice"
	"go-to-crawl-video/internal/service/serviceentity"
)

var CrawlVodTVTask = new(crawlVodTVTask)

type crawlVodTVTask struct {
}

func (crawlUrl *crawlVodTVTask) GenVodConfigTask(ctx gctx.Ctx) {
	vodConfig := crawlservice2.GetVodConfig()
	if vodConfig == nil {
		return
	}
	crawlservice2.UpdateVodConfig(vodConfig)

	configTask := new(entity2.CrawlVodConfigTask)
	configTask.VodConfigId = vodConfig.Id
	configTask.TaskStatus = crawlservice2.ConfigTaskStatusInit
	configTask.CreateTime = gtime.Now()

	dao2.CrawlVodConfigTask.Ctx(gctx.GetInitCtx()).Insert(configTask)

}

func (crawlUrl *crawlVodTVTask) VodTVTask(ctx gctx.Ctx) {
	locked := lockservice.TryLockSeleniumLong()

	if !locked {
		return
	}
	defer lockservice.ReleaseLockSelenium()
	vodConfigTaskDO := crawlservice2.GetVodConfigTaskDO()
	if vodConfigTaskDO != nil {
		crawlservice2.UpdateVodConfigTaskStatus(vodConfigTaskDO.CrawlVodConfigTask, crawlservice2.ConfigTaskStatusProcessing)
		DoStartCrawlVodTV(vodConfigTaskDO)
		crawlservice2.UpdateVodConfigTaskStatus(vodConfigTaskDO.CrawlVodConfigTask, crawlservice2.ConfigTaskStatusOk)
	}
}

// 填充视频基础信息
func (crawlUrl *crawlVodTVTask) VodTVPadInfoTask(ctx gctx.Ctx) {
	log := g.Log().Line()
	locked := lockservice.TryLockSelenium()
	if !locked {
		return
	}
	defer lockservice.ReleaseLockSelenium()

	vodTv := crawlservice2.GetVodTvByStatus(crawlservice2.CrawlTVInit)
	if vodTv == nil {
		return
	}

	log.Infof(gctx.GetInitCtx(), "更新vod tv. id = %v, to status = %v", vodTv.Id, crawlservice2.CrawlTVPadInfo)
	crawlservice2.UpdateVodTVStatus(vodTv, crawlservice2.CrawlTVPadInfo)

	DoStartCrawlVodPadInfo(vodTv)
}

// 填充视频ID
func (crawlUrl *crawlVodTVTask) VodTVPadIdTask(ctx gctx.Ctx) {
	log := g.Log().Line()

	// 把填充视频基础信息的记录拿出来填充视频ID
	vodTV := crawlservice2.GetVodTvByStatus(crawlservice2.CrawlTVPadInfoOK)
	if vodTV == nil {
		return
	}
	crawlservice2.UpdateVodTVStatus(vodTV, crawlservice2.CrawlTVPadId)

	log.Infof(gctx.GetInitCtx(), "更新vod tv. id = %v, to status = %v", vodTV.Id, crawlservice2.CrawlTVPadIdOk)
	crawlservice2.UpdateVodTVStatus(vodTV, crawlservice2.CrawlTVPadIdOk)
}

func (crawlUrl *crawlVodTVTask) VodTVItemPadIdTask(ctx gctx.Ctx) {
	log := g.Log().Line()

	vodTVItem := crawlservice2.GetPreparedVodTvItem()
	if vodTVItem == nil {
		return
	}
	vodTv := crawlservice2.GetVodTvById(vodTVItem.TvId)

	log.Infof(gctx.GetInitCtx(), "update vod tv item. id = %v, status = %v", vodTVItem.Id, vodTVItem.CrawlStatus)
	crawlservice2.UpdateVodTVItemStatus(vodTVItem, crawlservice2.CrawlTVItemPadId)

	if vodTv == nil {
		return
	}
	log.Infof(gctx.GetInitCtx(), "addSubId:%v", vodTv.Id)
	if vodTv.Id == 0 || g.NewVar(vodTv.Id).IsEmpty() {
		return
	}

	crawlservice2.UpdateVodTVItemStatus(vodTVItem, crawlservice2.CrawlTVItemPadIdOk)
	transToCrawlQueue(vodTv, vodTVItem)
}

// 转换到爬取队列走标准化抓取
func transToCrawlQueue(vodTv *entity2.CrawlVod, vodTvItem *entity2.CrawlVodItem) {

	vodConfig := crawlservice2.GetVodConfigById(vodTv.VodConfigId)
	businessType := crawlservice2.BusinessTypeNiVod
	if vodConfig != nil && vodConfig.BusinessType != 0 {
		businessType = vodConfig.BusinessType
	}
	crawlQueue := new(entity2.CrawlQueue)
	hostLabelVar, _ := g.Cfg().Get(gctx.GetInitCtx(), "crawl.hostLabel")
	crawlQueue.HostLabel = hostLabelVar.String()
	crawlQueue.BusinessType = businessType

	crawlQueue.VideoYear = gconv.Int(vodTv.VideoYear)
	crawlQueue.VideoCollId = vodTv.VideoCollId
	crawlQueue.CrawlType = crawlservice2.QueueTypePageUrl
	crawlQueue.CrawlStatus = crawlservice2.Init
	crawlQueue.CrawlSeedUrl = vodTvItem.SeedUrl
	crawlQueue.CountryCode = getCountryCodeByString(vodTv.VideoCountry)
	crawlQueue.CrawlSeedParams = ""
	crawlQueue.CreateTime = gtime.Now()

	_, _ = dao2.CrawlQueue.Ctx(gctx.GetInitCtx()).Save(crawlQueue)
}
func getCountryCodeByString(country string) string {
	//大陆 香港 台湾 日本 韩国 欧美 泰国 新马 其它
	defaultCry := "OTHER"

	if country == "" {
		return defaultCry
	}

	var cryMap = make(map[string]string)
	cryMap["大陆"] = "CN"
	cryMap["香港"] = "HK"
	cryMap["台湾"] = "TW"
	cryMap["日本"] = "JP"
	cryMap["韩国"] = "KR"
	cryMap["欧美"] = "US"
	cryMap["泰国"] = "TH"
	cryMap["新马"] = "MY"
	cryMap["其它"] = defaultCry
	if v, ok := cryMap[country]; ok {
		return v
	} else {
		return defaultCry
	}
}

func DoStartCrawlVodTV(configTaskDO *serviceentity.CrawlVodConfigTaskDO) {
	ctx := new(taskdto.BrowserContext)
	ctx.Log = g.Log().Line()
	ctx.VodConfigTaskDO = configTaskDO
	strategy := getCrawlVodTVStrategy(ctx.VodConfigTaskDO.CrawlVodConfig)
	if strategy.UseBrowser() {
		//g.Dump("使用浏览器")
		service, _ := chromeutil.GetChromeDriverService(chromeutil.DriverServicePort)
		ctx.Service = service
		defer ctx.Service.Stop()

		proxyUrl := ""
		if strategy.UseCrawlerProxy() {
			proxyUrl = crawlservice2.GetRandomProxyUrl()
			ctx.Log.Infof(gctx.GetInitCtx(), "visit list page via proxy. domain = %v, proxy = %v", configTaskDO.DomainKeyPart, proxyUrl)
		}
		caps := chromeutil.GetAllCapsChooseProxy(nil, proxyUrl)

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
			ctx.VodConfigTaskDO.CrawlVodConfig.ErrorMsg = "WebDriver Init Fail"
			ctx.Log.Error(gctx.GetInitCtx(), err)
			return
		}
		defer ctx.Wd.Quit()

		// 业务处理-start
		if ctx.VodConfigTaskDO.CrawlVodConfig.SeedParams != "" {
			json, _ := gjson.LoadJson(ctx.VodConfigTaskDO.CrawlVodConfig.SeedParams)
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
}

func DoStartCrawlVodPadInfo(vodTVItem *entity2.CrawlVod) {
	ctx := new(taskdto.BrowserContext)
	ctx.Log = g.Log().Line()
	ctx.VodTV = vodTVItem
	strategy := getCrawlVodPadInfoStrategy(ctx.VodTV)
	if strategy.UseBrowser() {
		//g.Dump("使用浏览器")
		service, _ := chromeutil.GetChromeDriverService(chromeutil.DriverServicePort)
		ctx.Service = service
		defer ctx.Service.Stop()

		proxyUrl := ""
		if strategy.UseCrawlerProxy() {
			proxyUrl = crawlservice2.GetRandomProxyUrl()
			ctx.Log.Infof(gctx.GetInitCtx(), "visit detail page via proxy. id = %v, proxy = %v", vodTVItem.Id, proxyUrl)
		}
		caps := chromeutil.GetAllCapsChooseProxy(nil, proxyUrl)

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
			ctx.VodTVItem.ErrorMsg = "WebDriver Init Fail"
			ctx.Log.Error(gctx.GetInitCtx(), err)
			return
		}
		defer ctx.Wd.Quit()

		// 业务处理-start
		if strategy.UseBrowserMobProxy() {
			browsermobutil.NewHarWait(ctx.Wd, ctx.XClient)
		}
		strategy.OpenBrowser(ctx)
		// 业务处理-end
	}

	// 把URL,Headers信息保存起来
	strategy.FillTargetRequest(ctx)
}
