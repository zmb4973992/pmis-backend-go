package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type OrganizationRouter struct{}

func (o *OrganizationRouter) InitOrganizationRouter(param *gin.RouterGroup) {
	organizationRouter := param.Group("/organization")

	organizationRouter.GET("/:organization-snow-id", controller.Organization.Get)       //获取部门详情
	organizationRouter.POST("", controller.Organization.Create)                         //新增部门
	organizationRouter.PATCH("/:organization-snow-id", controller.Organization.Update)  //修改部门
	organizationRouter.DELETE("/:organization-snow-id", controller.Organization.Delete) //删除部门
	//organizationRouter.POST("/array", controller.Organization.GetArray)            //获取部门数组
	organizationRouter.POST("/list", controller.Organization.GetList) //获取部门列表
}
