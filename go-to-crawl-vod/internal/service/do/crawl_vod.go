package do

import (
	"go-to-crawl-vod/internal/model/entity"
)

type CrawlVodConfigTaskDO struct {
	*entity.CrawlVodConfig
	*entity.CrawlVodConfigTask
}
