package model

// DepartmentAndUser 组织机构和用户的中间表
type DepartmentAndUser struct {
	BaseModel
	DepartmentID *int
	UserID       *int
}

// TableName 修改表名
func (*DepartmentAndUser) TableName() string {
	return "department_and_user"
}
