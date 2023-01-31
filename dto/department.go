package dto

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type DepartmentCreate struct {
	Creator      int
	LastModifier int
	Name         string `json:"name" binding:"required"`             //名称
	LevelName    string `json:"level_name" binding:"required"`       //级别，如公司、事业部、部门等
	SuperiorID   int    `json:"superior_id" binding:"required,gt=0"` //上级机构ID
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DepartmentUpdate struct {
	LastModifier int
	ID           int
	Name         *string `json:"name"`        //名称
	LevelName    *string `json:"level_name"`  //级别，如公司、事业部、部门等
	SuperiorID   *int    `json:"superior_id"` //上级机构ID
}

type DepartmentDelete struct {
	Deleter int
	ID      int
}

type DepartmentList struct {
	ListInput
	AuthInput
	SuperiorID int    `json:"superior_id,omitempty"`
	LevelName  string `json:"level_name,omitempty"`
	Name       string `json:"name,omitempty"`
	NameLike   string `json:"name_like,omitempty"`
}

//以下为出参

type DepartmentOutput struct {
	Creator      *int    `json:"creator" gorm:"creator"`
	LastModifier *int    `json:"last_modifier" gorm:"last_modifier"`
	ID           int     `json:"id" gorm:"id"`
	Name         string  `json:"name" gorm:"name"`               //名称
	LevelName    *string `json:"level_name" gorm:"level_name"`   //级别，如公司、事业部、部门等
	SuperiorID   *int    `json:"superior_id" gorm:"superior_id"` //上级机构id
}

type DepartmentOutputOld struct {
	Base       `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Name       string                   `json:"name"  mapstructure:"name"`              //部门名称
	Level      string                   `json:"level" mapstructure:"level"`             //级别，如公司、事业部、部门等
	SuperiorID *int                     `json:"superior_id" mapstructure:"superior_id"` //上级机构
}
