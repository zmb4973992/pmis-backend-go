package cron

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
)

func UpdateProjectCumulativeExpenditure() {
	var projects []model.Project
	err := global.DB.Find(&projects).Error
	if err != nil {
		global.SugaredLogger.Panicln("错误")
	}

	for i := range projects {
		var param service.ProjectCumulativeExpenditureUpdate
		param.ProjectID = projects[i].ID
		param.Update()
	}
}
