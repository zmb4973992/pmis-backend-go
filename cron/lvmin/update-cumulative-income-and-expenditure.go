package lvmin

import (
	"pmis-backend-go/service"
	"pmis-backend-go/util"
)

func updateCumulativeExpenditure(userID int64, projectIDs []int64, contractIDs []int64) error {
	projectIDs = util.RemoveDuplication(projectIDs)
	for i := range projectIDs {
		var param service.ProjectDailyAndCumulativeExpenditureUpdate
		param.UserID = userID
		param.ProjectID = projectIDs[i]

		errCode := param.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	contractIDs = util.RemoveDuplication(contractIDs)
	for i := range contractIDs {
		var param service.ContractDailyAndCumulativeExpenditureUpdate
		param.UserID = userID
		param.ContractID = contractIDs[i]

		errCode := param.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	return nil
}

func updateCumulativeIncome(userID int64, projectIDs []int64, contractIDs []int64) error {
	projectIDs = util.RemoveDuplication(projectIDs)
	for i := range projectIDs {
		var param service.ProjectDailyAndCumulativeIncomeUpdate
		param.UserID = userID
		param.ProjectID = projectIDs[i]

		errCode := param.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	contractIDs = util.RemoveDuplication(contractIDs)
	for i := range contractIDs {
		var param service.ContractDailyAndCumulativeIncomeUpdate
		param.UserID = userID
		param.ContractID = contractIDs[i]

		errCode := param.Update()
		if errCode != util.Success {
			return util.GenerateCustomError(errCode)
		}
	}

	return nil
}
