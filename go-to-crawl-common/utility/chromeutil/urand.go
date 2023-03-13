package chromeutil

import (
	"github.com/corpix/uarand"
	"github.com/gogf/gf/v2/text/gstr"
)

func GetRandomUA(includeMobile bool) string {
	ua := uarand.GetRandom()
	if includeMobile {
		return ua
	}

	if gstr.ContainsI(ua, "android") || gstr.ContainsI(ua, "iphone") {
		return GetRandomUA(includeMobile)
	}

	return ua
}
