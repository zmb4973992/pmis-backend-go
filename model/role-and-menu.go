package model

// MenuAndApi 角色和菜单的中间表
type RoleAndMenu struct {
	BasicModel
	RoleSnowID uint64 `gorm:"nut null;"`
	MenuSnowID uint64 `gorm:"nut null;"`
}

// TableName 修改表名
func (*RoleAndMenu) TableName() string {
	return "role_and_menu"
}
