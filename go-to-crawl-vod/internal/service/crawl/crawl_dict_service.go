package crawl

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"go-to-crawl-vod/internal/dao"
)

const (
	NsProxy       = "crawlerProxy"
	NsCountryCode = "countryCode"
)

const (
	DictEnable  = 1
	DictDisAble = 0
)

var (
	cdc = dao.CrawlDict.Columns()
)

func GetRandomProxyUrl() string {
	record, _ := dao.CrawlDict.Ctx(gctx.GetInitCtx()).Where(cdc.Namespace, NsProxy).And(cdc.DictStatus, DictEnable).OrderRandom().One()
	if record == nil {
		return ""
	}

	dict := new(model.CmsCrawlDict)
	_ = gconv.Struct(record, dict)

	return fmt.Sprintf("http://%s", dict.DictValue)
}
