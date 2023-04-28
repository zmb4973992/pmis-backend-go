package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type DepartmentRouter struct{}

func (d *DepartmentRouter) InitDepartmentRouter(param *gin.RouterGroup) {
	departmentRouter := param.Group("/department")

	departmentRouter.GET("/:department-id", controller.Department.Get)       //获取部门详情
	departmentRouter.POST("", controller.Department.Create)                  //新增部门
	departmentRouter.PATCH("/:department-id", controller.Department.Update)  //修改部门
	departmentRouter.DELETE("/:department-id", controller.Department.Delete) //删除部门
	departmentRouter.POST("/array", controller.Department.GetArray)          //获取部门数组
	departmentRouter.POST("/list", controller.Department.GetList)            //获取部门列表
}
