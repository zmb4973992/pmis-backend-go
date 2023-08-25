package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ProjectDailyAndCumulativeIncomeRouter struct{}

func (c *ProjectDailyAndCumulativeIncomeRouter) InitProjectDailyAndCumulativeIncomeRouter(param *gin.RouterGroup) {
	projectDailyAndCumulativeIncomeRouter := param.Group("/project-daily-and-cumulative-income")
	projectDailyAndCumulativeIncomeRouter.PATCH("", controller.ProjectDailyAndCumulativeIncome.Update)      //更新项目当日和累计收款详情
	projectDailyAndCumulativeIncomeRouter.POST("/list", controller.ProjectDailyAndCumulativeIncome.GetList) //获取项目当日和累计收款列表
}
