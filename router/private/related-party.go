package private

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
)

type RelatedPartyRouter struct{}

func (r *RelatedPartyRouter) InitRelatedPartyRouter(param *gin.RouterGroup) {
	relatedPartyRouter := param.Group("/related-party")

	relatedPartyRouter.GET("/:related-party-snow-id", controller.RelatedParty.Get)       //获取相关方详情
	relatedPartyRouter.PATCH("/:related-party-snow-id", controller.RelatedParty.Update)  //修改相关方
	relatedPartyRouter.POST("", controller.RelatedParty.Create)                          //新增相关方
	relatedPartyRouter.DELETE("/:related-party-snow-id", controller.RelatedParty.Delete) //删除相关方
	relatedPartyRouter.POST("/list", controller.RelatedParty.GetList)                    //获取相关方列表
}
