package user

import (
	"errors"
	"fmt"
	"gf-app/app/model/user"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
)

const (
	USER_SESSION_MARK = "user_info"
	SUCCEED		=	0
	FAIL		=	1
	ERROR		= 	2
)

type SessionUser struct {
	Passport	string
	Nickname	string
	Avatar		string
	Sex 		int
	Signature 	string
	CreatTime	*gtime.Time
}

// 注册输入参数
type SignUpInput struct {
	Passport string `v:"required#账号不能为空"`
	Password string `v:"required#请输入密码"`
	Nickname string
	Email    string `v:"email|required#请输入正确的邮箱|请输入邮箱"`
}

// 登录输入参数
type LoginInput struct {
	Passport string `v:"required#账号不能为空"`
	Password string `v:"required#请输入密码"`
}

// 用户注册
func SignUp(data *SignUpInput) error {
	// 输入参数检查
	if e := gvalid.CheckStruct(data, nil); e != nil {
		return errors.New(e.FirstString())
	}
	// 昵称为非必须，默认使用账号名称
	if data.Nickname == "" {
		data.Nickname = data.Passport
	}
	// 账号唯一性数据检查
	if !CheckPassport(data.Passport) {
		return errors.New(fmt.Sprintf("账号 %s 已被注册", data.Passport))
	}

	// 昵称唯一性检查
	if !CheckNickName(data.Nickname) {
		return errors.New(fmt.Sprintf("昵称 %s 已经存在", data.Nickname))
	}

	// 邮箱唯一性数据检查
	if !CheckEmail(data.Email) {
		return errors.New(fmt.Sprintf("邮箱 %s 已被注册", data.Email))
	}

	//// 对密码进行加密
	//data.Password, _ = gmd5.Encrypt(data.Password)

	// 将输入参数赋值到数据库实体对象上
	var entity *user.Entity
	if err := gconv.Struct(data, &entity); err != nil {
		return err
	}

	// 记录账号创建/注册时间
	entity.CreateTime = gtime.Now()
	entity.Sex = 2
	entity.Signature = "这人很懒，什么也没有写"
	entity.Avatar = "default.png"
	if _, err := user.Save(entity); err != nil {
		return err
	}
	return nil
}

// 用户登录校验
func Login(data *LoginInput, session *ghttp.Session) (error, int) {
	if e := gvalid.CheckStruct(data, nil); e != nil {
		return errors.New(e.FirstString()), FAIL
	}

	model, err := user.Model.FindOne("passport = ?", data.Passport)
	if err != nil {
		return errors.New("服务器正在出小差哦，请联系管理员吧"), ERROR
	}

	if model == nil {
		return errors.New("账号不存在"), FAIL
	}

	if model.Status == 1 {
		return errors.New("账号已被删除，请联系管理员"), FAIL
	} else if model.Status != 0 {
		return errors.New("账号存在问题，请联系管理员"), FAIL
	}

	//reqPw, err1 := gmd5.Encrypt(data.Password)

	//if err1 != nil {
	//	return errors.New("服务器正在出小差哦，请联系管理员吧"), ERROR
	//}

	if data.Password != model.Password {
		return errors.New("用户名或者密码错误"), FAIL
	}

	avatarUrl := g.Cfg().GetString("avatar.url") + model.Avatar
	sessionData := SessionUser{
		Passport: 	model.Passport,
		Nickname: 	model.Nickname,
		Sex: 		model.Sex,
		Signature:  model.Signature,
		Avatar:		avatarUrl,
		CreatTime: 	model.CreateTime,
	}
	if err2 := session.Set(USER_SESSION_MARK, sessionData); err2 != nil {
		return errors.New("服务器正在出小差哦，请联系管理员吧"), ERROR
	}

	return nil, SUCCEED
}

// 通过Session返回用户名
func GetPassport(session *ghttp.Session) string {
	var u *SessionUser
	session.GetStruct(USER_SESSION_MARK, &u)
	p := u.Passport
	return p
}

// 判断用户是否已经登录
func IsSignedIn(session *ghttp.Session) bool {
	return session.Contains(USER_SESSION_MARK)
}

// 用户注销
func SignOut(session *ghttp.Session) error {
	return session.Remove(USER_SESSION_MARK)
}

// 修改用户头像
func SetAvatar(session *ghttp.Session, img string) error {
	var u *SessionUser
	if err := session.GetStruct(USER_SESSION_MARK, &u); err != nil {
		return err
	}
	model, err1 := user.FindOne("passport = ?", u.Passport)
	if err1 != nil {
		return err1
	}
	model.Avatar = img
	model.Save()
	url := g.Cfg().GetString("avatar.url") + img
	u.Avatar = url
	session.Set(USER_SESSION_MARK, *u)
	return nil
}

// 用户改名接口,改名成功返回true，失败返回false
func SetNickName(newName string, session *ghttp.Session) bool {
	//// 判断新昵称是否存在
	//if !CheckNickName(newName) {
	//	return false
	//}
	data := (*SessionUser)(nil)
	if err := gconv.Struct(session.Get(USER_SESSION_MARK), &data); err!=nil {
		return false
	}
	model, err1 := user.Model.FindOne("passport = ?", data.Passport)
	if err1 != nil {
		return false
	}
	model.Nickname = newName
	model.Save()
	session.Remove(USER_SESSION_MARK)
	data.Nickname = newName
	session.Set(USER_SESSION_MARK, *data)
	return true
}

// 用户改性别接口
func SetSex(sex int, session *ghttp.Session) {
	data := (*SessionUser)(nil)
	if err := gconv.Struct(session.Get(USER_SESSION_MARK), &data); err!=nil {
		return
	}
	model, err1 := user.Model.FindOne("passport = ?", data.Passport)
	if err1 != nil {
		return
	}
	model.Sex = sex
	model.Save()
	session.Remove(USER_SESSION_MARK)
	data.Sex = sex
	session.Set(USER_SESSION_MARK, *data)
}

// 用户改签名
func SetSignature(s string, session *ghttp.Session) {
	data := (*SessionUser)(nil)
	if err := gconv.Struct(session.Get(USER_SESSION_MARK), &data); err!=nil {
		return
	}
	model, err1 := user.Model.FindOne("passport = ?", data.Passport)
	if err1 != nil {
		return
	}
	model.Signature = s
	model.Save()
	session.Remove(USER_SESSION_MARK)
	data.Signature = s
	session.Set(USER_SESSION_MARK, *data)
}

// 用户改密码接口
func SetPassword(p string, session *ghttp.Session) error {
	var u *SessionUser
	if err := session.GetStruct(USER_SESSION_MARK, &u); err != nil {
		return err
	}
	model, err1 := user.FindOne("passport = ?", u.Passport)
	if err1 != nil {
		return err1
	}
	model.Password = p
	model.Save()
	return nil
}

// 检查账号是否符合规范(目前仅检查唯一性),存在返回false,否则true
func CheckPassport(passport string) bool {
	if i, err := user.FindCount("passport", passport); err != nil {
		return false
	} else {
		return i == 0
	}
}

// 检查昵称是否符合规范(目前仅检查唯一性),存在返回false,否则true
func CheckNickName(nickname string) bool {
	if i, err := user.FindCount("nickname", nickname); err != nil {
		return false
	} else {
		return i == 0
	}
}

// 检查email是否唯一
func CheckEmail(email string) bool {
	if i, err := user.FindCount("email", email); err != nil {
		return false
	} else {
		return i == 0
	}
}

// 获得本用户的用户信息详情
func GetProfile(session *ghttp.Session) (u *SessionUser) {
	_ = session.GetStruct(USER_SESSION_MARK, &u)
	return
}

// 获得其他用户的信息
func GetOthersProfile(passport string) (u *SessionUser, err error) {
	model, err := user.Model.FindOne("passport = ?", passport)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, errors.New("用户不存在")
	}

	avatarUrl := g.Cfg().GetString("avatar.url") + model.Avatar
	sessionData := SessionUser{
		Passport: 	model.Passport,
		Nickname: 	model.Nickname,
		Sex: 		model.Sex,
		Signature:  model.Signature,
		Avatar:		avatarUrl,
		CreatTime: 	model.CreateTime,
	}
	u = &sessionData
	return u, nil
}
