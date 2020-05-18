package user

import (
	"gf-app/app/service/user"
	"gf-app/library/response"
	"github.com/gogf/gf/net/ghttp"
)

// 用户API管理对象
type Controller struct{}

// 注册请求参数，用于前后端交互参数格式约定
type SignUpRequest struct {
	user.SignUpInput
}

// 登录请求参数，用于前后端交互参数格式约定
type LoginRequest struct {
	user.LoginInput
}

// @summary 用户注册接口
// @tags    用户服务
// @produce json
// @param   passport  formData string  true "用户账号名称"
// @param   password  formData string  true "用户密码"
// @param   email     formData string true "用户邮箱"
// @param   nickname  formData string false "用户昵称"
// @router  /user/register [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) Register(r *ghttp.Request) {
	var data *SignUpRequest
	// 这里没有使用Parse而是仅用GetStruct获取对象，
	// 数据校验交给后续的service层统一处理
	if err := r.GetStruct(&data); err != nil {
		response.Fail(r, err.Error())
	}
	if err := user.SignUp(&data.SignUpInput); err != nil {
		response.Fail(r, err.Error())
	} else {
		response.Succ(r, "注册成功")
	}
}

// 登录请求参数，用于前后端交互参数格式约定
type SignInRequest struct {
	Passport string `v:"required#账号不能为空"`
	Password string `v:"required#密码不能为空"`
}

// @summary 用户登录接口
// @tags    用户服务
// @produce json
// @param   passport formData string true "用户账号"
// @param   password formData string true "用户密码"
// @router  /user/login [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) Login(r *ghttp.Request) {
	var data *LoginRequest
	if err := r.GetStruct(&data); err != nil {
		response.Error(r, err.Error())
	}
	if err1, code := user.Login(&data.LoginInput, r.Session); err1 != nil {
		if code == 1 {
			response.Fail(r, err1.Error())
		} else if code == 2 {
			response.Fail(r, err1.Error())
		}
	} else {
		response.Succ(r, "登录成功")
	}
}

// @summary 判断用户是否已经登录
// @tags    用户服务
// @produce json
// @router  /user/is-login [GET]
// @success 200 {object} response.JsonResponse "执行结果:`true/false`"
func (c *Controller) IsLogin(r *ghttp.Request) {
	response.JsonExit(r, 0, "", user.IsSignedIn(r.Session))
}

// @summary 用户注销/退出接口
// @tags    用户服务
// @produce json
// @router  /user/logout [GET]
// @success 200 {object} response.JsonResponse "执行结果, 1: 未登录"
func (c *Controller) Logout(r *ghttp.Request) {
	if err := user.SignOut(r.Session); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	response.JsonExit(r, 0, "ok")
}

// 账号唯一性检测请求参数，用于前后端交互参数格式约定
type CheckPassportRequest struct {
	Passport string
}

// @summary 检测用户账号接口(唯一性校验)
// @tags    用户服务
// @produce json
// @param   passport query string true "用户账号"
// @router  /user/check-passport [GET]
// @success 200 {object} response.JsonResponse "执行结果:`true/false`"
func (c *Controller) CheckPassport(r *ghttp.Request) {
	var data *CheckPassportRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if data.Passport != "" && !user.CheckPassport(data.Passport) {
		response.JsonExit(r, 0, "账号已经存在", false)
	}
	response.JsonExit(r, 0, "", true)
}

// 账号唯一性检测请求参数，用于前后端交互参数格式约定
type CheckNickNameRequest struct {
	Nickname string
}

// @summary 检测用户昵称接口(唯一性校验)
// @tags    用户服务
// @produce json
// @param   nickname query string true "用户昵称"
// @router  /user/check-nick-name [GET]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) CheckNickName(r *ghttp.Request) {
	var data *CheckNickNameRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if data.Nickname != "" && !user.CheckNickName(data.Nickname) {
		response.JsonExit(r, 0, "昵称已经存在", false)
	}
	response.JsonExit(r, 0, "ok", true)
}

// @summary 获取本用户详情信息
// @tags    用户服务
// @produce json
// @router  /user/profile [GET]
// @success 200 {object} user.Entity "用户信息"
func (c *Controller) Profile(r *ghttp.Request) {
	if b := user.IsSignedIn(r.Session); b == true {
		response.JsonExit(r, 0, "", user.GetProfile(r.Session))
	} else {
		response.Fail(r, "用户未登录")
	}
}

// @summary 修改本用户详情信息
// @tags    用户服务
// @produce json
// @param   nickname query string true "新昵称"
// @param   sex query string true "性别"
// @param   signature query string true "签名"
// @param   passport query string true "签名"
// @router  /user/set-profile [POST]
// @success 200 {object} response.JsonResponse "用户信息"
func (c *Controller) SetProfile(r *ghttp.Request) {
	nickname := r.GetString("nickname")
	sex := r.GetInt("sex")
	signature := r.GetString("signature")
	password := r.GetString("password")
	if err := user.SetNickName(nickname, r.Session); err == false {
		response.Fail(r, "该昵称已被使用")
	}
	if sex==0 || sex==1 || sex==2 {
		user.SetSex(sex, r.Session)
	}
	if signature != "" {
		user.SetSignature(signature, r.Session)
	}
	if password != "" {
		err1 := user.SetPassword(password, r.Session)
		if err1 != nil {
			response.Fail(r, err1.Error())
		}
	}
	response.Succ(r, "修改成功")
}


// 修改昵称
// @summary 用户改名接口
// @tags    用户服务
// @produce json
// @param   nickname  formData string  true "用户账号名称"
// @router  /user/set-nick-name [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) SetNickName(r *ghttp.Request) {
	nickname := r.GetString("nickname")
	if !user.SetNickName(nickname, r.Session) {
		response.Fail(r, "昵称已存在")
	} else {
		response.Succ(r, "改名成功")
	}
}

// 修改头像
// @tags    用户服务
// @produce json
// @router  /user/set-nickname [POST]
// @success 200 {object} response.JsonResponse "执行结果"
func (c *Controller) SetAvatar(r *ghttp.Request) {
	img := r.GetUploadFile("avatar")
	name, err := img.Save("./avatar", true)
	if err != nil {
		response.Fail(r, "请先裁剪完图片再上传")
	}
	if err1 := user.SetAvatar(r.Session, name); err1 != nil {
		response.Fail(r, "请先登录再执行这个操作")
	}
	response.SetAvatarSucc(r, name)
}