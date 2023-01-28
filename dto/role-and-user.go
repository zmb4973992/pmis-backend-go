package dto

//以下为入参

type RoleAndUserCreateOrUpdate struct {
	Base
	RoleIDs []int `json:"role_ids"`
	UserIDs []int `json:"user_ids"`
}

//以下为出参
