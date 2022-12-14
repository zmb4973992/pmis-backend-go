package main

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/router"
	"pmis-backend-go/util"
	"pmis-backend-go/util/logger"
	"pmis-backend-go/util/snowflake"
)

func main() {
	//加载全局变量
	global.Init()

	//加载日志记录器，使用的是zap
	logger.Init()

	//连接数据库
	model.Init()

	//初始化snowflake，用来生成唯一ID
	snowflake.Init()

	//创建保存上传文件的文件夹
	util.UploadInit()

	//开始采用自定义的方式生成引擎
	engine := router.Init()

	//global.Logger.Debug("系统配置正常", zap.String("当前运行模式：", global.Config.AppMode))
	//
	//运行服务
	err := engine.Run(":" + global.Config.HttpPort)

	if err != nil {
		panic(err)
	}
}
