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
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//开始解析token
	res, err := util.ParseToken(accessToken)
	//如果存在错误或token已过期
	if err != nil || res.ExpiresAt.Unix() < time.Now().Unix() {
		c.JSON(http.StatusOK, response.Failure(util.ErrorAccessTokenInvalid))
		c.Abort()
		return
	}

	userID := res.UserID
	res1 := map[string]int64{
		"user_id": userID,
	}

	c.JSON(http.StatusOK, response.SuccessWithData(res1))
	return
}
