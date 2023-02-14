package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/service"
)

type captcha struct {
}

func (*captcha) Get(c *gin.Context) {
	param := service.CaptchaGet{}
	res := param.Get()
	c.JSON(http.StatusOK, res)
	return

}
