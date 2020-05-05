package chat

import (
	"fmt"
	"gf-app/app/service/user"
	"gf-app/library/response"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
	"time"
)

const SendInterval = 2 * time.Second

var (
	// 在线的用户连接
	conn = gmap.New(true)
	// 在线用户的集合
	users = gset.NewStrSet(true)
	// 使用特定的缓存对象，不使用全局缓存对象
	cache = gcache.New()
)

// Msg 消息结构体
type Msg struct {
	Type string      `json:"type" v:"type@required#消息类型不能为空"`
	Data interface{} `json:"data" v:""`
	From string      `json:"name" v:""`
}

func WebSocket(r *ghttp.Request) {
	msg := &Msg{}

	// 初始化WebSocket请求
	ws, err := r.WebSocket()
	if err != nil {
		g.Log().Error()
		return
	}

	// 拿到登陆的用户名
	u := user.GetPassport(r.Session)
	if u == "" {
		response.Fail(r, "用户还未登录")
	}

	users.Add(u)
	conn.Set(ws, u)

	// 初始化后向所有客户端发送上线消息
	writeUsers()

	for {
		// 阻塞读取WS数据
		_, msgByte, err := ws.ReadMessage()
		if err != nil {
			// 如果失败，那么表示断开，这里清除用户信息
			// 为简单不实现失败重连机制
			users.Remove(u)
			conn.Remove(ws)
			// 通知所有客户端当前用户已下线

			break
		}
		// json参数解析
		if err := gjson.DecodeTo(msgByte, msg); err != nil {
			write(Msg{"error", "消息格式不正确"+err.Error(), ""})
			continue
		}
		// 数据校验
		if e := gvalid.CheckStruct(msg, nil); e != nil {
			write(Msg{"error", e.String(), ""})
		}
		msg.From = u

		// WS操作类型
		switch msg.Type {
		case "send":
			// 发送间隔检查
			intervalKey := fmt.Sprintf("%p", ws)
			if !cache.SetIfNotExist(intervalKey, struct{}{}, SendInterval) {
				write(Msg{"error", "您发送消息过于频繁，请休息下再重试", ""})
				continue
			}
			// 有消息时，群发消息
			if msg.Data != nil {
				if err = writeGroup(Msg{"send", ghtml.SpecialChars(gconv.String(msg.Data)), ghtml.SpecialChars(msg.From)}); err != nil {
					g.Log().Error(err)
				}
			}
		}
	}
}

// 向客户端写入信息
func write(msg Msg) error {
	b, err := gjson.Encode(msg)
	if err != nil {
		return err
	}
	conn.RLockFunc(func(m map[interface{}]interface{}) {
		for user := range m {
			user.(*ghttp.WebSocket).WriteMessage(ghttp.WS_MSG_TEXT, []byte(b))
		}
	})

	return nil
}

// 向所有客户端群发消息
func writeGroup(msg Msg) error {
	b, err := gjson.Encode(msg)
	if err != nil {
		return err
	}
	conn.RLockFunc(func(m map[interface{}]interface{}) {
		for user := range m {
			user.(*ghttp.WebSocket).WriteMessage(ghttp.WS_MSG_TEXT, []byte(b))
		}
	})

	return nil
}


// 向客户端返回用户列表。
// 内部方法不会自动注册到路由中。
func writeUsers() error {
	array := garray.NewSortedStrArray()
	users.Iterator(func(v string) bool {
		array.Add(v)
		return true
	})
	if err := writeGroup(Msg{"list", array.Slice(), ""}); err != nil {
		return err
	}
	return nil
}
