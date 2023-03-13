package crawltask

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/tebeka/selenium"
	"go-to-crawl-common/internal/consts"
	"go-to-crawl-common/internal/dao"
	"go-to-crawl-common/internal/logic/task/taskdto"
	"go-to-crawl-common/internal/model/entity"
	"go-to-crawl-common/internal/service/crawlservice"
	"go-to-crawl-common/internal/service/infra/configservice"
	"go-to-crawl-common/internal/service/infra/lockservice"
	"go-to-crawl-common/internal/service/infra/webproxyservice"
	"go-to-crawl-common/internal/service/serviceentity"
	"go-to-crawl-common/utility/browsermobutil"
	"go-to-crawl-common/utility/chromeutil"
)

var CrawlVodTVTask = new(crawlVodTVTask)

type crawlVodTVTask struct {
}

func (crawlUrl *crawlVodTVTask) GenVodConfigTask(ctx gctx.Ctx) {
	vodConfig := crawlservice.GetVodConfig()
	if vodConfig == nil {
		return
	}
	crawlservice.UpdateVodConfig(vodConfig)

	configTask := new(entity.CrawlVodConfigTask)
	configTask.VodConfigId = vodConfig.Id
	configTask.TaskStatus = crawlservice.ConfigTaskStatusInit
	configTask.CreateTime = gtime.Now()

	dao.CrawlVodConfigTask.Ctx(gctx.GetInitCtx()).Insert(configTask)

}

func (crawlUrl *crawlVodTVTask) VodTVTask(ctx gctx.Ctx) {
	locked := lockservice.TryLockSeleniumLong()

	if !locked {
		return
	}
	defer lockservice.ReleaseLockSelenium()
	vodConfigTaskDO := crawlservice.GetVodConfigTaskDO()
	if vodConfigTaskDO != nil {
		crawlservice.UpdateVodConfigTaskStatus(vodConfigTaskDO.CrawlVodConfigTask, crawlservice.ConfigTaskStatusProcessing)
		DoStartCrawlVodTV(vodConfigTaskDO)
		crawlservice.UpdateVodConfigTaskStatus(vodConfigTaskDO.CrawlVodConfigTask, crawlservice.ConfigTaskStatusOk)
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

	vodTv := crawlservice.GetVodTvByStatus(crawlservice.CrawlTVInit)
	if vodTv == nil {
		return
	}

	log.Infof(gctx.GetInitCtx(), "更新vod tv. id = %v, to status = %v", vodTv.Id, crawlservice.CrawlTVPadInfo)
	crawlservice.UpdateVodTVStatus(vodTv, crawlservice.CrawlTVPadInfo)

	DoStartCrawlVodPadInfo(vodTv)
}

// 填充视频ID
func (crawlUrl *crawlVodTVTask) VodTVPadIdTask(ctx gctx.Ctx) {
	log := g.Log().Line()

	// 把填充视频基础信息的记录拿出来填充视频ID
	vodTV := crawlservice.GetVodTvByStatus(crawlservice.CrawlTVPadInfoOK)
	if vodTV == nil {
		return
	}
	crawlservice.UpdateVodTVStatus(vodTV, crawlservice.CrawlTVPadId)

	log.Infof(gctx.GetInitCtx(), "更新vod tv. id = %v, to status = %v", vodTV.Id, crawlservice.CrawlTVPadIdOk)
	crawlservice.UpdateVodTVStatus(vodTV, crawlservice.CrawlTVPadIdOk)
}

func (crawlUrl *crawlVodTVTask) VodTVItemPadIdTask(ctx gctx.Ctx) {
	log := g.Log().Line()

	vodTVItem := crawlservice.GetPreparedVodTvItem()
	if vodTVItem == nil {
		return
	}
	vodTv := crawlservice.GetVodTvById(vodTVItem.TvId)

	log.Infof(gctx.GetInitCtx(), "update vod tv item. id = %v, status = %v", vodTVItem.Id, vodTVItem.CrawlStatus)
	crawlservice.UpdateVodTVItemStatus(vodTVItem, crawlservice.CrawlTVItemPadId)

	if vodTv == nil {
		return
	}
	log.Infof(gctx.GetInitCtx(), "addSubId:%v", vodTv.Id)
	if vodTv.Id == 0 || g.NewVar(vodTv.Id).IsEmpty() {
		return
	}

	crawlservice.UpdateVodTVItemStatus(vodTVItem, crawlservice.CrawlTVItemPadIdOk)
	transToCrawlQueue(vodTv, vodTVItem)
}

// 转换到爬取队列走标准化抓取
func transToCrawlQueue(vodTv *entity.CrawlVod, vodTvItem *entity.CrawlVodItem) {

	vodConfig := crawlservice.GetVodConfigById(vodTv.VodConfigId)
	businessType := crawlservice.BusinessTypeNiVod
	if vodConfig != nil && vodConfig.BusinessType != 0 {
		businessType = vodConfig.BusinessType
	}
	crawlQueue := new(entity.CrawlQueue)
	hostLabelVar, _ := g.Cfg().Get(gctx.GetInitCtx(), "crawl.hostLabel")
	crawlQueue.HostLabel = hostLabelVar.String()
	crawlQueue.BusinessType = businessType

	crawlQueue.VideoYear = gconv.Int(vodTv.VideoYear)
	crawlQueue.VideoCollId = vodTv.VideoCollId
	crawlQueue.CrawlType = crawlservice.QueueTypePageUrl
	crawlQueue.CrawlStatus = crawlservice.Init
	crawlQueue.CrawlSeedUrl = vodTvItem.SeedUrl
	crawlQueue.CountryCode = getCountryCodeByString(vodTv.VideoCountry)
	crawlQueue.CrawlSeedParams = ""
	crawlQueue.CreateTime = gtime.Now()

	_, _ = dao.CrawlQueue.Ctx(gctx.GetInitCtx()).Save(crawlQueue)
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
			proxyUrl = crawlservice.GetRandomProxyUrl()
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

func DoStartCrawlVodPadInfo(vodTVItem *entity.CrawlVod) {
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
			proxyUrl = crawlservice.GetRandomProxyUrl()
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
