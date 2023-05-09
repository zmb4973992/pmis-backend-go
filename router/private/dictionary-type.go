package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
)

type DictionaryTypeRouter struct{}

func (c *DictionaryTypeRouter) InitDictionaryTypeRouter(param *gin.RouterGroup) {
	dictionaryTypeRouter := param.Group("/dictionary-type")

	dictionaryTypeRouter.GET("/:dictionary-type-snow-id", controller.DictionaryType.Get)     //获取字典类型
	dictionaryTypeRouter.POST("", middleware.RequestLog(), controller.DictionaryType.Create) //新增字典类型
	//dictionaryTypeRouter.POST("/batch", middleware.RequestLog(), controller.DictionaryType.CreateInBatches)        //批量新增字典类型
	dictionaryTypeRouter.PATCH("/:dictionary-type-snow-id", middleware.RequestLog(), controller.DictionaryType.Update)  //修改字典类型
	dictionaryTypeRouter.DELETE("/:dictionary-type-snow-id", middleware.RequestLog(), controller.DictionaryType.Delete) //删除字典类型
	//dictionaryTypeRouter.POST("/array", controller.DictionaryType.GetArray)                                               //获取字典类型的数组
	dictionaryTypeRouter.POST("/list", controller.DictionaryType.GetList) //获取字典类型的列表
}
