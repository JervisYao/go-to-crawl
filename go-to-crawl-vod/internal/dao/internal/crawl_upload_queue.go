// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CrawlUploadQueueDao is the data access object for table crawl_upload_queue.
type CrawlUploadQueueDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns CrawlUploadQueueColumns // columns contains all the column names of Table for convenient usage.
}

// CrawlUploadQueueColumns defines and stores column names for table crawl_upload_queue.
type CrawlUploadQueueColumns struct {
	Id           string // 主键ID
	HostIp       string //
	CountryCode  string // 国家二字码.(eg: CN,US,SG等)
	VideoYear    string // 视频发布年份
	VideoCollId  string // 视频集ID（视频集ID，不限于电视剧,-1代表单集视频，或者说电影）
	VideoItemId  string // 视频集对应视频项ID（不限于电视剧的剧集）
	FileName     string // 文件标题
	FileType     string // 文件类型. 1-视频；2-大体积资源；（小文件无需用队列，直接用上传接口）
	FileSize     string // 文件大小. 单位KB
	Msg          string //
	UploadStatus string // 上传状态.0-创建任务;1-上传中;2-上传完成;3-流媒体处理中;4-流媒体处理结束;5-流媒体处理异常
	CreateUser   string // 添加人
	CreateTime   string // 添加时间
	UpdateUser   string // 更新人
	UpdateTime   string // 更新时间
}

// crawlUploadQueueColumns holds the columns for table crawl_upload_queue.
var crawlUploadQueueColumns = CrawlUploadQueueColumns{
	Id:           "id",
	HostIp:       "host_ip",
	CountryCode:  "country_code",
	VideoYear:    "video_year",
	VideoCollId:  "video_coll_id",
	VideoItemId:  "video_item_id",
	FileName:     "file_name",
	FileType:     "file_type",
	FileSize:     "file_size",
	Msg:          "msg",
	UploadStatus: "upload_status",
	CreateUser:   "create_user",
	CreateTime:   "create_time",
	UpdateUser:   "update_user",
	UpdateTime:   "update_time",
}

// NewCrawlUploadQueueDao creates and returns a new DAO object for table data access.
func NewCrawlUploadQueueDao() *CrawlUploadQueueDao {
	return &CrawlUploadQueueDao{
		group:   "default",
		table:   "crawl_upload_queue",
		columns: crawlUploadQueueColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CrawlUploadQueueDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CrawlUploadQueueDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CrawlUploadQueueDao) Columns() CrawlUploadQueueColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CrawlUploadQueueDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CrawlUploadQueueDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CrawlUploadQueueDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
