package model

// OrganizationAndUser 组织和用户的中间表
type OrganizationAndUser struct {
	BasicModel
	OrganizationID *int
	UserID         *int
}

// TableName 修改表名
func (*OrganizationAndUser) TableName() string {
	return "organization_and_user"
}
