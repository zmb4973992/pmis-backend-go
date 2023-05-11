package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ProgressRouter struct{}

func (p *ProgressRouter) InitProgressRouter(param *gin.RouterGroup) {
	progressRouter := param.Group("/progress")

	progressRouter.GET("/:progress-id", controller.Progress.Get)       //获取进度详情
	progressRouter.POST("", controller.Progress.Create)                //新增进度详情
	progressRouter.PATCH("/:progress-id", controller.Progress.Update)  //修改进度详情
	progressRouter.DELETE("/:progress-id", controller.Progress.Delete) //删除进度详情
	progressRouter.POST("/list", controller.Progress.GetList)          //获取进度列表
}
