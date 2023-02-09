package model

type File struct {
	BaseModel
	UUID            string
	InitialFileName string
	StoredFileName  string
	StoragePath     string
	Size            int // 文件大小(KB)
}

// TableName 修改数据库的表名
func (*File) TableName() string {
	return "file"
}
