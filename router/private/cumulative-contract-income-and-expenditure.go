package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type CumulativeContractIncomeAndExpenditureRouter struct{}

func (c *CumulativeContractIncomeAndExpenditureRouter) InitCumulativeContractIncomeAndExpenditureRouter(param *gin.RouterGroup) {
	cumulativeContractIncomeAndExpenditureRouter := param.Group("/cumulative-contract-income-and-expenditure")
	cumulativeContractIncomeAndExpenditureRouter.PATCH("", controller.CumulativeContractIncomeAndExpenditure.Update)      //更新项目累计收入和支出详情
	cumulativeContractIncomeAndExpenditureRouter.POST("/list", controller.CumulativeContractIncomeAndExpenditure.GetList) //获取项目累计收入和支出列表
}
