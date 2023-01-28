package model

type DisassemblyTemplate struct {
	BaseModel
	Name       *string  //名称
	ProjectID  *int     //项目id，外键
	SuperiorID *int     //上级id
	Level      *int     //层级
	Weight     *float64 //权重
}

// TableName 修改表名
func (DisassemblyTemplate) TableName() string {
	return "disassembly_template"
}
