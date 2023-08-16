package model

// UserAndDataAuthority 用户和数据范围的中间表
type UserAndDataAuthority struct {
	BasicModel
	UserID          int64 `gorm:"nut null;"`
	DataAuthorityID int64 `gorm:"nut null;"`
	ImportedByLDAP  *bool
}

// TableName 修改表名
func (u *UserAndDataAuthority) TableName() string {
	return "user_and_data_authority"
}
