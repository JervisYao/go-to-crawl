package crawltask

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/tebeka/selenium"
	"go-to-crawl-common/internal/service/crawlservice"
	"go-to-crawl-common/utility/chromeutil"
	"go-to-crawl-common/utility/selectorutil"
	"go-to-crawl-common/utility/timeutil"
	"time"
)

const (
	vipDescXpath   = "//*[@class='desc _vip_desc']"
	vipHistoryPage = "https://v.qq.com/biu/u/history/"
)

var CrawlCheckTask = new(crawlCheckTask)

type crawlCheckTask struct {
}

func (crawlUrl *crawlCheckTask) CheckQQLoginTask(ctx gctx.Ctx) {
	crawlUrl.LoginQQ(false)
}

// 手动登录QQ(N秒内操作)
// 然后把数据库里QQ类型爬取M3U8失败的记录状态重置为初始化
func (crawlUrl *crawlCheckTask) LoginQQManual() {
	crawlUrl.LoginQQ(true)
}

func (crawlUrl *crawlCheckTask) LoginQQ(waitScan bool) {
	log := g.Log().Line()
	log.Info(gctx.GetInitCtx(), "开始检测腾讯视频VIP登录，过期状态")
	service, _ := chromeutil.GetChromeDriverService(chromeutil.DriverServicePort)
	defer service.Stop()

	caps := chromeutil.GetAllCaps(nil)
	webDriver, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", chromeutil.DriverServicePort))
	if webDriver == nil {
		return
	}
	defer webDriver.Quit()

	_ = webDriver.WaitWithTimeout(selectorutil.GetXpathCondition(vipDescXpath), gtime.S*30)
	_ = webDriver.Get(vipHistoryPage)

	vipDesc := selectorutil.GetTextByXpath(webDriver, vipDescXpath)

	// 1、检测腾讯视频是否登录, 是否是VIP，是否马上到期
	if gstr.Contains(vipDesc, "到期") {
		// 是VIP
		dateStr, _ := gregex.MatchString("(\\d{4}-\\d{2}-\\d{2})", vipDesc)
		expireDate := gtime.ParseTimeFromContent(dateStr[1], timeutil.YYYY_MM_DD)

		milliSeconds := expireDate.TimestampMilli() - gtime.TimestampMilli()
		expireDay := milliSeconds / gtime.D.Milliseconds()

		log.Infof(gctx.GetInitCtx(), "腾讯视频VIP %d天后过期", expireDay)
		if expireDay < 7 {
			// 提前7天告警提示运维处理
			log.Error(gctx.GetInitCtx(), "VIP即将过期")
		}
	} else {

		if waitScan {
			time.Sleep(gtime.S * 100)
		} else {
			// 对接告警渠道
			log.Error(gctx.GetInitCtx(), "请用未过期VIP账号登录")
		}

	}
}

func (crawlUrl *crawlCheckTask) CrawlUrlFailNotifyTask(ctx gctx.Ctx) {
	list := crawlservice.GetNeedNotifySeedList()
	if len(list) == 0 {
		return
	}

	log := g.Log().Line()
	log.Errorf(gctx.GetInitCtx(), "重试多次失败短信告警. size = %d", len(list))

	for _, seed := range list {
		log.Errorf(gctx.GetInitCtx(), "失败seed. url = %s", seed.CrawlSeedUrl)
		seed.CrawlM3U8Notify = crawlservice.CrawlM3U8NotifyYes
		crawlservice.UpdateStatus(seed, seed.CrawlStatus)
	}
}
