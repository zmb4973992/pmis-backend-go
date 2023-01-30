package dto

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type DictionaryTypeCreate struct {
	Creator      int
	LastModifier int
	Name         string `json:"name" binding:"required"` //名称
	Sort         int    `json:"sort,omitempty"`          //顺序值
	Remarks      string `json:"remarks,omitempty"`       //备注
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DictionaryTypeUpdate struct {
	LastModifier int
	ID           int
	Name         *string `json:"name"`    //名称
	Sort         *int    `json:"sort"`    //顺序值
	Remarks      *string `json:"remarks"` //备注
}

type DictionaryTypeDelete struct {
	Deleter int
	ID      int
}

type DictionaryTypeList struct {
	ListInput
	NameInclude string `json:"name_include,omitempty"`
}

//以下为出参

type DictionaryTypeOutput struct {
	Creator      *int    `json:"creator" gorm:"creator"`
	LastModifier *int    `json:"last_modifier" gorm:"last_modifier"`
	ID           int     `json:"id" gorm:"id"`
	Name         string  `json:"name" gorm:"name"`       //名称
	Sort         *int    `json:"sort" gorm:"sort"`       //顺序值
	Remarks      *string `json:"remarks" gorm:"remarks"` //备注
}
