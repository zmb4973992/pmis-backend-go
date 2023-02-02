package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"pmis-backend-go/util/jwt"
	"time"
)

type token struct{}

func (token) Validate(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest,
			response.Fail(util.ErrorInvalidURIParameters))
		return
	}

	//开始解析token
	res, err := jwt.ParseToken(token)
	//如果存在错误或token已过期
	if err != nil || res.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusOK, response.Fail(util.ErrorAccessTokenInvalid))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.Succeed())
	return
}
