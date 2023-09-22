package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ProjectRouter struct{}

func (p *ProjectRouter) InitProjectRouter(param *gin.RouterGroup) {
	projectRouter := param.Group("/project")

	projectRouter.GET("/:project-id", controller.Project.Get)                    //获取项目详情
	projectRouter.POST("", controller.Project.Create)                            //新增项目
	projectRouter.PATCH("/:project-id", controller.Project.Update)               //修改项目
	projectRouter.DELETE("/:project-id", controller.Project.Delete)              //删除项目
	projectRouter.POST("/list", controller.Project.GetList)                      //获取项目列表
	projectRouter.POST("/simplified-list", controller.Project.GetSimplifiedList) //获取简化的项目列表
	projectRouter.POST("/count", controller.Project.GetCount)                    //获取项目数量
}
