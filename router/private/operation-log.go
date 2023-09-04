package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type OperationLogRouter struct{}

func (o *OperationLogRouter) InitOperationLogRouter(param *gin.RouterGroup) {
	operationLogRouter := param.Group("/operation-log")

	operationLogRouter.POST("/list", controller.OperationLog.GetList) //获取操作日志列表
}
