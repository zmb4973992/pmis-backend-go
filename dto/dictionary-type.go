package dto

//以下为入参格式

type DictionaryTypeCreateOrUpdateDTO struct {
	BaseDTO
	Name    string  `json:"name" binding:"required"`    //名称
	Sort    *int    `json:"sort" binding:"required"`    //顺序值
	Remarks *string `json:"remarks" binding:"required"` //备注
}

// DictionaryTypeListDTO 是list查询的过滤器
// 在dto传递给sqlCondition时，空值会被忽略
type DictionaryTypeListDTO struct {
	ListDTO
	Name        *string `json:"name"`
	NameInclude string  `json:"name_include"`
	SortGte     *int    `json:"sort_gte"`
}

//以下为出参格式

type DictionaryTypeOutputDTO struct {
	//BaseDTO `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Name    string  `json:"name"  mapstructure:"name"`      //名称
	Sort    *int    `json:"sort" mapstructure:"sort"`       //顺序值
	Remarks *string `json:"remarks" mapstructure:"remarks"` //备注
}
