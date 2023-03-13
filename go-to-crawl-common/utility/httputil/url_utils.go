package httputil

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	url2 "net/url"
	"strings"
)

// 通过Schema获取基地址
func GetBaseUrlBySchema(m3u8Url string) string {
	log := g.Log().Line()

	urlObj, err := url2.Parse(m3u8Url)
	if err != nil {
		log.Error(err)
		return ""
	}

	if urlObj.Scheme == "" || urlObj.Host == "" {
		log.Error("scheme or host is empty")
		return ""
	}

	baseUrl := fmt.Sprintf("%s://%s", urlObj.Scheme, urlObj.Host)
	log.Info("base url: ", baseUrl)
	return baseUrl
}

// 通过反斜杠获取基地址
func GetBaseUrlByBackslash(m3u8Url string) string {
	log := g.Log().Line()
	idx := strings.LastIndex(m3u8Url, "/")
	baseUrl := gstr.SubStr(m3u8Url, 0, idx+1)
	log.Info("base url: ", baseUrl)
	return baseUrl
}
