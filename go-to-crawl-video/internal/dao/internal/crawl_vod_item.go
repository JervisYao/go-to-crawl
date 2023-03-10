// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CrawlVodItemDao is the data access object for table crawl_vod_item.
type CrawlVodItemDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns CrawlVodItemColumns // columns contains all the column names of Table for convenient usage.
}

// CrawlVodItemColumns defines and stores column names for table crawl_vod_item.
type CrawlVodItemColumns struct {
	Id          string //
	TvId        string // VOD表ID
	TvItemMd5   string // 集数MD5
	CrawlStatus string // 抓取状态.0-INIT;1-自动补全视频信息中;2-补充视频信息失败;3-补充视频信息成功;4-补充TV ID信息中;5-补充TV ID信息失败;6-补充TV ID信息成功
	SeedUrl     string // 种子URL
	SeedParams  string // 种子URL参数
	ErrorMsg    string // 错误信息
	Episodes    string // 集数
	CreateUser  string // 创建者
	CreateTime  string // 创建时间
	UpdateUser  string // 修改者
	UpdateTime  string // 修改时间
}

// crawlVodItemColumns holds the columns for table crawl_vod_item.
var crawlVodItemColumns = CrawlVodItemColumns{
	Id:          "id",
	TvId:        "tv_id",
	TvItemMd5:   "tv_item_md5",
	CrawlStatus: "crawl_status",
	SeedUrl:     "seed_url",
	SeedParams:  "seed_params",
	ErrorMsg:    "error_msg",
	Episodes:    "episodes",
	CreateUser:  "create_user",
	CreateTime:  "create_time",
	UpdateUser:  "update_user",
	UpdateTime:  "update_time",
}

// NewCrawlVodItemDao creates and returns a new DAO object for table data access.
func NewCrawlVodItemDao() *CrawlVodItemDao {
	return &CrawlVodItemDao{
		group:   "default",
		table:   "crawl_vod_item",
		columns: crawlVodItemColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CrawlVodItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CrawlVodItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CrawlVodItemDao) Columns() CrawlVodItemColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CrawlVodItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CrawlVodItemDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CrawlVodItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
