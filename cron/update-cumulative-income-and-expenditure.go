package cron

import (
	"errors"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func updateCumulativeIncomeAndExpenditure() {
	err := updateProjectCumulativeIncomeAndExpenditure()
	if err != nil {
		param := service.ErrorLogCreate{Detail: err.Error()}
		param.Create()
	}

	err = updateContractCumulativeIncomeAndExpenditure()
	if err != nil {
		param := service.ErrorLogCreate{Detail: err.Error()}
		param.Create()
	}
}

func updateProjectCumulativeIncomeAndExpenditure() error {
	var projects []model.Project
	err := global.DB.Find(&projects).Error
	if err != nil {
		return err
	}

	for i := range projects {
		var param1 service.ProjectCumulativeIncomeUpdate
		param1.ProjectID = projects[i].ID
		res := param1.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}

		var param2 service.ProjectCumulativeExpenditureUpdate
		param2.ProjectID = projects[i].ID
		param2.Update()
		res = param2.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}
	}

	return nil
}

func updateContractCumulativeIncomeAndExpenditure() error {
	var contract []model.Contract
	err := global.DB.Find(&contract).Error
	if err != nil {
		return err
	}

	for i := range contract {
		var param1 service.ContractCumulativeIncomeUpdate
		param1.ContractID = contract[i].ID
		res := param1.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}

		var param2 service.ContractCumulativeExpenditureUpdate
		param2.ContractID = contract[i].ID
		param2.Update()
		res = param2.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}
	}

	return nil
}
