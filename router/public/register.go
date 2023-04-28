package public

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
)

type RegisterRouter struct{}

func (r *RegisterRouter) InitRegisterRouter(param *gin.RouterGroup) {
	registerRouter := param.Group("")
	registerRouter.Use(middleware.RateLimit())
	registerRouter.POST("/user", controller.User.Create)
}
