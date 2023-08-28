package middleware

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/util"
	"time"
)

var channelForRequestLog = make(chan model.RequestLog, 5)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行耗时（毫秒）
		timeElapsed := int(endTime.Sub(startTime).Milliseconds())

		//直接操作model更方便
		var requestLog model.RequestLog

		//处理creator、lastModifier、userID字段
		userID, exists := util.GetUserID(c)
		if exists {
			requestLog.Creator = &userID
			requestLog.LastModifier = &userID
		}

		//获取访问路径
		tempPath := c.FullPath()
		requestLog.Path = &tempPath

		//获取URI参数
		//requestLog.URIParams = c.Params

		//获取请求方式
		tempMethod := c.Request.Method
		requestLog.Method = &tempMethod

		//获取ip
		tempIP := c.ClientIP()
		requestLog.IP = &tempIP

		//获取响应码
		tempCode := c.Writer.Status()
		requestLog.ResponseCode = &tempCode

		//获取开始时间和执行耗时(毫秒)
		requestLog.StartTime = &startTime
		requestLog.TimeElapsed = &timeElapsed

		//获取用户的浏览器标识
		tempUserAgent := c.Request.UserAgent()
		requestLog.UserAgent = &tempUserAgent

		//把日志放到通道中，等待保存到数据库
		channelForRequestLog <- requestLog

	}
}

func SaveRequestLog() {
	for {
		select {
		case requestLog := <-channelForRequestLog:
			global.DB.Create(&requestLog)
		}
	}
}
