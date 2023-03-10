// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CrawlDict is the golang structure of table crawl_dict for DAO operations like Where/Data.
type CrawlDict struct {
	g.Meta     `orm:"table:crawl_dict, do:true"`
	Id         interface{} // ID
	Namespace  interface{} // 命名空间
	DictKey    interface{} // 键名
	DictValue  interface{} // 键值
	DictSort   interface{} //
	DictStatus interface{} // 状态. 0-停用；1-启用
	DictDesc   interface{} // 描述
}
