package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type MenuRouter struct{}

func (m *MenuRouter) InitMenuRouter(param *gin.RouterGroup) {
	menuRouter := param.Group("/menu")

	menuRouter.GET("/:menu-id", controller.Menu.Get)               //获取菜单详情
	menuRouter.POST("", controller.Menu.Create)                    //新增菜单
	menuRouter.PATCH("/:menu-id", controller.Menu.Update)          //修改菜单
	menuRouter.DELETE("/:menu-id", controller.Menu.Delete)         //删除菜单
	menuRouter.POST("/list", controller.Menu.GetList)              //获取菜单列表
	menuRouter.POST("/tree", controller.Menu.GetTree)              //获取菜单树
	menuRouter.POST("/:menu-id/apis", controller.Menu.UpdateUsers) //调整菜单的api
}
