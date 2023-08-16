package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type MessageRouter struct{}

func (m *MessageRouter) InitMessageRouter(param *gin.RouterGroup) {
	messageRouter := param.Group("/message")

	messageRouter.GET("/:message-id", controller.Message.Get)       //获取消息
	messageRouter.POST("", controller.Message.Create)               //新增消息
	messageRouter.PATCH("/:message-id", controller.Message.Update)  //修改消息
	messageRouter.DELETE("/:message-id", controller.Message.Delete) //删除消息
	messageRouter.POST("/list", controller.Message.GetList)         //获取消息列表
}
