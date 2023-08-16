package old_pmis

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/service"
	"time"
)

type progress struct {
	DisassemblyID               int64      `gorm:"column:拆解情况id"`
	Date                        *time.Time `gorm:"column:日期;type:date"`
	PlannedProgress             *float64   `gorm:"column:初始计划进度"`
	RemarksOfPlannedProgress    string     `gorm:"column:初始计划进度的备注"`
	ActualProgress              *float64   `gorm:"column:实际进度"`
	RemarksOfActualProgress     string     `gorm:"column:实际进度的备注"`
	ForecastedProgress          *float64   `gorm:"column:预测进度"`
	RemarksOfForecastedProgress string     `gorm:"column:预测进度的备注"`
}

func importProgress(userID int64) error {
	fmt.Println("★★★★★开始处理进度记录......★★★★★")

	var originalCountOfProgressRecords int64
	global.DB.Model(&model.Progress{}).
		Count(&originalCountOfProgressRecords)

	var oldProgresses []progress
	global.DBForOldPmis.Table("工作进度").
		Where("日期 is not null").
		Where("初始计划进度的数据来源 = '人工填写'").
		Or("实际进度的数据来源 = '人工填写'").
		Or("预测进度的数据来源 = '人工填写'").
		Find(&oldProgresses)

	var progressType model.DictionaryType
	err := global.DB.Where("name = '进度类型'").
		First(&progressType).Error
	if err != nil {
		return err
	}

	var planned model.DictionaryDetail
	err = global.DB.
		Where("dictionary_type_id = ?", progressType.ID).
		Where("name = '计划进度'").
		First(&planned).Error
	if err != nil {
		return err
	}

	var actual model.DictionaryDetail
	err = global.DB.
		Where("dictionary_type_id = ?", progressType.ID).
		Where("name = '实际进度'").
		First(&actual).Error
	if err != nil {
		return err
	}

	var forecasted model.DictionaryDetail
	err = global.DB.
		Where("dictionary_type_id = ?", progressType.ID).
		Where("name = '预测进度'").
		First(&forecasted).Error
	if err != nil {
		return err
	}

	var dataSourceOfProgress model.DictionaryType
	err = global.DB.Where("name = '进度的数据来源'").
		First(&dataSourceOfProgress).Error
	if err != nil {
		return err
	}

	var manualFilling model.DictionaryDetail
	err = global.DB.
		Where("dictionary_type_id = ?", dataSourceOfProgress.ID).
		Where("name = '人工填写'").
		First(&manualFilling).Error
	if err != nil {
		return err
	}

	for i := range oldProgresses {
		var newDisassembly model.Disassembly
		err = global.DB.
			Where("imported_id_from_old_pmis = ?", oldProgresses[i].DisassemblyID).
			First(&newDisassembly).Error
		if err != nil {
			continue
		}

		var newProgress model.Progress
		newProgress.Creator = &userID
		newProgress.DisassemblyID = &newDisassembly.ID
		newProgress.Date = oldProgresses[i].Date
		newProgress.DataSource = &manualFilling.ID
		if oldProgresses[i].PlannedProgress != nil {
			newProgress.Type = &planned.ID
			newProgress.Value = oldProgresses[i].PlannedProgress
			newProgress.Remarks = &oldProgresses[i].RemarksOfPlannedProgress
		} else if oldProgresses[i].ActualProgress != nil {
			newProgress.Type = &actual.ID
			newProgress.Value = oldProgresses[i].ActualProgress
			newProgress.Remarks = &oldProgresses[i].RemarksOfActualProgress
		} else if oldProgresses[i].ForecastedProgress != nil {
			newProgress.Type = &forecasted.ID
			newProgress.Value = oldProgresses[i].ForecastedProgress
			newProgress.Remarks = &oldProgresses[i].RemarksOfForecastedProgress
		}

		err = global.DB.
			Where("disassembly_id = ?", newProgress.DisassemblyID).
			Where("date = ?", newProgress.Date).
			Where("data_source = ?", manualFilling.ID).
			Where("type = ?", newProgress.Type).
			Where("value = ?", newProgress.Value).
			Where("remarks = ?", newProgress.Remarks).
			FirstOrCreate(&newProgress).Error

		if err != nil {
			global.SugaredLogger.Errorln(err)
			continue
		}

	}

	var NewCountOfProgressRecords int64
	global.DB.Model(&model.Progress{}).
		Count(&NewCountOfProgressRecords)

	if NewCountOfProgressRecords != originalCountOfProgressRecords {
		var projects []model.Project
		global.DB.Find(&projects)
		for i := range projects {
			var param service.ProgressUpdateByProjectID
			param.UserID = userID
			param.ProjectID = projects[i].ID
			err = param.UpdateByProjectID()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("★★★★★进度记录处理完成......★★★★★")

	return nil
}
