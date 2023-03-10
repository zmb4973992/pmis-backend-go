package model

import "time"

type Progress struct {
	BaseModel
	DisassemblyID *int       //拆解情况id
	Date          *time.Time `gorm:"type:date"` //日期

	ProgressType *int
	Progress     *float64
	Remark       *string
	DataSource   *string

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
func (*Progress) TableName() string {
	return "work_progress"
}
