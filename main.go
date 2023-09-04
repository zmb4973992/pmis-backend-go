package main

import (
	"pmis-backend-go/controller"
	"pmis-backend-go/global"
	"pmis-backend-go/middleware"
	"pmis-backend-go/model"
	"pmis-backend-go/router"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
)

func main() {
	//加载全局变量
	global.InitConfig()
	//加载日志记录器，使用的是zap
	util.InitLogger()
	//加载ID生成器
	util.InitIDGenerator()
	//连接数据库
	model.InitDatabase()
	//创建保存上传文件的文件夹
	controller.Init()
	//生成引擎
	engine := router.InitEngine()
	//开启2个协程，用来保存访问记录到数据库
	for i := 0; i < 2; i++ {
		go middleware.SaveRequestLog()
	}

	//第一次运行时，导入初始数据
	//disposable.Init()

	//开启定时任务
	//cron.Init()

	var param service.OperationLogCreate
	param.Operator = 454880052387845
	param.Date = "2020-03-03"
	param.OperationType = 455587328651269
	param.Detail = "sdkfjskfjd"
	param.ProjectID = 454882673385477
	param.Create()

	//运行服务，必须放最后
	err := engine.Run(":" + global.Config.AppConfig.HttpPort)

	if err != nil {
		global.SugaredLogger.Panicln(err)
	}
}
