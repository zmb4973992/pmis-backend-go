package model

import "time"

type OperationRecord struct {
	BaseModel
	ProjectID  *int       //项目id
	OperatorID *int       //操作人id
	Date       *time.Time `gorm:"type:date;"` //日期
	Action     *string    //动作
	Detail     *string    //详情
}

// TableName 修改数据库的表名
func (OperationRecord) TableName() string {
	return "operation_record"
}
