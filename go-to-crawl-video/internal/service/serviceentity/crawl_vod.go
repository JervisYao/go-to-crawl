package serviceentity

import (
	"go-to-crawl-video/internal/model/entity"
)

type CrawlVodConfigTaskDO struct {
	*entity.CrawlVodConfig
	*entity.CrawlVodConfigTask
}
