package model

type File struct {
	BasicModel
	Name   string
	SizeMB float64 // 文件大小(MB)
}

// TableName 修改数据库的表名
func (f *File) TableName() string {
	return "file"
}
