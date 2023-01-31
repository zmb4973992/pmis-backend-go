package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
)

func Login(c *gin.Context) {
	var loginDTO dto.Login
	err := c.ShouldBindJSON(&loginDTO)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		c.JSON(http.StatusOK, response.Fail(util.ErrorInvalidJSONParameters))
		return
	}
	res := service.LoginService.Login(loginDTO)
	c.JSON(http.StatusOK, res)
}
