package model

import (
	"time"
)

type Progress struct {
	BasicModel
	DisassemblyId *int64     //拆解情况id
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
