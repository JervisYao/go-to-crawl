// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CrawlVodItem is the golang structure for table crawl_vod_item.
type CrawlVodItem struct {
	Id          int         `json:"id"          ` //
	TvId        int         `json:"tvId"        ` // VOD表ID
	TvItemMd5   string      `json:"tvItemMd5"   ` // 集数MD5
	CrawlStatus int         `json:"crawlStatus" ` // 抓取状态.0-INIT;1-自动补全视频信息中;2-补充视频信息失败;3-补充视频信息成功;4-补充TV ID信息中;5-补充TV ID信息失败;6-补充TV ID信息成功
	SeedUrl     string      `json:"seedUrl"     ` // 种子URL
	SeedParams  string      `json:"seedParams"  ` // 种子URL参数
	ErrorMsg    string      `json:"errorMsg"    ` // 错误信息
	Episodes    string      `json:"episodes"    ` // 集数
	CreateUser  string      `json:"createUser"  ` // 创建者
	CreateTime  *gtime.Time `json:"createTime"  ` // 创建时间
	UpdateUser  string      `json:"updateUser"  ` // 修改者
	UpdateTime  *gtime.Time `json:"updateTime"  ` // 修改时间
}