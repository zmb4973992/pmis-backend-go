package dto

//以下为入参

type DepartmentCreateOrUpdate struct {
	Base
	Name       string `json:"name" binding:"required"`        //部门名称
	Level      string `json:"level" binding:"required"`       //级别，如公司、事业部、部门等
	SuperiorID *int   `json:"superior_id" binding:"required"` //上级机构ID
}

type DepartmentList struct {
	ListInput
	AuthInput
	ID         int     `form:"id"`
	SuperiorID *int    `json:"superior_id"`
	Level      *string `json:"level"`
	Name       *string `json:"name"`
	NameLike   *string `json:"name_like"`

	RoleNames           []string //用户的角色名称数组
	BusinessDivisionIDs []int    //用户所属的事业部id数组
	DepartmentIDs       []int    //用户所属的部门id数组
}

//以下为出参

type DepartmentOutput struct {
	Base       `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Name       string                   `json:"name"  mapstructure:"name"`              //部门名称
	Level      string                   `json:"level" mapstructure:"level"`             //级别，如公司、事业部、部门等
	SuperiorID *int                     `json:"superior_id" mapstructure:"superior_id"` //上级机构
}
