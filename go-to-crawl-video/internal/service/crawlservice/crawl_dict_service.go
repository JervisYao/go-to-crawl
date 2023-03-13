package crawlservice

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"go-to-crawl-video/internal/dao"
	"go-to-crawl-video/internal/model/entity"
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
	record, _ := dao.CrawlDict.Ctx(gctx.GetInitCtx()).Where(cdc.Namespace, NsProxy).Where(cdc.DictStatus, DictEnable).OrderRandom().One()
	if record == nil {
		return ""
	}

	dict := new(entity.CrawlDict)
	_ = gconv.Struct(record, dict)

	return fmt.Sprintf("http://%s", dict.DictValue)
}
