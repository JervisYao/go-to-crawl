package config

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func GetCrawlHostIp() string {
	return GetCrawlCfg("hostIp")
}

func GetCrawlCfg(key string) string {
	value, _ := g.Cfg().Get(gctx.GetInitCtx(), fmt.Sprintf("crawl.%s", key))
	return value.String()
}

func GetCrawlBool(key string) bool {
	value, _ := g.Cfg().Get(gctx.GetInitCtx(), fmt.Sprintf("crawl.%s", key), false)
	return value.Bool()
}

func GetCrawlDebugBool(key string) bool {
	value, _ := g.Cfg().Get(gctx.GetInitCtx(), fmt.Sprintf("crawl.debug.%s", key), false)
	return value.Bool()
}
