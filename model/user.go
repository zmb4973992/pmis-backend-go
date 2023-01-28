package model

type User struct {
	BaseModel
	Username          string
	Password          string
	IsValid           *bool   //用户为有效还是禁用
	FullName          *string //全名
	EmailAddress      *string //邮箱地址
	MobilePhoneNumber *string //手机号
	EmployeeNumber    *string //工号
	//这里是声名外键关系，并不是实际字段。不建议用gorm的多对多的设定，不好修改
	Roles           []RoleAndUser       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Departments     []DepartmentAndUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProjectAndUsers []ProjectAndUser    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OperationLogs   []OperationLog      `gorm:"constraint:OnUpdate:CASCADE;"`
}

// TableName 将表名改为user
func (User) TableName() string {
	return "user"
}
