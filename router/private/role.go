package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type RoleRouter struct{}

func (r *RoleRouter) InitRoleRouter(param *gin.RouterGroup) {
	roleRouter := param.Group("/role")

	roleRouter.GET("/:role-id", controller.Contract.Get)       //获取角色详情
	roleRouter.POST("", controller.Contract.Create)            //新增角色
	roleRouter.PATCH("/:role-id", controller.Contract.Update)  //修改角色
	roleRouter.DELETE("/:role-id", controller.Contract.Delete) //删除角色
	roleRouter.POST("/list", controller.Contract.GetList)      //获取角色列表
}
