package model

// UserAndRole 角色和用户的中间表
type UserAndRole struct {
	BasicModel
	RoleSnowID int64 `gorm:"nut null;"`
	UserSnowID int64 `gorm:"nut null;"`
}

// TableName 修改表名
func (*UserAndRole) TableName() string {
	return "user_and_role"
}
