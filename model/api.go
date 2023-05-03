package model

type Api struct {
	BasicModel
	Group     string `gorm:"not null;"`
	Name      string `gorm:"not null;"`
	Path      string `gorm:"not null;"`
	Method    string `gorm:"not null;"`
	WithParam bool   `gorm:"not null;"`
	//连接其他表的id

	//连接dictionary_item表的id

	//日期

	//数字(允许为0、nil)

	//数字(不允许为0、nil，必须有值)，暂无

	//字符串(允许为nil)

	//字符串(不允许为nil，必须有值)，暂无
}

// TableName 修改表名
func (*Api) TableName() string {
	return "api"
}
