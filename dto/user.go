package dto

//以下为入参

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserCreate struct {
	Base
	Username          string  `json:"username" binding:"required"`
	Password          string  `json:"password" binding:"required"`
	FullName          *string `json:"full_name"  binding:"required"`           //全名
	EmailAddress      *string `json:"email_address" binding:"required"`        //邮箱地址
	IsValid           *bool   `json:"is_valid" binding:"required"`             //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"  binding:"required"` //手机号
	EmployeeNumber    *string `json:"employee_number" binding:"required"`      //工号
}

type UserUpdate struct {
	Base
	FullName          *string `json:"full_name"  binding:"required"`           //全名
	EmailAddress      *string `json:"email_address" binding:"required"`        //邮箱地址
	IsValid           *bool   `json:"is_valid" binding:"required"`             //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"  binding:"required"` //手机号
	EmployeeNumber    *string `json:"employee_number" binding:"required"`      //工号
}

type UserList struct {
	ListInput
	IDGte        *int    `json:"id_gte"`
	IDLte        *int    `json:"id_lte"`
	IsValid      *bool   `json:"is_valid"`
	UsernameLike *string `json:"username_like"`
}

// 以下为出参

type UserOutput struct {
	Base              `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Username          string                   `json:"username" mapstructure:"username"`                       //用户名
	FullName          *string                  `json:"full_name" mapstructure:"full_name"`                     //全名
	EmailAddress      *string                  `json:"email_address" mapstructure:"email_address"`             //邮箱地址
	IsValid           *bool                    `json:"is_valid" mapstructure:"is_valid"`                       //是否有效
	MobilePhoneNumber *string                  `json:"mobile_phone_number" mapstructure:"mobile_phone_number"` //手机号
	EmployeeNumber    *string                  `json:"employee_number" mapstructure:"employee_number"`         //工号
}
