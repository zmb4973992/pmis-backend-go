package model

type Button struct {
	BasicModel
	Name string //名称
	Sort *int   //排序值
	//连接其他表的id
	MenuID *int //菜单id
	//连接dictionary_item表的id

	//日期

	//数字(允许为0、nil)

	//数字(不允许为0、nil，必须有值)，暂无

	//字符串(允许为nil)

	//字符串(不允许为nil，必须有值)，暂无
}

// TableName 修改表名
func (*Button) TableName() string {
	return "button"
}
