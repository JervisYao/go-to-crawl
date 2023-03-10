package browserutil

import (
	"easygoadmin/app/utils"
	"easygoadmin/appnewcms/service/configservice"
	"easygoadmin/appnewcms/task/taskdto"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/tebeka/selenium"
)

func NewRemote(capabilities selenium.Capabilities, port int, uriPrefix string) (selenium.WebDriver, error) {
	urlPrefix := fmt.Sprintf("http://localhost:%d%s", port, uriPrefix)
	return selenium.NewRemote(capabilities, urlPrefix)
}

func GetDriverService(port int) (*selenium.Service, error) {

	browserConfig := configservice.GetCrawlBrowser()
	browserInfoConfig := configservice.GetCrawlBrowserInfo()

	if utils.DriverTypeChrome == browserConfig.BrowserType {
		driverOpt := selenium.ChromeDriver(browserInfoConfig.DriverPath)
		return selenium.NewChromeDriverService(browserInfoConfig.DriverPath, port, driverOpt)
	} else if utils.DriverTypeFireFox == browserConfig.BrowserType {
		driverOpt := selenium.GeckoDriver(browserInfoConfig.DriverPath)
		outputOpt := selenium.Output(nil)
		return selenium.NewGeckoDriverService(browserInfoConfig.DriverPath, port, driverOpt, outputOpt)
	}

	return nil, nil
}

func GetAllCaps(browserCtx *taskdto.BrowserContext) selenium.Capabilities {
	return GetAllCapsChooseProxy(browserCtx, "")
}

func GetAllCapsChooseProxy(browserCtx *taskdto.BrowserContext, crawlerProxy string) selenium.Capabilities {

	browserConfig := configservice.GetCrawlBrowser()
	if browserConfig == nil {
		g.Log().Error("没有配置WebDriver参数")
		return nil
	}

	browserType := browserConfig.BrowserType
	caps := GetCommonCaps(browserType)

	if browserType == utils.DriverTypeChrome {
		specCaps := GetChromeCaps(browserCtx, crawlerProxy)
		caps.AddChrome(specCaps)
	} else if browserType == utils.DriverTypeFireFox {
		specCaps := GetFirefoxCaps(browserCtx, crawlerProxy)
		caps.AddFirefox(specCaps)
	} else if browserType == utils.DriverTypeEdge {
		// Edge为Chrome内核
		specCaps := GetEdgeCaps(browserCtx, crawlerProxy)
		caps.AddChrome(specCaps)
	}

	return caps
}

func GetCommonCaps(browser string) selenium.Capabilities {
	caps := selenium.Capabilities{
		"browserName": browser,
	}
	return caps
}

func appendProxyArgs(args []string, browserCtx *taskdto.BrowserContext, crawlerProxy string) []string {
	// proxy (mobProxy优先级高于crawlerProxy, 因为mobProxy是为了抓包，crawlerProxy是为了防爬)
	if browserCtx != nil && browserCtx.XClient != nil {
		args = append(args, "--proxy-server="+browserCtx.XClient.Proxy)
	} else {
		if crawlerProxy != "" {
			args = append(args, "--proxy-server="+crawlerProxy)
		}
	}
	return args
}

func appendConfigArgs(args []string) []string {
	browserConfig := configservice.GetCrawlBrowser()
	if browserConfig == nil {
		return args
	}

	// headless
	if browserConfig.Headless {
		args = append(args, "--headless")
	}

	// 谷歌缓存的用户信息，用于让selenium记录用户登录状态
	userDataDir := browserConfig.UserDataDir
	browserType := browserConfig.BrowserType

	if userDataDir != "" {
		args = append(args, "--user-data-dir="+fmt.Sprintf("%s\\%s", userDataDir, browserType))
	}

	return args
}
