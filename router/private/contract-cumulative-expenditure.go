package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ContractCumulativeExpenditureRouter struct{}

func (c *ContractCumulativeExpenditureRouter) InitContractCumulativeExpenditureRouter(param *gin.RouterGroup) {
	contractCumulativeExpenditureRouter := param.Group("/contract-cumulative-expenditure")
	contractCumulativeExpenditureRouter.PATCH("", controller.ContractCumulativeExpenditure.Update)      //更新合同累计支出详情
	contractCumulativeExpenditureRouter.POST("/list", controller.ContractCumulativeExpenditure.GetList) //获取合同累计支出列表
}
