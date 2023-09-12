package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

type token struct{}

func (t *token) Validate(c *gin.Context) {
	accessToken := c.Param("access-token")
	if accessToken == "" {
		c.JSON(
			http.StatusBadRequest,
			response.GenerateCommon(nil, util.ErrorInvalidURIParameters),
		)
		return
	}

	//开始解析token
	res, err := util.ParseToken(accessToken)
	//如果存在错误或token已过期
	if err != nil || res.ExpiresAt.Unix() < time.Now().Unix() {
		c.JSON(
			http.StatusOK,
			response.GenerateCommon(nil, util.ErrorAccessTokenInvalid),
		)
		c.Abort()
		return
	}

	userId := res.UserId

	c.JSON(
		http.StatusOK,
		response.GenerateCommon(
			gin.H{"user_id": userId},
			util.Success),
	)
	return
}
