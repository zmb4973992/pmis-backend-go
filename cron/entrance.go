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
		global.SugaredLogger.Errorln("添加定时任务失败，请检查")
	}

	//At 14:02 PM, every day
	_, err = c.AddFunc("02 14 * * *", func() {
		err = updateUsers()
		if err != nil {
			//这里要完善错误处理逻辑，以后再说
		}
	})
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务失败，请检查")
	}

	c.Start()
}

func test() {
	//fmt.Println("666")
}
