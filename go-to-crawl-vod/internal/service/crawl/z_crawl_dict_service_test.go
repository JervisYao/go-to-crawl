package crawl

import (
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

func TestGetRandomProxyUrl(t *testing.T) {
	g.Dump(GetRandomProxyUrl())
}

func TestGetVodConfigTaskDO(t *testing.T) {
	g.Dump(GetVodConfigTaskDO())
}
