package model

// MenuAndApi 角色和菜单的中间表

type RoleAndMenu struct {
	BasicModel
	RoleID int64 `gorm:"nut null;"`
	MenuID int64 `gorm:"nut null;"`
}

// TableName 修改表名
func (r *RoleAndMenu) TableName() string {
	return "role_and_menu"
}
