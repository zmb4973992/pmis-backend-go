package cron

import (
	"errors"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func UpdateProjectCumulativeIncomeAndExpenditure() error {
	var projects []model.Project
	err := global.DB.Find(&projects).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
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
