package model

// UserAndDataScope 用户和数据范围的中间表
type UserAndDataScope struct {
	BasicModel
	UserID      int64 `gorm:"nut null;"`
	DataScopeID int64 `gorm:"nut null;"`
}

// TableName 修改表名
func (*UserAndDataScope) TableName() string {
	return "user_and_data_scope"
}
