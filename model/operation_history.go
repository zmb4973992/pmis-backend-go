package model

type OperationHistory struct {
	BaseModel
	ProjectID  *int    //项目id
	OperatorID *int    //操作人id
	Date       *string //日期
	Action     *string //动作
	Detail     *string //详情
}

// TableName 修改数据库的表名
func (OperationHistory) TableName() string {
	return "operation_history"
}
