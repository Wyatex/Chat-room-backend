package router

import (
	"gf-app/app/api/chat"
	"gf-app/app/api/hello"
	"gf-app/app/api/user"
	"gf-app/app/service/middleware"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
	"time"
)

func init() {
	s := g.Server()
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  24 * time.Hour,
		"SessionStorage": gsession.NewStorageRedis(g.Redis()),
	})
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORS)
		ctlUser := new(user.Controller)
		group.ALL("/", hello.Hello)
		group.ALL("/user", ctlUser)
		group.ALL("/user/profile/:passport", user.GetOthersProfile)

		group.ALL("/ws", chat.WebSocket)


	})
	s.SetIndexFolder(true)
	s.AddStaticPath("/resource/avatar", "./avatar")
}
