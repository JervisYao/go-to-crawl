package crawlservice

import (
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

func TestGetPreparedVodTvItem(t *testing.T) {
	item := GetPreparedVodTvItem()
	g.Dump(item)
}
