package public

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
)

type JWTRouter struct{}

func (j *JWTRouter) InitJWTRouter(param *gin.RouterGroup) {
	jwtRouter := param.Group("")
	jwtRouter.Use(middleware.RateLimit())
	//校验token是否有效
	jwtRouter.GET("/jwt/:token", controller.Token.Validate)
}
