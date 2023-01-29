package dto

type Base struct {
	Creator      *int `json:"creator" mapstructure:"creator"`
	LastModifier *int `json:"last_modifier" mapstructure:"last_modifier"`
	ID           int  `json:"id" mapstructure:"id"`
}
