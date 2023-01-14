package config

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func GetCrawlHostLabel() string {
	return GetCrawlCfg("hostLabel")
}

func GetCrawlCfg(key string) string {
	key = fmt.Sprintf("crawl.%s", key)
	value, _ := g.Cfg().Get(gctx.GetInitCtx(), key)
	return value.String()
}

func GetCrawlBool(key string) bool {
	key = fmt.Sprintf("crawl.%s", key)
	value, _ := g.Cfg().Get(gctx.GetInitCtx(), key, false)
	return value.Bool()
}

func GetCrawlDebugBool(key string) bool {
	key = fmt.Sprintf("crawl.debug.%s", key)
	value, _ := g.Cfg().Get(gctx.GetInitCtx(), key, false)
	return value.Bool()
}
