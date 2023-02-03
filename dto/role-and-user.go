package dto

//以下为入参

type RoleAndUserCreateOrUpdate struct {
	Creator      *int `json:"creator" mapstructure:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" mapstructure:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" mapstructure:"id" gorm:"id"`

	RoleIDs []int `json:"role_ids"`
	UserIDs []int `json:"user_ids"`
}

//以下为出参
