package boot

import (
	"github.com/gogf/gf-swagger/swagger"
	"github.com/gogf/gf/frame/g"
)

func init() {
	s := g.Server()
	s.Plugin(&swagger.Swagger{})


}

