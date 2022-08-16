package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
)

type operationRecordDAO struct{}

func (operationRecordDAO) Get(operationRecordID int) *dto.OperationRecordGetDTO {
	var paramOut dto.OperationRecordGetDTO
	//把基础的拆解信息查出来
	var operationRecord model.OperationRecord
	err := global.DB.Where("id = ?", operationRecordID).First(&operationRecord).Error
	if err != nil {
		return nil
	}
	//把所有查出的结果赋值给输出变量
	if operationRecord.ProjectID != nil {
		paramOut.ProjectID = operationRecord.ProjectID
	}
	if operationRecord.OperatorID != nil {
		paramOut.OperatorID = operationRecord.OperatorID
	}
	if operationRecord.Date != nil {
		date := operationRecord.Date.Format("2006-01-02")
		paramOut.Date = &date
	}
	if operationRecord.Action != nil {
		paramOut.Action = operationRecord.Action
	}
	if operationRecord.Detail != nil {
		paramOut.Detail = operationRecord.Detail
	}

	return &paramOut
}

// Create 这里是只负责新增，不写任何业务逻辑。只要收到参数就创建数据库记录，然后返回错误
func (operationRecordDAO) Create(param *model.OperationRecord) error {
	err := global.DB.Create(param).Error
	return err
}

// Update 这里是只负责更新，不写任何业务逻辑。只要收到id和更新参数，然后返回错误
func (operationRecordDAO) Update(param *model.OperationRecord) error {
	//注意，这里就算没有找到记录，也不会报错，只有更新字段出现问题才会报错。详见gorm的update用法
	err := global.DB.Where("id = ?", param.ID).Omit("created_at").Save(param).Error
	return err
}

func (operationRecordDAO) Delete(operationRecordID int) error {
	//注意，这里就算没有找到记录，也不会报错。详见gorm的delete用法
	err := global.DB.Delete(&model.OperationRecord{}, operationRecordID).Error
	return err
}
