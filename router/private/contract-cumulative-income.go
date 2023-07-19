package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ContractCumulativeIncomeRouter struct{}

func (c *ContractCumulativeIncomeRouter) InitContractCumulativeIncomeRouter(param *gin.RouterGroup) {
	contractCumulativeIncomeRouter := param.Group("/contract-cumulative-income")
	contractCumulativeIncomeRouter.PATCH("", controller.ContractCumulativeIncome.Update)      //更新合同累计收款详情
	contractCumulativeIncomeRouter.POST("/list", controller.ContractCumulativeIncome.GetList) //获取合同累计收款列表
}
