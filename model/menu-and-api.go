package model

// MenuAndApi 菜单和api的中间表
type MenuAndApi struct {
	BasicModel
	MenuSnowID int64 `gorm:"nut null;"`
	ApiSnowID  int64 `gorm:"nut null;"`
}

// TableName 修改表名
func (*MenuAndApi) TableName() string {
	return "menu_and_api"
}
