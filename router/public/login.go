package public

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
)

type LoginRouter struct{}

func (l *LoginRouter) InitLoginRouter(param *gin.RouterGroup) {
	loginRouter := param.Group("")
	loginRouter.Use(middleware.RateLimit())
	loginRouter.POST("/login", controller.Login)
}
