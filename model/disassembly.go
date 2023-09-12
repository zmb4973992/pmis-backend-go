package model

type Disassembly struct {
	BasicModel
	Name                  *string  //名称
	ProjectId             *int64   //项目id
	SuperiorId            *int64   //上级id
	Level                 *int     //层级
	Weight                *float64 //权重
	Sort                  *int     //排序值
	ImportedIdFromOldPmis *int64   //老PMIS的拆解情况id
}

// TableName 修改表名
func (d *Disassembly) TableName() string {
	return "disassembly"
}
