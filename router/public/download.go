package public

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type DownLoadRouter struct{}

func (d *DownLoadRouter) InitDownLoadRouter(param *gin.RouterGroup) {
	downloadRouter := param.Group("/download")

	downloadRouter.GET("/:file-id", controller.File.Get) //下载文件
}
