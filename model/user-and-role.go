package model

// UserAndRole 角色和用户的中间表
type UserAndRole struct {
	BasicModel
	RoleId int64 `gorm:"nut null;"`
	UserId int64 `gorm:"nut null;"`
}

// TableName 修改表名
func (u *UserAndRole) TableName() string {
	return "user_and_role"
}
