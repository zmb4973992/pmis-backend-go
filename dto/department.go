package dto

type DepartmentGetDTO struct {
	Name     string `json:"name"`     //部门名称
	Level    string `json:"level"`    //级别，如公司、事业部、部门等
	Superior any    `json:"superior"` //上级机构
}

type DepartmentCreateAndUpdateDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name" binding:"required"`        //部门名称
	Level      string `json:"level" binding:"required"`       //级别，如公司、事业部、部门等
	SuperiorID *int   `json:"superior_id" binding:"required"` //上级机构ID
}

// DepartmentListDTO 待实现
type DepartmentListDTO struct {
	ID    int  `json:"id"`
	IDGte *int `json:"id_gte"`
	IDLte *int `json:"id_lte"`

	Name        *string `json:"name"`
	NameInclude *string `json:"name_include"`

	ListDTO
}
