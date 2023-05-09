package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type RequestLogRouter struct{}

func (r *RequestLogRouter) InitRequestLogRouter(param *gin.RouterGroup) {
	requestLogRouter := param.Group("/request-log")

	requestLogRouter.GET("/:request-log-snow-id", controller.RequestLog.Get)       //获取请求日志详情
	requestLogRouter.DELETE("/:request-log-snow-id", controller.RequestLog.Delete) //删除请求日志详情
	requestLogRouter.POST("/list", controller.RequestLog.GetList)                  //获取请求日志列表
}
