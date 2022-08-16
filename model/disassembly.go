package model

type Disassembly struct {
	BaseModel
	Name           *string        //名称
	ProjectID      *int           //项目id，外键
	SuperiorID     *int           //上级id
	Level          *int           //层级
	Weight         *float64       //权重
	WorkProgresses []WorkProgress `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 修改表名
func (Disassembly) TableName() string {
	return "disassembly"
}
