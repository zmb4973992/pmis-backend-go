package middleware

import (
	"github.com/gin-gonic/gin"
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		//不会写，暂缓
	}
}
