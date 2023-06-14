package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ProjectCumulativeIncomeRouter struct{}

func (c *ProjectCumulativeIncomeRouter) InitProjectCumulativeIncomeRouter(param *gin.RouterGroup) {
	projectCumulativeIncomeRouter := param.Group("/project-cumulative-income")
	projectCumulativeIncomeRouter.PATCH("", controller.ProjectCumulativeIncome.Update)      //更新项目累计收入详情
	projectCumulativeIncomeRouter.POST("/list", controller.ProjectCumulativeIncome.GetList) //获取项目累计收入列表
}
