package dto

type RoleAndUserCreateOrUpdateDTO struct {
	BaseDTO
	RoleIDs []int `json:"role_ids"`
	UserIDs []int `json:"user_ids"`
}
