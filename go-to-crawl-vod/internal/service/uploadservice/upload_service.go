package uploadservice

import (
	"github.com/gogf/gf/v2/os/gctx"
	"go-to-crawl-vod/internal/dao"
	"go-to-crawl-vod/internal/model/entity"
)

var (
	columns = dao.CrawlUploadQueue.Columns()
)

func GetByVideoItemId(videoItemId int64, status int) *entity.CrawlUploadQueue {
	where := dao.CrawlUploadQueue.Ctx(gctx.GetInitCtx()).Where(columns.VideoItemId, videoItemId)
	where = where.Where(columns.UploadStatus, status)

	var do *entity.CrawlUploadQueue
	where.Scan(&do)
	return do
}

func GetById(id int64) *entity.CrawlUploadQueue {
	where := dao.CrawlUploadQueue.Ctx(gctx.GetInitCtx()).Where(columns.Id, id)

	var do *entity.CrawlUploadQueue
	where.Scan(&do)
	return do
}

func UpdateById(queue *entity.CrawlUploadQueue, status int) {
	queue.UploadStatus = status
	where := dao.CrawlUploadQueue.Ctx(gctx.GetInitCtx()).Data(queue).Where(columns.Id, queue.Id)
	where.Update()
}
