package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pmis-backend-go/global"
	"time"
)

// ZapLogger 接管gin框架默认的日志，用作中间件
func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		global.Logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			//暂时不需要显示user-agent，以后需要了可以再打开
			//zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
