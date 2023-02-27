package model

type File struct {
	BaseModel
	UUID string
	Name string
	Mode string
	Path string
	Size int // 文件大小(KB)
}

// TableName 修改数据库的表名
func (*File) TableName() string {
	return "file"
}
