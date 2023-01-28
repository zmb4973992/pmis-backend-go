package model

type DisassemblyView struct {
	BaseModel
	Name       *string  //名称
	ProjectID  *int     //项目id，外键
	SuperiorID *int     //上级id
	Level      *int     //层级
	Weight     *float64 //权重
}
