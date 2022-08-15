package dto

type ProjectDisassemblyGetDTO struct {
	Name       *string  `json:"name" mapstructure:"name"`               //名称
	ProjectID  *int     `json:"project_id" mapstructure:"project_id"`   //所属项目id
	Level      *int     `json:"level" mapstructure:"level"`             //层级
	Weight     *float64 `json:"weight" mapstructure:"weight"`           //权重
	SuperiorID *int     `json:"superior_id" mapstructure:"superior_id"` //上级拆解项id
}

// ProjectDisassemblyCreateOrUpdateDTO
// 除id外，所有字段都设置为必须绑定
type ProjectDisassemblyCreateOrUpdateDTO struct {
	ID         int      `json:"id"`
	Name       *string  `json:"name" binding:"required"`        //拆解项名称
	ProjectID  *int     `json:"project_id" binding:"required"`  //所属项目id
	Level      *int     `json:"level" binding:"required"`       //层级
	Weight     *float64 `json:"weight" binding:"required"`      //权重
	SuperiorID *int     `json:"superior_id" binding:"required"` //上级拆解项ID
}

// ProjectDisassemblyListDTO 是list查询的过滤器
// 在dto传递给sqlCondition时，空值会被忽略
type ProjectDisassemblyListDTO struct {
	ID    int  `json:"id"`
	IDGte *int `json:"id_gte"`
	IDLte *int `json:"id_lte"`

	Name        *string `json:"name"`
	NameInclude *string `json:"name_include"`

	ListDTO
}
