package response

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

const (
	SUCCEED		=	0
	FAIL		=	-1
	ERROR		= 	-2
)

// 数据返回通用JSON数据结构
type JsonResponse struct {
	Code    int         `json:"code"`    // 错误码((0:成功, -1:失败, <-1:错误码))
	Message string      `json:"msg"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据(业务接口定义具体数据结构)
}

// 标准返回结果数据结构封装。
func Json(r *ghttp.Request, code int, message string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(JsonResponse{
		Code:    code,
		Message: message,
		Data:    responseData,
	})
}

// 返回JSON数据并退出当前HTTP执行函数。
func JsonExit(r *ghttp.Request, err int, msg string, data ...interface{}) {
	Json(r, err, msg, data...)
	r.Exit()
}

// 返回失败信息
func Fail(r *ghttp.Request, msg string) {
	JsonExit(r, FAIL, msg)
}

// 返回出错信息
func Error(r *ghttp.Request, msg string) {
	JsonExit(r, ERROR, msg)
}

// 返回成功信息
func Succ(r *ghttp.Request, msg string) {
	JsonExit(r, SUCCEED, msg)
}

// 返回头像修改成功信息
func SetAvatarSucc(r *ghttp.Request, url string) {
	url = g.Cfg().GetString("avatar.url") + url
	JsonExit(r, SUCCEED, "修改头像成功", url)
}