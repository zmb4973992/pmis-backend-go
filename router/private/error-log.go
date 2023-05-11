package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ErrorLogRouter struct{}

func (e *ErrorLogRouter) InitErrorLogRouter(param *gin.RouterGroup) {
	errorLogRouter := param.Group("/error-log")

	errorLogRouter.GET("/:error-log-id", controller.ErrorLog.Get)       //获取错误日志详情
	errorLogRouter.POST("", controller.ErrorLog.Create)                 //新增错误日志
	errorLogRouter.PATCH("/:error-log-id", controller.ErrorLog.Update)  //修改错误日志
	errorLogRouter.DELETE("/:error-log-id", controller.ErrorLog.Delete) //删除错误日志
	errorLogRouter.POST("/list", controller.ErrorLog.GetList)           //获取错误日志列表
}
