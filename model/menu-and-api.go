package model

// MenuAndApi 菜单和api的中间表
type MenuAndApi struct {
	BasicModel
	MenuId int64 `gorm:"nut null;"`
	ApiId  int64 `gorm:"nut null;"`
}

// TableName 修改表名
func (m *MenuAndApi) TableName() string {
	return "menu_and_api"
}
