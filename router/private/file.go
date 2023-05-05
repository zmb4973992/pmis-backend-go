package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type FileRouter struct{}

func (f *FileRouter) InitFileRouter(param *gin.RouterGroup) {
	fileRouter := param.Group("/file")

	fileRouter.POST("/upload/single", controller.FileManagement.UploadSingleFile)      //上传单个文件
	fileRouter.POST("/upload/multiple", controller.FileManagement.UploadMultipleFiles) //上传多个文件
	fileRouter.DELETE("/:file-snow-id", controller.FileManagement.DeleteFile)          //删除单个文件
}
