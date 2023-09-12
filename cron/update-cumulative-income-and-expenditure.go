package cron

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
)

func UpdateCumulativeIncomeAndExpenditureForCron() {
	var user model.User
	err := global.DB.Where("username = 'z0030975'").
		First(&user).Error
	if err != nil {
		global.SugaredLogger.Panicln(err)
	}

	err = UpdateCumulativeIncomeAndExpenditure(user.Id)
	if err != nil {
		param := service.ErrorLogCreate{Detail: err.Error()}
		param.Create()
	}
}

func UpdateCumulativeIncomeAndExpenditure(userId int64) error {
	err := updateProjectCumulativeIncomeAndExpenditure(userId)
	if err != nil {
		return err
	}

	err = updateContractCumulativeIncomeAndExpenditure(userId)
	if err != nil {
		return err
	}

	return nil
}

func updateProjectCumulativeIncomeAndExpenditure(userId int64) error {
	var projects []model.Project
	err := global.DB.Find(&projects).Error
	if err != nil {
		return err
	}

	for i := range projects {
		var param1 service.ProjectDailyAndCumulativeIncomeUpdate
		param1.UserId = userId
		param1.ProjectId = projects[i].Id

		errCode := param1.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}

		var param2 service.ProjectDailyAndCumulativeExpenditureUpdate
		param2.UserId = userId
		param2.ProjectId = projects[i].Id

		errCode = param2.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	return nil
}

func updateContractCumulativeIncomeAndExpenditure(userId int64) error {
	var contract []model.Contract
	err := global.DB.Find(&contract).Error
	if err != nil {
		return err
	}

	for i := range contract {
		var param1 service.ContractDailyAndCumulativeIncomeUpdate
		param1.UserId = userId
		param1.ContractId = contract[i].Id

		errCode := param1.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}

		var param2 service.ContractDailyAndCumulativeExpenditureUpdate
		param2.UserId = userId
		param2.ContractId = contract[i].Id

		errCode = param2.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	return nil
}
