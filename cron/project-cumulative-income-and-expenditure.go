package cron

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func UpdateProjectCumulativeIncomeAndExpenditure() {
	var projects []model.Project
	err := global.DB.Find(&projects).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
	}

	for i := range projects {
		var param1 service.ProjectCumulativeIncomeUpdate
		param1.ProjectID = projects[i].ID
		param1.Update()

		var param2 service.ProjectCumulativeExpenditureUpdate
		param2.ProjectID = projects[i].ID
		param2.Update()
	}
}
