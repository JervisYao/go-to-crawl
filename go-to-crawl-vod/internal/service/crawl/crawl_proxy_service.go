package crawl

import (
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"go-to-crawl-vod/internal/dao"
	netUrl "net/url"
	"strings"
)

var (
	C      = dao.CrawlProxy.Columns()
	regTop = "[^.]+\\.(com.cn|com|net.cn|net|org.cn|org|gov.cn|gov|cn|mobi|me|info|name|biz|cc|tv|asia|hk|网络|公司|中国)"
)

func GetProxyByUrl(requestUrl string) string {

	if requestUrl == "" {
		return ""
	}

	url, err := netUrl.Parse(requestUrl)
	if err != nil {
		return ""
	}

	host := url.Host
	index := strings.LastIndex(host, ":")
	if index > 0 {
		host = gstr.SubStr(host, 0, index)
	}
	matches, _ := gregex.MatchString(regTop, host)
	do, _ := dao.CrawlProxy.Where(C.TargetDomain, matches[0]).And(C.ProxyStatus, CrawProxyOpen).FindOne()
	if do != nil {
		return do.ProxyUrl
	}
	return ""
}
