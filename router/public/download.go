package public

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
)

type DownloadRouter struct{}

func (d *DownloadRouter) InitDownloadRouter(param *gin.RouterGroup) {
	downloadRouter := param.Group("")
	downloadRouter.Use(middleware.RateLimit())
	downloadRouter.GET("/download/:file-name", controller.Download)
}
