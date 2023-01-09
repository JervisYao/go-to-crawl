package crawl

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"go-to-crawl-vod/internal/dao"
	"time"
)

var (
	vt  = dao.CrawlVod.Columns()
	vti = dao.CrawlVodItem.Columns()
	vc  = dao.CrawlVodConfig.Columns()
	vct = dao.CrawlVodConfigTask.Columns()
)

const (
	ConfigTaskStatusInit       = 0
	ConfigTaskStatusProcessing = 1
	ConfigTaskStatusErr        = 2
	ConfigTaskStatusOk         = 3
)

func GetVodConfigById(id int) *model.CmsCrawlVodConfig {
	config, _ := dao.CrawlVodConfig.FindOne(vc.Id, id)
	return config
}

func GetVodConfig() *model.CmsCrawlVodConfig {
	hourBefore := time.Now().Add(-gtime.H).Format(timeUtil.YYYY_MM_DD_HH_MM_SS)
	where := dao.CrawlVodConfig.Where(fmt.Sprintf("%v < '%v' or %v is null", vc.UpdateTime, hourBefore, vc.UpdateTime))
	where := dao.CrawlVodConfig.Ctx(gctx.GetInitCtx()).Where(fmt.Sprintf("%v < '%v' or %v is null", vc.UpdateTime, hourBefore, vc.UpdateTime))
	where = where.And(vc.SeedStatus, 1)
	where.Order(vc.UpdateTime)
	config, _ := where.FindOne()
	return config
}

func UpdateVodConfig(vodConfig *model.CmsCrawlVodConfig) {
	vodConfig.UpdateTime = gtime.Now()
	dao.CrawlVodConfig.Data(vodConfig).Where(vc.Id, vodConfig.Id).Update()
}

func UpdateVodConfigTaskStatus(configTask *model.CmsCrawlVodConfigTask, status int) {
	configTask.TaskStatus = status
	configTask.UpdateTime = gtime.Now()
	dao.CrawlVodConfigTask.Data(configTask).Where(vct.Id, configTask.Id).Update()
}

func GetVodConfigTaskDO() *do.CrawlVodConfigTaskDO {
	configTask, _ := dao.CrawlVodConfigTask.One(vct.TaskStatus, ConfigTaskStatusInit)
	if configTask == nil {
		return nil
	}
	config, _ := dao.CrawlVodConfig.One(vc.Id, configTask.VodConfigId)

	taskDO := new(do.CrawlVodConfigTaskDO)
	taskDO.CmsCrawlVodConfigTask = configTask
	taskDO.CmsCrawlVodConfig = config

	return taskDO
}

func GetVodTvById(id int) *model.CmsCrawlVodTv {
	one, _ := dao.CrawlVodTv.Where(vt.Id, id).FindOne()
	return one
}

func GetVodTvByStatus(crawlStatus int) *model.CmsCrawlVodTv {
	one, _ := dao.CrawlVodTv.Where(vt.CrawlStatus, crawlStatus).FindOne()
	return one
}

func GetVodTvByMd5(vodMd5 string) *model.CmsCrawlVodTv {
	one, _ := dao.CrawlVodTv.Where(vt.VodMd5, vodMd5).FindOne()
	return one
}

func UpdateVodTVStatus(vodTv *model.CmsCrawlVodTv, status int) {
	vodTv.CrawlStatus = status
	vodTv.UpdateTime = gtime.Now()
	dao.CrawlVodTv.Data(vodTv).Where(vt.Id, vodTv.Id).Update()
}

func GetPreparedVodTvItem() *model.CmsCrawlVodTvItem {
	join := g.Model(dao.CrawlVodTvItem.Table+" vti").LeftJoin(dao.CrawlVodTv.Table+" vt", fmt.Sprintf("vti.%s = vt.%s", vti.TvId, vt.Id))
	record, _ := join.Fields("vti.*").One(fmt.Sprintf("vti.%s = %d and vt.%s = %d", vti.CrawlStatus, CrawlTVItemInit, vt.CrawlStatus, CrawlTVPadIdOk))
	if record == nil {
		return nil
	}

	tvItem := new(model.CmsCrawlVodTvItem)
	_ = gconv.Struct(record, tvItem)
	return tvItem
}

func GetVodTvItemByMd5(vodItemMd5 string) *model.CmsCrawlVodTvItem {
	one, _ := dao.CrawlVodTvItem.Where(vti.TvItemMd5, vodItemMd5).FindOne()
	return one
}

func GetVodTvItemByVideoItemId(videoItemId string) *model.CmsCrawlVodTvItem {
	one, _ := dao.CrawlVodTvItem.Where(vti.VideoItemId, videoItemId).FindOne()
	return one
}

func UpdateVodTVItemStatus(vodTvItem *model.CmsCrawlVodTvItem, status int) {
	vodTvItem.CrawlStatus = status
	vodTvItem.UpdateTime = gtime.Now()
	dao.CrawlVodTvItem.Data(vodTvItem).Where(vti.Id, vodTvItem.Id).Update()
}
