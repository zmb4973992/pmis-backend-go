package dto

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type DepartmentCreate struct {
	Creator      int
	LastModifier int
	Name         string `json:"name" binding:"required"`             //名称
	Level        string `json:"level" binding:"required"`            //级别，如公司、事业部、部门等
	SuperiorID   int    `json:"superior_id" binding:"required,gt=0"` //上级机构ID
}

type DepartmentCreateOrUpdateOld struct {
	Base
	Name       string `json:"name" binding:"required"`        //部门名称
	Level      string `json:"level" binding:"required"`       //级别，如公司、事业部、部门等
	SuperiorID *int   `json:"superior_id" binding:"required"` //上级机构ID
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DepartmentUpdate struct {
	LastModifier int
	ID           int
	Name         *string `json:"name"`        //名称
	Level        *string `json:"level"`       //级别，如公司、事业部、部门等
	SuperiorID   *int    `json:"superior_id"` //上级机构ID
}

type DepartmentDelete struct {
	Deleter int
	ID      int
}

type DepartmentList struct {
	ListInput
	RoleNames           []string //用户的角色名称数组
	BusinessDivisionIDs []int    //用户所属的事业部id数组
	DepartmentIDs       []int    //用户所属的部门id数组

	SuperiorID int    `json:"superior_id,omitempty"`
	Level      string `json:"level,omitempty"`
	Name       string `json:"name,omitempty"`
	NameLike   string `json:"name_like,omitempty"`
}

//以下为出参

type DepartmentOutput struct {
	Creator      *int    `json:"creator" gorm:"creator"`
	LastModifier *int    `json:"last_modifier" gorm:"last_modifier"`
	ID           int     `json:"id" gorm:"id"`
	Name         string  `json:"name" gorm:"name"`               //名称
	Level        *string `json:"level" gorm:"level"`             //级别，如公司、事业部、部门等
	SuperiorID   *int    `json:"superior_id" gorm:"superior_id"` //上级机构id
}

type DepartmentOutputOld struct {
	Base       `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Name       string                   `json:"name"  mapstructure:"name"`              //部门名称
	Level      string                   `json:"level" mapstructure:"level"`             //级别，如公司、事业部、部门等
	SuperiorID *int                     `json:"superior_id" mapstructure:"superior_id"` //上级机构
}

//以下为入参

type DepartmentListOld struct {
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
