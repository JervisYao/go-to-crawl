package crawlservice

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"go-to-crawl-common/utility/timeutil"
	dao2 "go-to-crawl-video/internal/dao"
	entity2 "go-to-crawl-video/internal/model/entity"
	"go-to-crawl-video/internal/service/serviceentity"
	"time"
)

var (
	vt  = dao2.CrawlVod.Columns()
	vti = dao2.CrawlVodItem.Columns()
	vc  = dao2.CrawlVodConfig.Columns()
	vct = dao2.CrawlVodConfigTask.Columns()
)

const (
	ConfigTaskStatusInit       = 0
	ConfigTaskStatusProcessing = 1
	ConfigTaskStatusErr        = 2
	ConfigTaskStatusOk         = 3
)

func GetVodConfigById(id int) *entity2.CrawlVodConfig {
	var config *entity2.CrawlVodConfig
	dao2.CrawlVodConfig.Ctx(gctx.GetInitCtx()).Scan(&config, vc.Id, id)
	return config
}

func GetVodConfig() *entity2.CrawlVodConfig {
	hourBefore := time.Now().Add(-gtime.H).Format(timeutil.YYYY_MM_DD_HH_MM_SS)
	where := dao2.CrawlVodConfig.Ctx(gctx.GetInitCtx()).Where(fmt.Sprintf("%v < '%v' or %v is null", vc.UpdateTime, hourBefore, vc.UpdateTime))
	where = where.Where(vc.SeedStatus, 1)
	where.Order(vc.UpdateTime)

	var config *entity2.CrawlVodConfig
	where.Scan(&config)
	return config
}

func UpdateVodConfig(vodConfig *entity2.CrawlVodConfig) {
	vodConfig.UpdateTime = gtime.Now()
	dao2.CrawlVodConfig.Ctx(gctx.GetInitCtx()).Data(vodConfig).Where(vc.Id, vodConfig.Id).Update()
}

func UpdateVodConfigTaskStatus(configTask *entity2.CrawlVodConfigTask, status int) {
	configTask.TaskStatus = status
	configTask.UpdateTime = gtime.Now()
	dao2.CrawlVodConfigTask.Ctx(gctx.GetInitCtx()).Data(configTask).Where(vct.Id, configTask.Id).Update()
}

func GetVodConfigTaskDO() *serviceentity.CrawlVodConfigTaskDO {
	var configTask *entity2.CrawlVodConfigTask
	dao2.CrawlVodConfigTask.Ctx(gctx.GetInitCtx()).Where(vct.TaskStatus, ConfigTaskStatusInit).Scan(&configTask)
	if configTask == nil {
		return nil
	}
	var config *entity2.CrawlVodConfig
	dao2.CrawlVodConfig.Ctx(gctx.GetInitCtx()).One(vc.Id, configTask.VodConfigId)

	taskDO := new(serviceentity.CrawlVodConfigTaskDO)
	taskDO.CrawlVodConfigTask = configTask
	taskDO.CrawlVodConfig = config

	return taskDO
}

func GetVodTvById(id int) *entity2.CrawlVod {
	var queue *entity2.CrawlVod
	_ = dao2.CrawlVod.Ctx(gctx.GetInitCtx()).Where(vt.Id, id).Scan(&queue)
	return queue
}

func GetVodTvByStatus(crawlStatus int) *entity2.CrawlVod {
	var one *entity2.CrawlVod
	dao2.CrawlVod.Ctx(gctx.GetInitCtx()).Where(vt.CrawlStatus, crawlStatus).Scan(&one)
	return one
}

func GetVodTvByMd5(vodMd5 string) *entity2.CrawlVod {
	var one *entity2.CrawlVod
	dao2.CrawlVod.Ctx(gctx.GetInitCtx()).Where(vt.VodMd5, vodMd5).Scan(&one)
	return one
}

func UpdateVodTVStatus(vodTv *entity2.CrawlVod, status int) {
	vodTv.CrawlStatus = status
	vodTv.UpdateTime = gtime.Now()
	dao2.CrawlVod.Ctx(gctx.GetInitCtx()).Data(vodTv).Where(vt.Id, vodTv.Id).Update()
}

func GetPreparedVodTvItem() *entity2.CrawlVodItem {
	join := g.Model(dao2.CrawlVodItem.Table()+" vti").LeftJoin(dao2.CrawlVod.Table()+" vt", fmt.Sprintf("vti.%s = vt.%s", vti.TvId, vt.Id))
	record, _ := join.Fields("vti.*").One(fmt.Sprintf("vti.%s = %d and vt.%s = %d", vti.CrawlStatus, CrawlTVItemInit, vt.CrawlStatus, CrawlTVPadIdOk))
	if record == nil {
		return nil
	}

	tvItem := new(entity2.CrawlVodItem)
	_ = gconv.Struct(record, tvItem)
	return tvItem
}

func GetVodTvItemByMd5(vodItemMd5 string) *entity2.CrawlVodItem {
	var one *entity2.CrawlVodItem
	dao2.CrawlVodItem.Ctx(gctx.GetInitCtx()).Where(vti.TvItemMd5, vodItemMd5).Scan(&one)
	return one
}

func UpdateVodTVItemStatus(vodTvItem *entity2.CrawlVodItem, status int) {
	vodTvItem.CrawlStatus = status
	vodTvItem.UpdateTime = gtime.Now()
	dao2.CrawlVodItem.Ctx(gctx.GetInitCtx()).Data(vodTvItem).Where(vti.Id, vodTvItem.Id).Update()
}
