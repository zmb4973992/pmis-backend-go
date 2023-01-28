package dto

//以下为入参

type DictionaryItemCreateOrUpdate struct {
	Base
	DictionaryTypeID int     `json:"dictionary_type_id" binding:"required"` //字典类型id
	Name             string  `json:"name" binding:"required"`               //名称
	Sort             *int    `json:"sort" binding:"required"`               //顺序值
	Remarks          *string `json:"remarks" binding:"required"`            //备注
}

//以下为出参

type DictionaryItemOutput struct {
	//Base `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Name    string  `json:"name"  mapstructure:"name"`      //名称
	Sort    *int    `json:"sort" mapstructure:"sort"`       //顺序值
	Remarks *string `json:"remarks" mapstructure:"remarks"` //备注
}
