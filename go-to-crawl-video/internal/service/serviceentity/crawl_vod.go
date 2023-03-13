package serviceentity

import (
	entity2 "go-to-crawl-video/internal/model/entity"
)

type CrawlVodConfigTaskDO struct {
	*entity2.CrawlVodConfig
	*entity2.CrawlVodConfigTask
}
