package user

import (
	"gf-app/app/service/user"
	"gf-app/library/response"
	"github.com/gogf/gf/net/ghttp"
)

func SendMail(r *ghttp.Request)  {
	passport := r.GetString("passport")
	err := user.Sendmail(passport)
	if err != nil {
		response.Error(r, "err")
	}
	response.Succ(r, "验证码已发送")
}

func SetNewPassport(r *ghttp.Request) {
	passport := r.GetString("passport")
	code := r.GetString("code")
	password := r.GetString("password")
	err := user.SetNewPassport(passport, code, password)
	if err != nil {
		response.Error(r, "err")
	}
	response.Succ(r, "修改密码成功")
}