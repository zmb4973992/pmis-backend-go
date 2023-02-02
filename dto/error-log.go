package dto

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type ErrorLogCreate struct {
	Creator      int
	LastModifier int

	Detail        string `json:"detail,omitempty" `
	Date          string `json:"date,omitempty"`
	MajorCategory string `json:"major_category,omitempty"`
	MinorCategory string `json:"minor_category,omitempty"`
	IsResolved    bool   `json:"is_resolved,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ErrorLogUpdate struct {
	LastModifier int
	ID           int

	Detail        *string `json:"detail"`
	Date          *string `json:"date"`
	MajorCategory *string `json:"major_category"`
	MinorCategory *string `json:"minor_category"`
	IsResolved    *bool   `json:"is_resolved"`
}

type ErrorLogDelete struct {
	Deleter int
	ID      int
}

type ErrorLogList struct {
	ListInput

	DetailInclude string `json:"detail_include,omitempty" `
	Date          string `json:"date,omitempty"`
	MajorCategory string `json:"major_category,omitempty"`
	MinorCategory string `json:"minor_category,omitempty"`
	IsResolved    bool   `json:"is_resolved,omitempty"`
}

//以下为出参

type ErrorLogOutput struct {
	Creator      *int `json:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" gorm:"id"`

	Detail        *string `json:"detail" gorm:"detail"`
	Date          *string `json:"date" gorm:"date"`
	MajorCategory *string `json:"major_category" gorm:"major_category"`
	MinorCategory *string `json:"minor_category" gorm:"minor_category"`
	IsResolved    *bool   `json:"is_resolved" gorm:"is_resolved"`
}
