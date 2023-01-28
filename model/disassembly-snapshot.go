package model

import "time"

type DisassemblySnapshot struct {
	BaseModel
	Name               *string    //名称
	Date               *time.Time `gorm:"type:datetime"` //添加记录的日期
	IDWithDate         *string    //带日期的拆解情况id
	ProjectID          *int       //项目id，外键
	SuperiorIDWithDate *string    //带日期的上级id
	Level              *int       //层级
	Weight             *float64   //权重
	//WorkProgresses     []WorkProgress `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 修改表名
func (DisassemblySnapshot) TableName() string {
	return "disassembly_snapshot"
}
