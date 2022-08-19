package dto

type RoleAndUserCreateOrUpdateDTO struct {
	RoleIDs []int `json:"role_ids"`
	UserIDs []int `json:"user_ids"`
}
