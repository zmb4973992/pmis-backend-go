package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/service"
)

type captcha struct {
}

func (ca *captcha) Get(c *gin.Context) {
	var param service.CaptchaGet
	res := param.Get()
	c.JSON(http.StatusOK, res)
	return
}
