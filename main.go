package main

import (
	"pmis-backend-go/cron"
	"pmis-backend-go/global"
	"pmis-backend-go/middleware"
	"pmis-backend-go/model"
	"pmis-backend-go/router"
	"pmis-backend-go/util"
	"pmis-backend-go/util/logger"
)

func main() {
	//加载全局变量,如应用基础设置、数据库连接信息、jwt信息、日志设置登
	global.InitConfig()

	//加载日志记录器，使用的是zap
	logger.InitLogger()

	//连接数据库
	model.InitDatabase()

	//创建保存上传文件的文件夹
	util.UploadInit()

	//开始采用自定义的方式生成引擎
	engine := router.Init()

	//开启4个协程，用来保存访问记录到数据库
	for i := 0; i < 4; i++ {
		go middleware.SaveOperationLog()
	}

	//开启定时任务
	cron.Init()

	//运行服务
	err := engine.Run(":" + global.Config.AppConfig.HttpPort)

	if err != nil {
		global.SugaredLogger.Panicln(err)
	}
}
