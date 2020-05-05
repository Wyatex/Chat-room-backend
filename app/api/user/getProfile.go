package user

import (
	"gf-app/app/service/user"
	"gf-app/library/response"
	"github.com/gogf/gf/net/ghttp"
)

// @summary 获取其他用户详情信息
// @tags    用户服务
// @produce json
// @router  /user/profile/:passport [GET]
// @success 200 {object} user.Entity "用户信息"
func GetOthersProfile(r *ghttp.Request) {
	passport := r.GetString("passport")
	if p, err := user.GetOthersProfile(passport); err != nil {
		response.Fail(r, err.Error())
	} else {
		response.JsonExit(r, 0, "", p)
	}
}