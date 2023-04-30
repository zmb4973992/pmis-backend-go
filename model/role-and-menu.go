package model

// RoleAndMenu 角色和菜单的中间表
type RoleAndMenu struct {
	BasicModel
	RoleID *int
	MenuID *int
}

// TableName 修改表名
func (*RoleAndMenu) TableName() string {
	return "role_and_menu"
}
