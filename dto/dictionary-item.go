package dto

// DictionaryItemCreateOrUpdateDTO
// 除id外，所有字段都设置为必须绑定
type DictionaryItemCreateOrUpdateDTO struct {
	BaseDTO
	DictionaryTypeID int     `json:"dictionary_type_id" binding:"required"` //字典类型id
	Name             string  `json:"name" binding:"required"`               //名称
	Sort             *int    `json:"sort" binding:"required"`               //顺序值
	Remarks          *string `json:"remarks" binding:"required"`            //备注
}

type DictionaryItemOutputDTO struct {
	//BaseDTO `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Name    string  `json:"name"  mapstructure:"name"`      //名称
	Sort    *int    `json:"sort" mapstructure:"sort"`       //顺序值
	Remarks *string `json:"remarks" mapstructure:"remarks"` //备注
}
