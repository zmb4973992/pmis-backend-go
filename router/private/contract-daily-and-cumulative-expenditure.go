package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ContractDailyAndCumulativeExpenditureRouter struct{}

func (c *ContractDailyAndCumulativeExpenditureRouter) InitContractDailyAndCumulativeExpenditureRouter(param *gin.RouterGroup) {
	contractDailyAndCumulativeExpenditureRouter := param.Group("/contract-daily-and-cumulative-expenditure")
	contractDailyAndCumulativeExpenditureRouter.PATCH("", controller.ContractDailyAndCumulativeExpenditure.Update)      //更新合同累计支出详情
	contractDailyAndCumulativeExpenditureRouter.POST("/list", controller.ContractDailyAndCumulativeExpenditure.GetList) //获取合同累计支出列表
}
