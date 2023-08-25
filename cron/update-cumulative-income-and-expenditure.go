package cron

import (
	"errors"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func UpdateCumulativeIncomeAndExpenditureForCron() {
	var user model.User
	err := global.DB.Where("username = 'z0030975'").First(&user).Error
	if err != nil {
		global.SugaredLogger.Panicln(err)
	}

	err = UpdateCumulativeIncomeAndExpenditure(user.ID)
	if err != nil {
		param := service.ErrorLogCreate{Detail: err.Error()}
		param.Create()
	}
}

func UpdateCumulativeIncomeAndExpenditure(userID int64) error {
	err := updateProjectCumulativeIncomeAndExpenditure(userID)
	if err != nil {
		return err
	}

	err = updateContractCumulativeIncomeAndExpenditure(userID)
	if err != nil {
		return err
	}

	return nil
}

func updateProjectCumulativeIncomeAndExpenditure(userID int64) error {
	var projects []model.Project
	err := global.DB.Find(&projects).Error
	if err != nil {
		return err
	}

	for i := range projects {
		var param1 service.ProjectDailyAndCumulativeIncomeUpdate
		param1.UserID = userID
		param1.ProjectID = projects[i].ID
		res := param1.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}

		var param2 service.ProjectDailyAndCumulativeExpenditureUpdate
		param2.UserID = userID
		param2.ProjectID = projects[i].ID
		res = param2.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}
	}

	return nil
}

func updateContractCumulativeIncomeAndExpenditure(userID int64) error {
	var contract []model.Contract
	err := global.DB.Find(&contract).Error
	if err != nil {
		return err
	}

	for i := range contract {
		var param1 service.ContractDailyAndCumulativeIncomeUpdate
		param1.UserID = userID
		param1.ContractID = contract[i].ID
		res := param1.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}

		var param2 service.ContractDailyAndCumulativeExpenditureUpdate
		param2.UserID = userID
		param2.ContractID = contract[i].ID
		res = param2.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}
	}

	return nil
}
