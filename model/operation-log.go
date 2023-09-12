package model

import (
	"time"
)

type OperationLog struct {
	BasicModel
	ProjectId     *int64
	Operator      *int64
	Date          *time.Time `gorm:"type:date"`
	OperationType *int64
	Detail        *string
}

// TableName 修改数据库的表名
func (r *OperationLog) TableName() string {
	return "operation_log"
}
