package model

// ProjectAndUser 组织机构和用户的中间表
type ProjectAndUser struct {
	BaseModel
	ProjectID *int
	UserID    *int
	Title     *string
	StartDate *string `gorm:"type:date;"`
	EndDate   *string `gorm:"type:date;"`
}

// TableName 修改表名
func (ProjectAndUser) TableName() string {
	return "project_and_user"
}
