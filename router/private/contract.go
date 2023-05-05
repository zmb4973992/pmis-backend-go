package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ContractRouter struct{}

func (c *ContractRouter) InitContractRouter(param *gin.RouterGroup) {
	contractRouter := param.Group("/contract")

	contractRouter.GET("/:contract-snow-id", controller.Contract.Get)       //获取合同详情
	contractRouter.POST("", controller.Contract.Create)                     //新增合同
	contractRouter.PATCH("/:contract-snow-id", controller.Contract.Update)  //修改合同
	contractRouter.DELETE("/:contract-snow-id", controller.Contract.Delete) //删除合同
	contractRouter.POST("/list", controller.Contract.GetList)               //获取合同列表
}
