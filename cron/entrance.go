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

	//_, err = c.AddFunc("5 * * * *", UpdateProjectCumulativeIncomeAndExpenditure)
	//if err != nil {
	//	global.SugaredLogger.Panicln("添加定时任务失败，请检查")
	//}
	//
	//_, err = c.AddFunc("* * * * *", UpdateContractCumulativeIncomeAndExpenditure)
	//if err != nil {
	//	global.SugaredLogger.Panicln("添加定时任务失败，请检查")
	//}
	//
	////At 20 minutes past the hour, every hour, every day
	//_, err = c.AddFunc("20 * * * *", importDataFromLvmin)
	//if err != nil {
	//	global.SugaredLogger.Panicln("添加定时任务失败，请检查")
	//}

	//At 15:50, every day
	_, err = c.AddFunc("50 15 * * *", updateUsers)
	if err != nil {
		global.SugaredLogger.Errorln("添加定时任务失败，请检查")
	}

	c.Start()
}

func importDataFromLvmin() error {
	err := importRelatedParty()
	if err != nil {
		return err
	}

	err = importProject()
	if err != nil {
		return err
	}

	err = importContract()
	if err != nil {
		return err
	}

	err = importActualExpenditure()
	if err != nil {
		return err
	}

	err = importForecastedExpenditure()
	if err != nil {
		return err
	}

	err = importPlannedExpenditure()
	if err != nil {
		return err
	}

	err = importActualIncome()
	if err != nil {
		return err
	}

	return nil
}
