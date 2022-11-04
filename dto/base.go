package dto

type BaseDTO struct {
	Creator      *int `json:"creator"`
	LastModifier *int `json:"last_modifier"`
	ID           int  `json:"id"`
}
