package model

// MenuAndApi 角色和菜单的中间表

type RoleAndMenu struct {
	BasicModel
	RoleId int64 `gorm:"nut null;"`
	MenuId int64 `gorm:"nut null;"`
}

// TableName 修改表名
func (r *RoleAndMenu) TableName() string {
	return "role_and_menu"
}
