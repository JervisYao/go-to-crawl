package httputil

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// 返回结果对象
type JsonResult struct {
	Code  int         `json:"code"`   // 响应编码：0成功 401请登录 403无权限 500错误
	Msg   string      `json:"msg"`    // 消息提示语
	AddID int64       `json:"add_id"` //新增成功时最后一条记录的ID
	Data  interface{} `json:"data"`   // 数据对象
	Count int         `json:"count"`  // 记录总数
}

func ParseParam(req *ghttp.Request, dto interface{}) {
	err := req.Parse(dto)
	if err != nil {
		Error(req, err.Error())
	}
}

func Error(r *ghttp.Request, msg string) {
	r.Response.WriteJsonExit(JsonResult{
		Code: -1,
		Msg:  msg,
	})
}

func Success(r *ghttp.Request) {
	SuccessData(r, nil)
}

func SuccessData(r *ghttp.Request, Data interface{}) {
	SuccessMsgData(r, "SUCCESS", Data)
}

func SuccessMsgData(r *ghttp.Request, msg string, Data interface{}) {
	r.Response.WriteJsonExit(JsonResult{
		Code: 0,
		Msg:  msg,
		Data: Data,
	})
}
