package dto

type DictionaryTypeCreateOrUpdateDTO struct {
	BaseDTO
	Name    string  `json:"name" binding:"required"`    //名称
	Sort    *int    `json:"sort" binding:"required"`    //顺序值
	Remarks *string `json:"remarks" binding:"required"` //备注
}

type DictionaryTypeOutputDTO struct {
	//BaseDTO `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Name    string  `json:"name"  mapstructure:"name"`      //名称
	Sort    *int    `json:"sort" mapstructure:"sort"`       //顺序值
	Remarks *string `json:"remarks" mapstructure:"remarks"` //备注
}
