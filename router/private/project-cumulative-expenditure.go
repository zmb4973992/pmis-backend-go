package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ProjectCumulativeExpenditureRouter struct{}

func (c *ProjectCumulativeExpenditureRouter) InitProjectCumulativeExpenditureRouter(param *gin.RouterGroup) {
	projectCumulativeExpenditureRouter := param.Group("/project-cumulative-expenditure")
	projectCumulativeExpenditureRouter.PATCH("", controller.ProjectCumulativeExpenditure.Update)      //更新项目累计支出详情
	projectCumulativeExpenditureRouter.POST("/list", controller.ProjectCumulativeExpenditure.GetList) //获取项目累计支出列表
}
