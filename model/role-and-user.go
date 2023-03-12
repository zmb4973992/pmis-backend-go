package model

// RoleAndUser 角色和用户的中间表
type RoleAndUser struct {
	BasicModel
	RoleID *int
	UserID *int
}

// TableName 修改表名
func (*RoleAndUser) TableName() string {
	return "role_and_user"
}
