package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type IncomeAndExpenditureRouter struct{}

func (c *IncomeAndExpenditureRouter) InitIncomeAndExpenditureRouter(param *gin.RouterGroup) {
	incomeAndExpenditureRouter := param.Group("/income-and-expenditure")

	incomeAndExpenditureRouter.GET("/:income-and-expenditure-id", controller.IncomeAndExpenditure.Get)         //获取收入和支出详情
	incomeAndExpenditureRouter.POST("", controller.IncomeAndExpenditure.Create)                                //新增收入和支出详情
	incomeAndExpenditureRouter.PATCH("/:income-and-expenditure-id", controller.IncomeAndExpenditure.Update)    //修改收入和支出详情
	incomeAndExpenditureRouter.DELETE("/:income-and-expenditure-id", controller.IncomeAndExpenditure.Delete)   //删除收入和支出详情
	incomeAndExpenditureRouter.POST("/list", controller.IncomeAndExpenditure.GetList)                          //获取收入和支出列表
	incomeAndExpenditureRouter.POST("/total-amount", controller.IncomeAndExpenditure.GetCumulativeTotalAmount) //获取累计总额
}
