package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type FileRouter struct{}

func (f *FileRouter) InitFileRouter(param *gin.RouterGroup) {
	fileRouter := param.Group("/file")

	fileRouter.POST("", controller.File.Create)            //上传单个文件
	fileRouter.DELETE("/:file-id", controller.File.Delete) //删除单个文件
}
