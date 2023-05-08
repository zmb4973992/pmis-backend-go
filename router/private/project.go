package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ProjectRouter struct{}

func (p *ProjectRouter) InitProjectRouter(param *gin.RouterGroup) {
	projectRouter := param.Group("/project")

	projectRouter.GET("/:project-snow-id", controller.Project.Get)       //获取项目详情
	projectRouter.POST("", controller.Project.Create)                    //新增项目
	projectRouter.PATCH("/:project-snow-id", controller.Project.Update)  //修改项目
	projectRouter.DELETE("/:project-snow-id", controller.Project.Delete) //删除项目
	projectRouter.POST("/list", controller.Project.GetList)              //获取项目列表
}