package cron

import (
	"errors"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func UpdateContractCumulativeIncomeAndExpenditure() error {
	var contracts []model.Contract
	err := global.DB.Find(&contracts).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}

	for i := range contracts {
		var param1 service.ContractCumulativeIncomeUpdate
		param1.ContractID = contracts[i].ID
		res := param1.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}

		var param2 service.ContractCumulativeExpenditureUpdate
		param2.ContractID = contracts[i].ID
		res = param2.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}
	}

	return nil
}
