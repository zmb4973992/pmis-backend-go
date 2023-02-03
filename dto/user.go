package dto

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserCreate struct {
	Creator           int
	LastModifier      int
	Username          string `json:"username" binding:"required"`
	Password          string `json:"password" binding:"required"`
	FullName          string `json:"full_name,omitempty"`           //全名
	EmailAddress      string `json:"email_address,omitempty"`       //邮箱地址
	IsValid           *bool  `json:"is_valid"`                      //是否有效
	MobilePhoneNumber string `json:"mobile_phone_number,omitempty"` //手机号
	EmployeeNumber    string `json:"employee_number,omitempty"`     //工号
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type UserUpdate struct {
	LastModifier      int
	ID                int
	FullName          *string `json:"full_name"`           //全名
	EmailAddress      *string `json:"email_address"`       //邮箱地址
	IsValid           *bool   `json:"is_valid"`            //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"` //手机号
	EmployeeNumber    *string `json:"employee_number"`     //工号
}

type UserDelete struct {
	Deleter int
	ID      int
}

type UserList struct {
	ListInput
	IsValid         *bool  `json:"is_valid"`
	UsernameInclude string `json:"username_include,omitempty"`
}

//以下为出参

type UserOutput struct {
	Creator      *int `json:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" gorm:"id"`

	Username          string  `json:"username" gorm:"username"`                       //用户名
	FullName          *string `json:"full_name" gorm:"full_name"`                     //全名
	EmailAddress      *string `json:"email_address" gorm:"email_address"`             //邮箱地址
	IsValid           *bool   `json:"is_valid" gorm:"is_valid"`                       //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number" gorm:"mobile_phone_number"` //手机号
	EmployeeNumber    *string `json:"employee_number" gorm:"employee_number"`         //工号
}
