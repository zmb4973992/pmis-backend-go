package dto

type DepartmentGetDTO struct {
	Name       string `json:"name"  mapstructure:"name"`              //部门名称
	Level      string `json:"level" mapstructure:"level"`             //级别，如公司、事业部、部门等
	SuperiorID *int   `json:"superior_id" mapstructure:"superior_id"` //上级机构
}

// DepartmentCreateOrUpdateDTO
// 除id外，所有字段都设置为必须绑定
type DepartmentCreateOrUpdateDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name" binding:"required"`        //部门名称
	Level      string `json:"level" binding:"required"`       //级别，如公司、事业部、部门等
	SuperiorID *int   `json:"superior_id" binding:"required"` //上级机构ID
}

// DepartmentListDTO 待实现
type DepartmentListDTO struct {
	ID int `json:"id"`

	SuperiorID *int    `json:"superior_id"`
	Level      *string `json:"level"`

	Name        *string `json:"name"`
	NameInclude *string `json:"name_include"`

	ListDTO
}
