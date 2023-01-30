package dto

type Base struct {
	//mapstructure不能删，用于service list中从model到dto的转换
	Creator      *int `json:"creator" mapstructure:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" mapstructure:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" mapstructure:"id" gorm:"id"`
}
