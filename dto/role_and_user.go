package dto

// RoleAndUserGetDTO mapstructure用于list
// list方法中返回的数据为[]map，需要借助mapstructure转换为struct，再返回给前端
// []map中的键为数据库的字段名，mapstructure需要和[]map中的键名保持一致
type RoleAndUserGetDTO struct {
	RoleID *int `json:"role_id" mapstructure:"role_id"`
	UserID *int `json:"user_id" mapstructure:"user_id"`
}

//这里不需要update的dto，因为update拆分成先delete、再create

type RoleAndUserCreateDTO struct {
	RoleID *int `json:"role_id" binding:"required"`
	UserID *int `json:"user_id" binding:"required"`
}

type RoleAndUserCreateOrUpdateDTO struct {
	RoleIDs []int `json:"role_ids"`
	UserIDs []int `json:"user_ids"`
}

type RoleAndUserDeleteDTO struct {
	RoleID *int `json:"role_id"`
	UserID *int `json:"user_id"`
}

type RoleAndUserListDTO struct {
	RoleID *int `json:"role_id"`
	UserID *int `json:"user_id"`

	ListDTO
}
