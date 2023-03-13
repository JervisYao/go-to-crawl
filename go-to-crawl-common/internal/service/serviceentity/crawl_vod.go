package serviceentity

import (
	"go-to-crawl-common/internal/model/entity"
)

type CrawlVodConfigTaskDO struct {
	*entity.CrawlVodConfig
	*entity.CrawlVodConfigTask
}
