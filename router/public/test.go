package public

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
)

type TestRouter struct{}

func (t *TestRouter) InitTestRouter(param *gin.RouterGroup) {
	testRouter := param.Group("")
	testRouter.Use(middleware.RateLimit())
	//将存储路径下的所有文件设置为静态文件，可以直接访问
	testRouter.GET("/test", controller.Test)
}
