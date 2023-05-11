package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type RoleRouter struct{}

func (r *RoleRouter) InitRoleRouter(param *gin.RouterGroup) {
	roleRouter := param.Group("/role")

	roleRouter.GET("/:role-id", controller.Role.Get)                //获取角色详情
	roleRouter.POST("", controller.Role.Create)                     //新增角色
	roleRouter.PATCH("/:role-id", controller.Role.Update)           //修改角色
	roleRouter.DELETE("/:role-id", controller.Role.Delete)          //删除角色
	roleRouter.POST("/list", controller.Role.GetList)               //获取角色列表
	roleRouter.POST("/:role-id/users", controller.Role.UpdateUsers) //调整角色的用户

}
