package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"pmis-backend-go/global"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

func RateLimit() gin.HandlerFunc {
	//采用令牌桶算法，生成限流器。每秒往令牌桶放几个令牌，令牌桶最大容量
	limiter := rate.NewLimiter(
		rate.Limit(global.Config.RateLimit.Limit),
		global.Config.RateLimit.Burst,
	)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		} else {
			c.JSON(
				http.StatusOK,
				response.GenerateCommon(nil, util.ErrorRequestFrequencyTooHigh),
			)
			c.Abort()
			return
		}

	}
}
