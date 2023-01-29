package dto

//以下为入参
//指针的话允许不传参（即json格式下值为null）
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type DictionaryTypeCreate struct {
	Creator      int
	LastModifier int
	Name         string  `json:"name" binding:"required"` //名称
	Sort         *int    `json:"sort"`                    //顺序值
	Remarks      *string `json:"remarks"`                 //备注
}

type DictionaryTypeUpdate struct {
	LastModifier int
	ID           int
	Name         string  `json:"name" binding:"required"` //名称
	Sort         *int    `json:"sort"`                    //顺序值
	Remarks      *string `json:"remarks"`                 //备注
}

type DictionaryTypeDelete struct {
	Deleter int
	ID      int
}

type DictionaryTypeList struct {
	ListInput
	Name        *string `json:"name"`
	NameInclude *string `json:"name_include"`
	SortGte     *int    `json:"sort_gte"`
}

//以下为出参

type DictionaryTypeOutput struct {
	Name    string  `json:"name"  mapstructure:"name"`      //名称
	Sort    *int    `json:"sort" mapstructure:"sort"`       //顺序值
	Remarks *string `json:"remarks" mapstructure:"remarks"` //备注
}
