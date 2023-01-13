// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CrawlUploadQueue is the golang structure of table crawl_upload_queue for DAO operations like Where/Data.
type CrawlUploadQueue struct {
	g.Meta       `orm:"table:crawl_upload_queue, do:true"`
	Id           interface{} // 主键ID
	HostIp       interface{} //
	CountryCode  interface{} // 国家二字码.(eg: CN,US,SG等)
	VideoYear    interface{} // 视频发布年份
	VideoCollId  interface{} // 视频集ID（视频集ID，不限于电视剧,-1代表单集视频，或者说电影）
	VideoItemId  interface{} // 视频集对应视频项ID（不限于电视剧的剧集）
	FileName     interface{} // 文件标题
	FileType     interface{} // 文件类型. 1-视频；2-大体积资源；（小文件无需用队列，直接用上传接口）
	FileSize     interface{} // 文件大小. 单位KB
	Msg          interface{} //
	UploadStatus interface{} // 上传状态.0-创建任务;1-上传中;2-上传完成;3-流媒体处理中;4-流媒体处理结束;5-流媒体处理异常
	CreateUser   interface{} // 添加人
	CreateTime   *gtime.Time // 添加时间
	UpdateUser   interface{} // 更新人
	UpdateTime   *gtime.Time // 更新时间
}
