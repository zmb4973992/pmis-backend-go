package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type UserAndRoleRouter struct{}

func (u *UserAndRoleRouter) InitUserAndRoleRouter(param *gin.RouterGroup) {
	userAndRoleRouter := param.Group("/user-and-role")

	userAndRoleRouter.POST("/role/:role-id", controller.UserAndRole.UpdateByRoleID) //根据roleID修改中间表
	userAndRoleRouter.POST("/user/:user-id", controller.UserAndRole.UpdateByUserID) //根据userID修改中间表
}
