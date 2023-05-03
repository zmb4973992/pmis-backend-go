package model

// UserAndRole 角色和用户的中间表
type UserAndRole struct {
	BasicModel
	RoleID int
	UserID int
}

// TableName 修改表名
func (*UserAndRole) TableName() string {
	return "user_and_role"
}
