// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-to-crawl-video/internal/dao/internal"
)

// internalCrawlVodItemDao is internal type for wrapping internal DAO implements.
type internalCrawlVodItemDao = *internal.CrawlVodItemDao

// crawlVodItemDao is the data access object for table crawl_vod_item.
// You can define custom methods on it to extend its functionality as you wish.
type crawlVodItemDao struct {
	internalCrawlVodItemDao
}

var (
	// CrawlVodItem is globally public accessible object for table crawl_vod_item operations.
	CrawlVodItem = crawlVodItemDao{
		internal.NewCrawlVodItemDao(),
	}
)

// Fill with you ideas below.
