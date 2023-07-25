package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type DisassemblyRouter struct{}

func (d *DisassemblyRouter) InitDisassemblyRouter(param *gin.RouterGroup) {
	disassemblyRouter := param.Group("/disassembly")

	disassemblyRouter.GET("/:disassembly-id", controller.Disassembly.Get)       //获取项目拆解详情
	disassemblyRouter.POST("/tree", controller.Disassembly.Tree)                //获取项目拆解的节点树
	disassemblyRouter.POST("", controller.Disassembly.Create)                   //新增项目拆解
	disassemblyRouter.PATCH("/:disassembly-id", controller.Disassembly.Update)  //修改项目拆解
	disassemblyRouter.DELETE("/:disassembly-id", controller.Disassembly.Delete) //删除项目拆解
	//disassemblyRouter.DELETE("/with-inferiors/:disassembly-id", controller.Disassembly.DeleteWithInferiors) //删除项目拆解（子项一并删除）
	disassemblyRouter.POST("/list", controller.Disassembly.GetList) //获取项目拆解列表
}
