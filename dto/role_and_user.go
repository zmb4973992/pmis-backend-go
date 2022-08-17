package dto

// RoleAndUserGetDTO mapstructure用于list
// list方法中返回的数据为[]map，需要借助mapstructure转换为struct，再返回给前端
// []map中的键为数据库的字段名，mapstructure需要和[]map中的键名保持一致
type RoleAndUserGetDTO struct {
	RoleID *int `json:"role_id"`
	UserID *int `json:"user_id"`
}

type RoleAndUserCreateDTO struct {
	Username          string  `json:"username" binding:"required"`
	Password          string  `json:"password" binding:"required"`
	FullName          *string `json:"full_name"  binding:"required"`           //全名
	EmailAddress      *string `json:"email_address" binding:"required"`        //邮箱地址
	IsValid           *bool   `json:"is_valid" binding:"required"`             //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"  binding:"required"` //手机号
	EmployeeNumber    *string `json:"employee_number" binding:"required"`      //工号
	Roles             []int   `json:"roles" binding:"required"`                //角色
	Departments       []int   `json:"departments" binding:"required"`          //部门
}

type RoleAndUserUpdateDTO struct {
	ID                int     `json:"id"`
	FullName          *string `json:"full_name"  binding:"required"`           //全名
	EmailAddress      *string `json:"email_address" binding:"required"`        //邮箱地址
	IsValid           *bool   `json:"is_valid" binding:"required"`             //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"  binding:"required"` //手机号
	EmployeeNumber    *string `json:"employee_number" binding:"required"`      //工号
	Roles             []int   `json:"roles" binding:"required"`                //角色
	Departments       []int   `json:"departments" binding:"required"`          //部门
}

type RoleAndUserListDTO struct {
	RoleID *int `json:"role_id"`
	UserID *int `json:"user_id"`
}
