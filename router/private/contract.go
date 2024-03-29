package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type ContractRouter struct{}

func (c *ContractRouter) InitContractRouter(param *gin.RouterGroup) {
	contractRouter := param.Group("/contract")

	contractRouter.GET("/:contract-id", controller.Contract.Get)       //获取合同详情
	contractRouter.POST("", controller.Contract.Create)                //新增合同
	contractRouter.PATCH("/:contract-id", controller.Contract.Update)  //修改合同
	contractRouter.DELETE("/:contract-id", controller.Contract.Delete) //删除合同
	contractRouter.POST("/list", controller.Contract.GetList)          //获取合同列表
	contractRouter.POST("/count", controller.Contract.GetCount)        //获取合同数量
}
