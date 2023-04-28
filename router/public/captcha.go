package public

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/controller"
	"pmis-backend-go/middleware"
)

type CaptchaRouter struct{}

func (l *CaptchaRouter) InitCaptchaRouter(param *gin.RouterGroup) {
	captchaRouter := param.Group("")
	captchaRouter.Use(middleware.RateLimit())
	captchaRouter.GET("/captcha", controller.Captcha.Get)
}
