package controller

import (
	"github.com/gin-gonic/gin"
	"learn-go/serializer/response"
	"learn-go/util"
	"learn-go/util/jwt"
	"net/http"
	"time"
)

type tokenController struct{}

func (tokenController) Validate(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//开始解析token
	res, err := jwt.ParseToken(token)
	//如果存在错误或token已过期
	if err != nil || res.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusOK, response.Failure(util.ErrorAccessTokenInvalid))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.Success())
	return
}
