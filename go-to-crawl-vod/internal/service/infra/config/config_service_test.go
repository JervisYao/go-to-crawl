package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
)

func TestGetCrawlHostLabel(t *testing.T) {
	g.Log().Infof(gctx.GetInitCtx(), GetCrawlHostLabel())
}
