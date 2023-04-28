package public

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
)

type ValidateTokenRouter struct{}

func (l *ValidateTokenRouter) InitValidateTokenRouter(param *gin.RouterGroup) {
	validateTokenRouter := param.Group("")
	validateTokenRouter.Use(middleware.RateLimit())
	//校验token是否有效
	validateTokenRouter.GET("/validate-token/:token", controller.Token.Validate)
}
