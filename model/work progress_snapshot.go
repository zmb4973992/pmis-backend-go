package model

type WorkProgressSnapshot struct {
	BaseModel
	DisassemblyID *int    //项目拆解id，外键
	Date          *string `gorm:"type:date;"` //日期  默认格式为2020-02-02

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
func (WorkProgressSnapshot) TableName() string {
	return "work_progress_snapshot"
}
