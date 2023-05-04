package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"time"
)

var channelOfOperationLogs = make(chan model.OperationLog, 50)

func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行耗时（毫秒）
		timeElapsed := int(endTime.Sub(startTime).Milliseconds())
		//fmt.Println("命令耗时:", timeElapsed, "毫秒")

		//这里不用dto了，因为都是程序内部操作，不对外暴露接口和数据
		//直接操作model更方便
		var operationLog model.OperationLog

		//处理creator、lastModifier、userID字段
		tempUserID1, exists := c.Get("user_id")
		if exists {
			tempUserID2 := tempUserID1.(int)
			operationLog.UserSnowID = &tempUserID2
			operationLog.Creator = &tempUserID2
			operationLog.LastModifier = &tempUserID2
		}
		//获取访问路径
		tempPath := c.FullPath()
		operationLog.Path = &tempPath

		//获取URI参数
		operationLog.URIParams = c.Params

		//获取请求方式
		tempMethod := c.Request.Method
		operationLog.Method = &tempMethod

		fmt.Println(c.Request.Context())

		//获取ip
		tempIP := c.ClientIP()
		operationLog.IP = &tempIP

		//获取响应码
		tempCode := c.Writer.Status()
		operationLog.ResponseCode = &tempCode

		//获取开始时间和执行耗时(毫秒)
		operationLog.StartTime = &startTime
		operationLog.TimeElapsed = &timeElapsed

		//获取用户的浏览器标识
		tempUserAgent := c.Request.UserAgent()
		operationLog.UserAgent = &tempUserAgent

		//把日志放到通道中，等待保存到数据库
		channelOfOperationLogs <- operationLog

	}
}

func SaveOperationLog() {
	for {
		select {
		case log := <-channelOfOperationLogs:
			global.DB.Create(&log)
		}
	}
}
