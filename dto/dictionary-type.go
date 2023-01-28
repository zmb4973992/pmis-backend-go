package dto

//以下为入参

type DictionaryTypeCreateOrUpdate struct {
	Base
	Name    string  `json:"name" binding:"required"`    //名称
	Sort    *int    `json:"sort" binding:"required"`    //顺序值
	Remarks *string `json:"remarks" binding:"required"` //备注
}

type DictionaryTypeDelete struct {
	Base
	//这里不用json tag，因为不从前端读取，而是在controller中处理
	DictionaryTypeID int
}

type DictionaryTypeList struct {
	ListInput
	Name        *string `json:"name"`
	NameInclude string  `json:"name_include"`
	SortGte     *int    `json:"sort_gte"`
}

//以下为出参

type DictionaryTypeOutput struct {
	//Base `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Name    string  `json:"name"  mapstructure:"name"`      //名称
	Sort    *int    `json:"sort" mapstructure:"sort"`       //顺序值
	Remarks *string `json:"remarks" mapstructure:"remarks"` //备注
}
