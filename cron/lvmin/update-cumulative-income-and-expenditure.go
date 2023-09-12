package lvmin

import (
	"pmis-backend-go/service"
	"pmis-backend-go/util"
)

func updateCumulativeExpenditure(userId int64, projectIds []int64, contractIds []int64) error {
	projectIds = util.RemoveDuplication(projectIds)
	for i := range projectIds {
		var param service.ProjectDailyAndCumulativeExpenditureUpdate
		param.UserId = userId
		param.ProjectId = projectIds[i]

		errCode := param.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	contractIds = util.RemoveDuplication(contractIds)
	for i := range contractIds {
		var param service.ContractDailyAndCumulativeExpenditureUpdate
		param.UserId = userId
		param.ContractId = contractIds[i]

		errCode := param.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	return nil
}

func updateCumulativeIncome(userId int64, projectIds []int64, contractIds []int64) error {
	projectIds = util.RemoveDuplication(projectIds)
	for i := range projectIds {
		var param service.ProjectDailyAndCumulativeIncomeUpdate
		param.UserId = userId
		param.ProjectId = projectIds[i]

		errCode := param.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	contractIds = util.RemoveDuplication(contractIds)
	for i := range contractIds {
		var param service.ContractDailyAndCumulativeIncomeUpdate
		param.UserId = userId
		param.ContractId = contractIds[i]

		errCode := param.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	return nil
}
