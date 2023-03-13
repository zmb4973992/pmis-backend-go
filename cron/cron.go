package cron

import (
	"github.com/robfig/cron/v3"
	"pmis-backend-go/global"
)

func Init() {
	//默认是5位格式: * * * * *
	c := cron.New()
	//添加每分钟执行一次的任务
	_, err := c.AddFunc("* * * * ?", test)
	if err != nil {
		global.SugaredLogger.Panicln("添加定时任务失败，请检查")
	}

	_, err = c.AddFunc("* * * * ?", updateUser)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务失败，请检查")
	}

	c.Start()
}

func test() {
	//fmt.Println("666")
}
