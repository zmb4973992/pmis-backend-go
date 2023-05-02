package model

// MenuAndApi 菜单和api的中间表
type MenuAndApi struct {
	BasicModel
	MenuID *int
	ApiID  *int
}

// TableName 修改表名
func (*MenuAndApi) TableName() string {
	return "menu_and_api"
}
