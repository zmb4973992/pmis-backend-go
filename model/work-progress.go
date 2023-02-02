package model

import "time"

type WorkProgress struct {
	BaseModel
	DisassemblyID         *int       //拆解情况id，外键
	DisassemblyIDWithDate *string    //带日期的拆解情况id
	FillingDate           *time.Time `gorm:"type:date"` //添加记录的日期
	Date                  *time.Time `gorm:"type:date"` //日期

	PlannedProgress             *float64 //初始计划进度
	RemarkOfPlannedProgress     *string  //初始计划进度的备注
	DataSourceOfPlannedProgress *string  //初始计划进度的数据来源

	ActualProgress                    *float64 //实际进度
	RemarkOfActualProgress            *string  //实际进度的备注
	DataSourceOfActualProgress        *string  //实际进度的数据来源
	ActualProgressIsTemporarilyFilled *bool    //实际进度是否为临时填充

	PredictedProgress             *float64 //预测进度
	RemarkOfPredictedProgress     *string  //预测进度的备注
	DataSourceOfPredictedProgress *string  //预测进度的数据来源
}

// TableName 修改表名
func (*WorkProgress) TableName() string {
	return "work_progress"
}
