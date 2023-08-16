package model

import (
	"time"
)

type Progress struct {
	BasicModel
	DisassemblyID *int64     //拆解情况ID
	Date          *time.Time `gorm:"type:date"`
	Type          *int64
	Value         *float64
	DataSource    *int64
	Remarks       *string
}

// TableName 修改表名
func (p *Progress) TableName() string {
	return "progress"
}
