package dto

type LoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserCreateDTO struct {
	BaseDTO
	Username          string  `json:"username" binding:"required"`
	Password          string  `json:"password" binding:"required"`
	FullName          *string `json:"full_name"  binding:"required"`           //全名
	EmailAddress      *string `json:"email_address" binding:"required"`        //邮箱地址
	IsValid           *bool   `json:"is_valid" binding:"required"`             //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"  binding:"required"` //手机号
	EmployeeNumber    *string `json:"employee_number" binding:"required"`      //工号
}

type UserUpdateDTO struct {
	BaseDTO
	ID                int     `json:"id"`
	FullName          *string `json:"full_name"  binding:"required"`           //全名
	EmailAddress      *string `json:"email_address" binding:"required"`        //邮箱地址
	IsValid           *bool   `json:"is_valid" binding:"required"`             //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"  binding:"required"` //手机号
	EmployeeNumber    *string `json:"employee_number" binding:"required"`      //工号
}

// UserGetDTO mapstructure用于list
// list方法中返回的数据为[]map，需要借助mapstructure转换为struct，再返回给前端
// []map中的键为数据库的字段名，mapstructure需要和[]map中的键名保持一致
type UserGetDTO struct {
	Username          string  `json:"username" mapstructure:"username"`                       //用户名
	FullName          *string `json:"full_name" mapstructure:"full_name"`                     //全名
	EmailAddress      *string `json:"email_address" mapstructure:"email_address"`             //邮箱地址
	IsValid           *bool   `json:"is_valid" mapstructure:"is_valid"`                       //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number" mapstructure:"mobile_phone_number"` //手机号
	EmployeeNumber    *string `json:"employee_number" mapstructure:"employee_number"`         //工号
	//Roles             []string `json:"roles" mapstructure:"-"`                                 //角色
	//Departments       []string `json:"departments" mapstructure:"-"`                           //部门
}

// UserListDTO form用于接收query传参
type UserListDTO struct {
	ListDTO
	IDGte           *int    `form:"id_gte"`
	IDLte           *int    `form:"id_lte"`
	IsValid         *int    `form:"is_valid"`
	UsernameInclude *string `form:"username_include"`
}
