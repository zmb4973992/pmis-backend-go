package dto

type LoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserGetDTO struct {
	Username          string   `json:"username"`            //用户名
	FullName          *string  `json:"full_name"`           //全名
	EmailAddress      *string  `json:"email_address"`       //邮箱地址
	IsValid           *bool    `json:"is_valid"`            //是否有效
	MobilePhoneNumber *string  `json:"mobile_phone_number"` //手机号
	EmployeeNumber    *string  `json:"employee_number"`     //工号
	Roles             []string `json:"roles"`               //角色
	Departments       []string `json:"departments"`         //部门
}

type UserCreateDTO struct {
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

type UserUpdateDTO struct {
	ID                int     `json:"id"`
	FullName          *string `json:"full_name"  binding:"required"`           //全名
	EmailAddress      *string `json:"email_address" binding:"required"`        //邮箱地址
	IsValid           *bool   `json:"is_valid" binding:"required"`             //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"  binding:"required"` //手机号
	EmployeeNumber    *string `json:"employee_number" binding:"required"`      //工号
	Roles             []int   `json:"roles" binding:"required"`                //角色
	Departments       []int   `json:"departments" binding:"required"`          //部门
}

type UserListDTO struct {
	ID    int  `json:"id"`
	IDGte *int `json:"id_gte"`
	IDLte *int `json:"id_lte"`

	Username        *string `json:"username"`
	UsernameInclude *string `json:"username_include"`

	ListDTO
}
