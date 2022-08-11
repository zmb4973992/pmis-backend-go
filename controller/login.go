package controller

import (
	"github.com/gin-gonic/gin"
	"learn-go/dto"
	"learn-go/serializer/response"
	"learn-go/service"
	"learn-go/util"
	"net/http"
)

func Login(c *gin.Context) {
	var loginDTO dto.LoginDTO
	err := c.ShouldBindJSON(&loginDTO)
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.LoginService.Login(loginDTO)
	c.JSON(http.StatusOK, res)
}
