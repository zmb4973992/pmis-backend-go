package cron

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func UpdateContractCumulativeIncomeAndExpenditure() {
	var contracts []model.Contract
	err := global.DB.Find(&contracts).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
	}

	for i := range contracts {
		var param1 service.ContractCumulativeIncomeUpdate
		param1.ContractID = contracts[i].ID
		param1.Update()

		var param2 service.ContractCumulativeExpenditureUpdate
		param2.ContractID = contracts[i].ID
		param2.Update()
	}
}
