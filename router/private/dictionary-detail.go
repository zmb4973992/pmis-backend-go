package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type DictionaryDetailRouter struct{}

func (d *DictionaryDetailRouter) InitDictionaryDetailRouter(param *gin.RouterGroup) {
	dictionaryDetailRouter := param.Group("/dictionary-detail")

	dictionaryDetailRouter.GET("/:dictionary-detail-snow-id", controller.DictionaryDetail.Get)       //获取字典项的值
	dictionaryDetailRouter.POST("", controller.DictionaryDetail.Create)                              //新增字典项的值
	dictionaryDetailRouter.POST("/batch", controller.DictionaryDetail.CreateInBatches)               //批量新增字典项的值
	dictionaryDetailRouter.PATCH("/:dictionary-detail-snow-id", controller.DictionaryDetail.Update)  //修改字典项的值
	dictionaryDetailRouter.DELETE("/:dictionary-detail-snow-id", controller.DictionaryDetail.Delete) //删除字典项的值
	dictionaryDetailRouter.POST("/list", controller.DictionaryDetail.GetList)                        //获取字典项的列表
}
