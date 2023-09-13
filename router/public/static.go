package public

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/global"
	"pmis-backend-go/middleware"
)

type StaticRouter struct{}

func (s *StaticRouter) InitStaticRouter(param *gin.RouterGroup) {
	staticRouter := param.Group("")
	staticRouter.Use(middleware.RateLimit())
	//将存储路径下的所有文件设置为静态文件，可以直接访问
	staticRouter.Static("/static-test", global.Config.Upload.StoragePath)
}
