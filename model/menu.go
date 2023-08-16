package model

type Menu struct {
	BasicModel
	SuperiorID *int64
	Path       *string //路由访问路径
	Group      string
	Name       string  //名称
	Component  *string //组件路径
	Sort       *int    //排序值
	Meta
	//连接其他表的id

	//连接dictionary_item表的id

	//日期

	//数字(允许为0、nil)

	//数字(不允许为0、nil，必须有值)，暂无

	//字符串(允许为nil)

	//字符串(不允许为nil，必须有值)，暂无
}

type Meta struct {
	Hidden    bool    //在侧边栏内是否隐藏
	KeepAlive *bool   //是否缓存
	Title     *string //菜单名
	Icon      *string //图标
}

// TableName 修改表名
func (m *Menu) TableName() string {
	return "menu"
}
