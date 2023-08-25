package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ContractDailyAndCumulativeIncomeRouter struct{}

func (c *ContractDailyAndCumulativeIncomeRouter) InitContractDailyAndCumulativeIncomeRouter(param *gin.RouterGroup) {
	contractDailyAndCumulativeIncomeRouter := param.Group("/contract-daily-and-cumulative-income")
	contractDailyAndCumulativeIncomeRouter.PATCH("", controller.ContractDailyAndCumulativeIncome.Update)      //更新合同累计收款详情
	contractDailyAndCumulativeIncomeRouter.POST("/list", controller.ContractDailyAndCumulativeIncome.GetList) //获取合同累计收款列表
}
