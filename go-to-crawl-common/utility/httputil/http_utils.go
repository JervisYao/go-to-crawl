package httputil

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"go-to-crawl-common/app/utils/common"
	"go-to-crawl-common/entity/reqobj"
)

func ParseParam(req *ghttp.Request, dto interface{}) {
	err := req.Parse(dto)
	if err != nil {
		Error(req, err.Error())
	}
}

func ParsePageParam(r *ghttp.Request, dto interface{}) {
	ParseParam(r, dto)
	parser, ok := dto.(reqobj.PageParam)
	if ok {
		parser.InitPageParam()
	}
}

func Error(r *ghttp.Request, msg string) {
	_ = r.Response.WriteJsonExit(common.JsonResult{
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
	_ = r.Response.WriteJsonExit(common.JsonResult{
		Code: 0,
		Msg:  msg,
		Data: Data,
	})
}
