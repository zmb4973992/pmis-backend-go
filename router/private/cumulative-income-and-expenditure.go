package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type CumulativeIncomeAndExpenditureRouter struct{}

func (c *CumulativeIncomeAndExpenditureRouter) InitCumulativeIncomeAndExpenditureRouter(param *gin.RouterGroup) {
	incomeAndExpenditureRouter := param.Group("/cumulative-income-and-expenditure")

	incomeAndExpenditureRouter.PATCH("", controller.CumulativeIncomeAndExpenditure.Update)      //更新累计收入和支出详情
	incomeAndExpenditureRouter.POST("/list", controller.CumulativeIncomeAndExpenditure.GetList) //获取累计收入和支出列表
}
