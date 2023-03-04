package chromeutil

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"go-to-crawl-vod/internal/consts"
	"go-to-crawl-vod/internal/service/infra/configservice"
	"go-to-crawl-vod/internal/service/infra/webproxyservice"
)

const (
	DriverServicePort = 8088
	ScrollBottomJs    = "window.scrollTo(0,document.body.scrollHeight);"
)

func GetChromeDriverService(port int) (*selenium.Service, error) {
	opts := []selenium.ServiceOption{
		//selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	chromeDriverPath := configservice.GetCrawlCfg(consts.CrawlBrowserDriverPath)
	service, _ := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	return service, nil
}

func GetAllCaps(mobProxy *webproxyservice.Client) selenium.Capabilities {
	return GetAllCapsChooseProxy(mobProxy, "")
}

func GetAllCapsChooseProxy(mobProxy *webproxyservice.Client, crawlerProxy string) selenium.Capabilities {
	driverType := configservice.GetCrawlCfg(consts.CrawlBrowserDriverType)
	caps := GetCommonCaps(driverType)

	if driverType == consts.DriverTypeChrome {
		chromeCaps := GetChromeCaps(mobProxy, crawlerProxy)
		caps.AddChrome(chromeCaps)
	} else if driverType == consts.DriverTypeEdge {
		// TODO 待兼容Edge浏览器
		chromeCaps := GetChromeCaps(mobProxy, crawlerProxy)
		caps.AddChrome(chromeCaps)
	}

	return caps
}

func GetCommonCaps(browser string) selenium.Capabilities {
	caps := selenium.Capabilities{
		"browserName": browser,
	}
	return caps
}

func GetChromeCaps(mobProxy *webproxyservice.Client, crawlerProxy string) chrome.Capabilities {
	args := []string{
		"--no-sandbox",
		"--ignore-certificate-errors",
		"--disable-blink-features=AutomationControlled", // 隐藏自己是selenium. window.navigator.webdrive=true
		"--user-agent=" + GetRandomUA(false),
		"--acceptSslCerts=true",
	}

	// headless
	headless := configservice.GetCrawlBool(consts.CrawlBrowserHeadless)
	if headless {
		args = append(args, "--headless")
	}

	// proxy (mobProxy优先级高于crawlerProxy, 因为mobProxy是为了抓包，crawlerProxy是为了防爬)
	if mobProxy != nil {
		args = append(args, "--proxy-server="+mobProxy.Proxy)
	} else {
		if crawlerProxy != "" {
			args = append(args, "--proxy-server="+crawlerProxy)
		}
	}

	// 缓存的用户信息，用于让selenium记录用户登录状态
	userDataDir := configservice.GetCrawlCfg(consts.CrawlBrowserUserDataDir)
	driverType := configservice.GetCrawlCfg(consts.CrawlBrowserDriverType)

	if userDataDir != "" {
		args = append(args, "--user-data-dir="+fmt.Sprintf("%s\\%s", userDataDir, driverType))
	}

	chromeCaps := chrome.Capabilities{
		Path:  configservice.GetCrawlCfg(consts.CrawlBrowserExecutorPath),
		Args:  args,
		Prefs: map[string]interface{}{
			//"profile.managed_default_content_settings.images": 2,
			//"permissions.default.stylesheet": 2,
		},
		ExcludeSwitches: []string{"enable-automation"},
	}
	return chromeCaps
}
