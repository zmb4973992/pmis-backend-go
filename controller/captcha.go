package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
)

type captcha struct {
}

func (ca *captcha) Get(c *gin.Context) {
	var param service.CaptchaGet
	output, errCode := param.Get()
	c.JSON(
		http.StatusOK,
		response.GenerateCommon(output, errCode),
	)
	return
}
