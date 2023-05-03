package model

type Api struct {
	BasicModel
	SnowID uint64
	Group  *string
	Name   *string
	Path   *string
	Method *string
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
