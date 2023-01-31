package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"pmis-backend-go/util/jwt"
	"time"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("access_token")
		//如果请求头没有携带access_token
		if token == "" {
			c.JSON(http.StatusOK, response.Fail(util.ErrorAccessTokenNotFound))
			c.Abort()
			return
		}
		//开始校验access_token
		res, err := jwt.ParseToken(token)
		//如果存在错误或token已过期
		if err != nil || res.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusOK, response.Fail(util.ErrorAccessTokenInvalid))
			c.Abort()
			return
		}
		//如果access_token校验通过
		c.Set("user_id", res.UserID)
		c.Next()
		return
	}
}
