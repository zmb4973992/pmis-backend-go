package model

// OrganizationAndUser 组织和用户的中间表
type OrganizationAndUser struct {
	BasicModel
	OrganizationSnowID int64 `gorm:"nut null;"`
	UserSnowID         int64 `gorm:"nut null;"`
	ImportedByLDAP     *bool
}

// TableName 修改表名
func (*OrganizationAndUser) TableName() string {
	return "organization_and_user"
}
