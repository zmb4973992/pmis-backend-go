package lvmin

import (
	"errors"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
)

func updateCumulativeExpenditure(userID int64, projectIDs []int64, contractIDs []int64) error {
	projectIDs = util.RemoveDuplication(projectIDs)
	for i := range projectIDs {
		var param service.ProjectDailyAndCumulativeExpenditureUpdate
		param.UserID = userID
		param.ProjectID = projectIDs[i]
		res := param.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}
	}

	contractIDs = util.RemoveDuplication(contractIDs)
	for i := range contractIDs {
		var param service.ContractCumulativeExpenditureUpdate
		param.UserID = userID
		param.ContractID = contractIDs[i]
		res := param.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
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
		res := param.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}
	}

	contractIDs = util.RemoveDuplication(contractIDs)
	for i := range contractIDs {
		var param service.ContractCumulativeIncomeUpdate
		param.UserID = userID
		param.ContractID = contractIDs[i]
		res := param.Update()
		if res.Code != 0 {
			return errors.New(res.Message)
		}
	}

	return nil
}
