package crawl

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/tebeka/selenium"
	"go-to-crawl-vod/internal/dao"
	"go-to-crawl-vod/internal/logic/task/dto"
	"go-to-crawl-vod/internal/service/crawl"
	"go-to-crawl-vod/internal/service/infra/config"
	"go-to-crawl-vod/internal/service/infra/lock"
	"go-to-crawl-vod/utility/browsermob"
	"go-to-crawl-vod/utility/chrome"
)

var CrawlVodTVTask = new(crawlVodTVTask)

type crawlVodTVTask struct {
}

func (crawlUrl *crawlVodTVTask) GenVodConfigTask(ctx gctx.Ctx) {
	vodConfig := crawl.GetVodConfig()
	if vodConfig == nil {
		return
	}
	crawl.UpdateVodConfig(vodConfig)

	configTask := new(model.CmsCrawlVodConfigTask)
	configTask.VodConfigId = vodConfig.Id
	configTask.TaskStatus = crawl.ConfigTaskStatusInit
	configTask.CreateTime = gtime.Now()

	dao.CrawlVodConfigTask.Insert(configTask)

}

func (crawlUrl *crawlVodTVTask) VodTVTask(ctx gctx.Ctx) {
	locked := lock.TryLockSeleniumLong()

	if !locked {
		return
	}
	defer lock.ReleaseLockSelenium()
	vodConfigTaskDO := crawl.GetVodConfigTaskDO()
	if vodConfigTaskDO != nil {
		crawl.UpdateVodConfigTaskStatus(vodConfigTaskDO.CmsCrawlVodConfigTask, crawl.ConfigTaskStatusProcessing)
		DoStartCrawlVodTV(vodConfigTaskDO)
		crawl.UpdateVodConfigTaskStatus(vodConfigTaskDO.CmsCrawlVodConfigTask, crawl.ConfigTaskStatusOk)
	}
}

// 填充视频基础信息
func (crawlUrl *crawlVodTVTask) VodTVPadInfoTask(ctx gctx.Ctx) {
	log := g.Log().Line()
	locked := lock.TryLockSelenium()
	if !locked {
		return
	}
	defer lock.ReleaseLockSelenium()

	vodTv := crawl.GetVodTvByStatus(crawl.CrawlTVInit)
	if vodTv == nil {
		return
	}

	log.Infof("更新vod tv. id = %v, to status = %v", vodTv.Id, crawl.CrawlTVPadInfo)
	crawl.UpdateVodTVStatus(vodTv, crawl.CrawlTVPadInfo)

	DoStartCrawlVodPadInfo(vodTv)
}

// 填充视频ID
func (crawlUrl *crawlVodTVTask) VodTVPadIdTask(ctx gctx.Ctx) {
	log := g.Log().Line()

	// 把填充视频基础信息的记录拿出来填充视频ID
	vodTV := crawl.GetVodTvByStatus(crawl.CrawlTVPadInfoOK)
	if vodTV == nil {
		return
	}
	crawl.UpdateVodTVStatus(vodTV, crawl.CrawlTVPadId)

	addSubInfo := appoldcms.EpgCmsDb.GetSubclassFromVodTv(vodTV)
	videoCollId, err := appoldcms.EpgCmsDb.AddSubClass(addSubInfo)
	log.Infof(gctx.GetInitCtx(), "填充的剧集Id = %v", videoCollId)
	if err != nil {
		//	获取数据鼠标
		log.Infof(gctx.GetInitCtx(), "新增剧集失败:%v", err)
		crawl.UpdateVodTVStatus(vodTV, crawl.CrawlTVPadIdErr)
		return
	}

	// 先填充剧集ID
	vodTV.VideoCollId = videoCollId
	log.Infof(gctx.GetInitCtx(), "更新vod tv. id = %v, to status = %v", vodTV.Id, crawl.CrawlTVPadIdOk)
	crawl.UpdateVodTVStatus(vodTV, crawl.CrawlTVPadIdOk)
}

func (crawlUrl *crawlVodTVTask) VodTVItemPadIdTask(ctx gctx.Ctx) {
	log := g.Log().Line()

	vodTVItem := crawl.GetPreparedVodTvItem()
	if vodTVItem == nil {
		return
	}
	vodTv := crawl.GetVodTvById(vodTVItem.TvId)

	vodTVItem.VideoCollId = vodTv.VideoCollId // 修复异常情况
	log.Infof(gctx.GetInitCtx(), "update vod tv item. id = %v, status = %v", vodTVItem.Id, vodTVItem.CrawlStatus)
	crawl.UpdateVodTVItemStatus(vodTVItem, crawl.CrawlTVItemPadId)

	if vodTv == nil {
		return
	}
	log.Infof(gctx.GetInitCtx(), "addSubId:%v", vodTv.Id)
	if vodTv.Id == 0 || g.NewVar(vodTv.Id).IsEmpty() {
		return
	}

	var addSubInfo = modelOld.Subclass{
		Id:       g.NewVar(vodTv.VideoCollId).Uint(),
		Showname: vodTv.VideoName,
	}

	videoItemId, err := appoldcms.EpgCmsDb.AddSubClassCon(addSubInfo, vodTVItem.Episodes)
	if err != nil {
		//	获取数据鼠标
		log.Infof(gctx.GetInitCtx(), "新增级数失败:%v", err)
		crawl.UpdateVodTVItemStatus(vodTVItem, crawl.CrawlTVItemPadIdErr)
		return
	}

	vodTVItem.VideoItemId = videoItemId
	crawl.UpdateVodTVItemStatus(vodTVItem, crawl.CrawlTVItemPadIdOk)
	transToCrawlQueue(vodTv, vodTVItem)
}

// 转换到爬取队列走标准化抓取
func transToCrawlQueue(vodTv *model.CmsCrawlVodTv, vodTvItem *model.CmsCrawlVodTvItem) {

	vodConfig := crawl.GetVodConfigById(vodTv.VodConfigId)
	hostType := crawl.HostTypeNiVod
	if vodConfig != nil && vodConfig.HostType != 0 {
		hostType = vodConfig.HostType
	}
	crawlQueue := new(model.CmsCrawlQueue)
	crawlQueue.HostIp = g.Cfg().GetString("crawl.down_ip")
	crawlQueue.HostType = hostType

	crawlQueue.VideoYear = gconv.Int(vodTv.VideoYear)
	crawlQueue.VideoCollId = vodTv.VideoCollId
	crawlQueue.VideoItemId = vodTvItem.VideoItemId
	crawlQueue.CrawlType = crawl.QueueTypePageUrl
	crawlQueue.CrawlStatus = crawl.Init
	crawlQueue.CrawlSeedUrl = vodTvItem.SeedUrl
	crawlQueue.CountryCode = getCountryCodeByString(vodTv.VideoCountry)
	crawlQueue.CrawlSeedParams = ""
	crawlQueue.CreateTime = gtime.Now()

	_, _ = dao.CrawlQueue.Save(crawlQueue)
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

func DoStartCrawlVodTV(configTaskDO *do.CrawlVodConfigTaskDO) {
	ctx := new(dto.BrowserContext)
	ctx.Log = g.Log().Line()
	ctx.VodConfigTaskDO = configTaskDO
	strategy := getCrawlVodTVStrategy(ctx.VodConfigTaskDO.CmsCrawlVodConfig)
	if strategy.UseBrowser() {
		//g.Dump("使用浏览器")
		service, _ := chrome.GetChromeDriverService(chrome.DriverServicePort)
		ctx.Service = service
		defer ctx.Service.Stop()

		proxyUrl := ""
		if strategy.UseCrawlerProxy() {
			proxyUrl = crawl.GetRandomProxyUrl()
			ctx.Log.Infof(gctx.GetInitCtx(), "visit list page via proxy. domain = %v, proxy = %v", configTaskDO.DomainKeyPart, proxyUrl)
		}
		caps := chrome.GetAllCapsChooseProxy(nil, proxyUrl)

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
			ctx.VodConfigTaskDO.CmsCrawlVodConfig.ErrorMsg = "WebDriver Init Fail"
			ctx.Log.Error(gctx.GetInitCtx(), err)
			return
		}
		defer ctx.Wd.Quit()

		// 业务处理-start
		if ctx.VodConfigTaskDO.CmsCrawlVodConfig.SeedParams != "" {
			json, _ := gjson.LoadJson(ctx.VodConfigTaskDO.CmsCrawlVodConfig.SeedParams)
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
}

func DoStartCrawlVodPadInfo(vodTVItem *model.CmsCrawlVodTv) {
	ctx := new(dto.BrowserContext)
	ctx.Log = g.Log().Line()
	ctx.VodTV = vodTVItem
	strategy := getCrawlVodPadInfoStrategy(ctx.VodTV)
	if strategy.UseBrowser() {
		//g.Dump("使用浏览器")
		service, _ := chrome.GetChromeDriverService(chrome.DriverServicePort)
		ctx.Service = service
		defer ctx.Service.Stop()

		proxyUrl := ""
		if strategy.UseCrawlerProxy() {
			proxyUrl = crawl.GetRandomProxyUrl()
			ctx.Log.Infof(gctx.GetInitCtx(), "visit detail page via proxy. id = %v, proxy = %v", vodTVItem.Id, proxyUrl)
		}
		caps := chrome.GetAllCapsChooseProxy(nil, proxyUrl)

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
			ctx.VodTVItem.ErrorMsg = "WebDriver Init Fail"
			ctx.Log.Error(gctx.GetInitCtx(), err)
			return
		}
		defer ctx.Wd.Quit()

		// 业务处理-start
		if strategy.UseBrowserMobProxy() {
			browsermob.NewHarWait(ctx.Wd, ctx.XClient)
		}
		strategy.OpenBrowser(ctx)
		// 业务处理-end
	}

	// 把URL,Headers信息保存起来
	strategy.FillTargetRequest(ctx)
}
