package crawl

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"go-to-crawl-vod/internal/dao"
	"go-to-crawl-vod/internal/model/entity"
	"go-to-crawl-vod/internal/service/do"
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

func GetVodConfigById(id int) *entity.CrawlVodConfig {
	config, _ := dao.CrawlVodConfig.Ctx(gctx.GetInitCtx()).FindOne(vc.Id, id)
	return config
}

func GetVodConfig() *entity.CrawlVodConfig {
	hourBefore := time.Now().Add(-gtime.H).Format(timeUtil.YYYY_MM_DD_HH_MM_SS)
	where := dao.CrawlVodConfig.Ctx(gctx.GetInitCtx()).Where(fmt.Sprintf("%v < '%v' or %v is null", vc.UpdateTime, hourBefore, vc.UpdateTime))
	where = where.Where(vc.SeedStatus, 1)
	where.Order(vc.UpdateTime)
	config, _ := where.FindOne()
	return config
}

func UpdateVodConfig(vodConfig *entity.CrawlVodConfig) {
	vodConfig.UpdateTime = gtime.Now()
	dao.CrawlVodConfig.Ctx(gctx.GetInitCtx()).Data(vodConfig).Where(vc.Id, vodConfig.Id).Update()
}

func UpdateVodConfigTaskStatus(configTask *entity.CrawlVodConfigTask, status int) {
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
	taskDO.CrawlVodConfigTask = configTask
	taskDO.CrawlVodConfig = config

	return taskDO
}

func GetVodTvById(id int) *entity.CrawlVod {
	one, _ := dao.CrawlVod.Ctx(gctx.GetInitCtx()).Where(vt.Id, id).FindOne()
	return one
}

func GetVodTvByStatus(crawlStatus int) *entity.CrawlVod {
	one, _ := dao.CrawlVod.Ctx(gctx.GetInitCtx()).Where(vt.CrawlStatus, crawlStatus).FindOne()
	return one
}

func GetVodTvByMd5(vodMd5 string) *entity.CrawlVod {
	one, _ := dao.CrawlVod.Ctx(gctx.GetInitCtx()).Where(vt.VodMd5, vodMd5).FindOne()
	return one
}

func UpdateVodTVStatus(vodTv *entity.CrawlVod, status int) {
	vodTv.CrawlStatus = status
	vodTv.UpdateTime = gtime.Now()
	dao.CrawlVod.Ctx(gctx.GetInitCtx()).Data(vodTv).Where(vt.Id, vodTv.Id).Update()
}

func GetPreparedVodTvItem() *entity.CrawlVodItem {
	join := g.Model(dao.CrawlVodItem.Table()+" vti").LeftJoin(dao.CrawlVod.Table()+" vt", fmt.Sprintf("vti.%s = vt.%s", vti.TvId, vt.Id))
	record, _ := join.Fields("vti.*").One(fmt.Sprintf("vti.%s = %d and vt.%s = %d", vti.CrawlStatus, CrawlTVItemInit, vt.CrawlStatus, CrawlTVPadIdOk))
	if record == nil {
		return nil
	}

	tvItem := new(entity.CrawlVodItem)
	_ = gconv.Struct(record, tvItem)
	return tvItem
}

func GetVodTvItemByMd5(vodItemMd5 string) *entity.CrawlVodItem {
	one, _ := dao.CrawlVodItem.Ctx(gctx.GetInitCtx()).Where(vti.TvItemMd5, vodItemMd5).FindOne()
	return one
}

func GetVodTvItemByVideoItemId(videoItemId string) *entity.CrawlVodItem {
	one, _ := dao.CrawlVodItem.Ctx(gctx.GetInitCtx()).Where(vti.VideoItemId, videoItemId).FindOne()
	return one
}

func UpdateVodTVItemStatus(vodTvItem *entity.CrawlVodItem, status int) {
	vodTvItem.CrawlStatus = status
	vodTvItem.UpdateTime = gtime.Now()
	dao.CrawlVodItem.Ctx(gctx.GetInitCtx()).Data(vodTvItem).Where(vti.Id, vodTvItem.Id).Update()
}
