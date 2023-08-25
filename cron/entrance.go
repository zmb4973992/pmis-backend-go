package cron

import (
	"github.com/robfig/cron/v3"
	"pmis-backend-go/cron/lvmin"
	"pmis-backend-go/cron/old-pmis"
	"pmis-backend-go/cron/windows-ad"
	"pmis-backend-go/global"
)

func Init() {
	//默认是5位格式: * * * * *
	c := cron.New()

	_, err := c.AddFunc("33 * * * *", windows_ad.UpdateUsersForCron)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务”windows_ad.UpdateUsersForCron“失败，请检查")
	}

	_, err = c.AddFunc("35 * * * *", lvmin.ImportDataForCron)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务”lvmin.ImportDataForCron“失败，请检查")
	}

	_, err = c.AddFunc("44 * * * *", old_pmis.ImportDataForCron)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务”old_pmis.ImportDataForCron“失败，请检查")
	}

	_, err = c.AddFunc("* * * * *", UpdateCumulativeIncomeAndExpenditureForCron)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务”UpdateCumulativeIncomeAndExpenditureForCron“失败，请检查")
	}

	_, err = c.AddFunc("12 22 * * *", clearUnusedFiles)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务”clearUnusedFiles“失败，请检查")
	}

	c.Start()
}
