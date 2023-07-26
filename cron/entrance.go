package cron

import (
	"github.com/robfig/cron/v3"
	"pmis-backend-go/global"
)

func Init() {
	//默认是5位格式: * * * * *
	c := cron.New()

	_, err := c.AddFunc("5 * * * *", connectToLvmin)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务失败，请检查")
	}

	_, err = c.AddFunc("35 * * * *", updateUsersByLDAP)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务失败，请检查")
	}

	_, err = c.AddFunc("45 14 * * *", updateCumulativeIncomeAndExpenditure)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务失败，请检查")
	}

	c.Start()
}
