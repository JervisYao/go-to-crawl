// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CrawlQueue is the golang structure for table crawl_queue.
type CrawlQueue struct {
	Id              int         `json:"id"              ` // 主键ID
	BusinessType    int         `json:"businessType"    ` // 业务类型。0-普通类型；1-抓付费资源类型
	HostLabel       string      `json:"hostLabel"       ` // 主机标签。任务处理的主机的标签(config.yaml配置)。由哪台机器领取的M3U8下载任务就不能变更了
	CountryCode     string      `json:"countryCode"     ` // 国家二字码.(eg: CN,US,SG等)
	VideoYear       int         `json:"videoYear"       ` // 视频发布年份
	VideoCollId     int64       `json:"videoCollId"     ` // 视频集ID（视频集ID，不限于电视剧,-1代表单集视频，或者说电影）
	VideoItemId     int64       `json:"videoItemId"     ` // 视频集对应视频项ID（不限于电视剧的剧集）
	CrawlType       int         `json:"crawlType"       ` // 抓取类型.1-页面URL;2-文件m3u8;3-MP4地址
	CrawlStatus     int         `json:"crawlStatus"     ` // //抓取状态.0-创建任务;1-M3U8 URL抓取中;2-M3U8 URL抓取失败;3-M3U8 URL抓取完成;4-M3U8下载中;5-M3U8下载异常;6-M3U8下载结束
	CrawlSeedUrl    string      `json:"crawlSeedUrl"    ` // 种子URL
	CrawlSeedParams string      `json:"crawlSeedParams" ` // 种子URL携带的参数。保存Json串
	CrawlM3U8Url    string      `json:"crawlM3U8Url"    ` // m3u8 url
	CrawlM3U8Text   string      `json:"crawlM3U8Text"   ` // M3U8文本
	CrawlM3U8Cnt    int         `json:"crawlM3U8Cnt"    ` // m3u8 url抓取次数
	CrawlM3U8Notify int         `json:"crawlM3U8Notify" ` // crawl_m3u8_cnt次数超过阈值告警,需要人工介入,大概率要优化代码了
	ErrorMsg        string      `json:"errorMsg"        ` // 错误信息
	CreateUser      string      `json:"createUser"      ` // 创建者
	CreateTime      *gtime.Time `json:"createTime"      ` // 创建时间
	UpdateUser      string      `json:"updateUser"      ` // 修改者
	UpdateTime      *gtime.Time `json:"updateTime"      ` // 修改时间
}
