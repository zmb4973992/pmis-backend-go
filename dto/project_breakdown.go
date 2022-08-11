package dto

type ProjectBreakdownGetDTO struct {
	Name       *string  `json:"name"`        //名称
	ProjectID  *int     `json:"project_id"`  //所属项目id
	Level      *int     `json:"level"`       //层级
	Weight     *float64 `json:"weight"`      //权重
	SuperiorID *int     `json:"superior_id"` //上级拆解项id
}

type ProjectBreakdownCreateAndUpdateDTO struct {
	ID         int      `json:"id"`
	Name       *string  `json:"name" binding:"required"`        //拆解项名称
	ProjectID  *int     `json:"project_id" binding:"required"`  //所属项目id
	Level      *int     `json:"level" binding:"required"`       //层级
	Weight     *float64 `json:"weight" binding:"required"`      //权重
	SuperiorID *int     `json:"superior_id" binding:"required"` //上级拆解项ID
}

// ProjectBreakdownListDTO 是list查询的过滤器
// 这里不用指针，如果前端没传字段或者只传字段没传值，那么该字段默认为空
// 在dto传递给sqlCondition时，空值会被忽略
type ProjectBreakdownListDTO struct {
	ID    int  `json:"id"`
	IDGte *int `json:"id_gte"`
	IDLte *int `json:"id_lte"`

	Name        *string `json:"name"`
	NameInclude *string `json:"name_include"`

	ListDTO
}
