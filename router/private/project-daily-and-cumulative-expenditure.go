package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ProjectDailyAndCumulativeExpenditureRouter struct{}

func (c *ProjectDailyAndCumulativeExpenditureRouter) InitProjectDailyAndCumulativeExpenditureRouter(param *gin.RouterGroup) {
	projectDailyAndCumulativeExpenditureRouter := param.Group("/project-daily-and-cumulative-expenditure")
	projectDailyAndCumulativeExpenditureRouter.PATCH("", controller.ProjectDailyAndCumulativeExpenditure.Update)      //更新项目当日和累计付款详情
	projectDailyAndCumulativeExpenditureRouter.POST("/list", controller.ProjectDailyAndCumulativeExpenditure.GetList) //获取项目当日和累计付款列表
}
