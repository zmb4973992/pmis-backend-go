package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
)

type UserRouter struct{}

func (u *UserRouter) InitUserRouter(param *gin.RouterGroup) {
	userRouter := param.Group("/user")

	userRouter.GET("", controller.User.GetByToken)                                //根据header里的token获取用户详情
	userRouter.GET("/:user-id", controller.User.Get)                              //获取用户详情
	userRouter.PATCH("/:user-id", controller.User.Update)                         //修改用户（目前为全功能，考虑改成：修改用户基本信息）
	userRouter.DELETE("/:user-id", middleware.NeedAuth(), controller.User.Delete) //删除用户
	userRouter.POST("/list", controller.User.List)                                //获取用户列表
}
