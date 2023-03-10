// +----------------------------------------------------------------------
// | EasyGoAdmin敏捷开发框架 [ 赋能开发者，助力企业发展 ]
// +----------------------------------------------------------------------
// | 版权所有 2019~2022 深圳EasyGoAdmin研发中心
// +----------------------------------------------------------------------
// | Licensed LGPL-3.0 EasyGoAdmin并不是自由软件，未经许可禁止去掉相关版权
// +----------------------------------------------------------------------
// | 官方网站: http://www.easygoadmin.vip
// +----------------------------------------------------------------------
// | Author: @半城风雨 团队荣誉出品
// +----------------------------------------------------------------------
// | 版权和免责声明:
// | 本团队对该软件框架产品拥有知识产权（包括但不限于商标权、专利权、著作权、商业秘密等）
// | 均受到相关法律法规的保护，任何个人、组织和单位不得在未经本团队书面授权的情况下对所授权
// | 软件框架产品本身申请相关的知识产权，禁止用于任何违法、侵害他人合法权益等恶意的行为，禁
// | 止用于任何违反我国法律法规的一切项目研发，任何个人、组织和单位用于项目研发而产生的任何
// | 意外、疏忽、合约毁坏、诽谤、版权或知识产权侵犯及其造成的损失 (包括但不限于直接、间接、
// | 附带或衍生的损失等)，本团队不承担任何法律责任，本软件框架禁止任何单位和个人、组织用于
// | 任何违法、侵害他人合法利益等恶意的行为，如有发现违规、违法的犯罪行为，本团队将无条件配
// | 合公安机关调查取证同时保留一切以法律手段起诉的权利，本软件框架只能用于公司和个人内部的
// | 法律所允许的合法合规的软件产品研发，详细声明内容请阅读《框架免责声明》附件；
// +----------------------------------------------------------------------

/**
 * 公共函数库
 * @author 半城风雨
 * @since 2021/3/2
 * @File : common
 */
package common

type BunissType int

//业务类型
const (
	BOther BunissType = 0 //0其它
	BAdd   BunissType = 1 //1新增
	BEdit  BunissType = 2 //2修改
	BDel   BunissType = 3 //3删除
)

// 返回结果对象
type JsonResult struct {
	Code  int         `json:"code"`   // 响应编码：0成功 401请登录 403无权限 500错误
	Msg   string      `json:"msg"`    // 消息提示语
	AddID int64       `json:"add_id"` //新增成功时最后一条记录的ID
	Data  interface{} `json:"data"`   // 数据对象
	Count int         `json:"count"`  // 记录总数
	Btype BunissType  `json:"btype"`  // 业务类型
}

// 验证码
type CaptchaRes struct {
	Code  int         `json:"code"`  //响应编码 0 成功 500 错误 403 无权限
	Msg   string      `json:"msg"`   //消息
	Data  interface{} `json:"data"`  //数据内容
	IdKey string      `json:"idkey"` //验证码ID
}

// 部门类型
var DEPT_TYPE_LIST = map[int]string{
	1: "公司",
	2: "子公司",
	3: "部门",
	4: "小组",
}

// 菜单类型
var MENU_TYPE_LIST = map[int]string{
	0: "菜单",
	1: "节点",
}

// 城市等级
var CITY_LEVEL = map[int]string{
	1: "省份",
	2: "城市",
	3: "县区",
	4: "街道",
}

// 配置项类型
var CONFIG_DATA_TYPE_LIST = map[string]string{
	"text":     "单行文本",
	"textarea": "多行文本",
	"ueditor":  "富文本编辑器",
	"date":     "日期",
	"datetime": "时间",
	"number":   "数字",
	"select":   "下拉框",
	"radio":    "单选框",
	"checkbox": "复选框",
	"image":    "单张图片",
	"images":   "多张图片",
	"password": "密码",
	"icon":     "字体图标",
	"file":     "单个文件",
	"files":    "多个文件",
	"hidden":   "隐藏",
	"readonly": "只读文本",
}

// 友链类型
var LINK_TYPE_LIST = map[int]string{
	1: "友情链接",
	2: "合作伙伴",
}

// 友链形式
var LINK_FORM_LIST = map[int]string{
	1: "文字链接",
	2: "图片链接",
}

// 友链平台
var LINK_PLATFORM_LIST = map[int]string{
	1: "PC站",
	2: "WAP站",
	3: "小程序",
	4: "APP应用",
}

// 站点类型
var ITEM_TYPE_LIST = map[int]string{
	1: "国内站点",
	2: "国外站点",
	3: "其他站点",
}

// 广告位所属平台
var ADSORT_PLATFORM_LIST = map[int]string{
	1: "PC站",
	2: "WAP站",
	3: "小程序",
	4: "APP应用",
}

// 广告类型
var AD_TYPE_LIST = map[int]string{
	1: "图片",
	2: "文字",
	3: "视频",
	4: "其他",
}

// 通知来源
var NOTICE_SOURCE_LIST = map[int]string{
	1: "内部通知",
	2: "外部通知",
}

// 会员设备类型
var MEMBER_DEVICE_LIST = map[int]string{
	1: "苹果",
	2: "安卓",
	3: "WAP站",
	4: "PC站",
	5: "后台添加",
}

// 会员来源
var MEMBER_SOURCE_LIST = map[int]string{
	1: "注册会员",
	2: "马甲会员",
}
