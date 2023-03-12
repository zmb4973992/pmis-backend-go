package model

import "time"

// ProjectAndUser 组织机构和用户的中间表
type ProjectAndUser struct {
	BasicModel
	ProjectID *int
	UserID    *int
	Title     *string
	StartDate *time.Time `gorm:"type:date"`
	EndDate   *time.Time `gorm:"type:date"`
}

// TableName 修改表名
func (*ProjectAndUser) TableName() string {
	return "project_and_user"
}
