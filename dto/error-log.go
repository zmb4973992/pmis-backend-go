package dto

//以下为入参

type ErrorLogCreateOrUpdate struct {
	Base
	Detail        *string `json:"detail" binding:"required"`
	Date          *string `json:"date" binding:"required"`
	MajorCategory *string `json:"major_category" binding:"required"`
	MinorCategory *string `json:"minor_category" binding:"required"`
	IsResolved    *bool   `json:"is_resolved" binding:"required"`
}

type ErrorLogList struct {
	ListInput
	Detail        *string `json:"detail"`
	Date          *string `json:"date"`
	MajorCategory *string `json:"major_category"`
	MinorCategory *string `json:"minor_category"`
	IsResolved    *bool   `json:"is_resolved"`
}

//以下为出参

type ErrorLogOutput struct {
	Base          `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	Detail        *string                  `json:"detail" mapstructure:"detail"`
	Date          *string                  `json:"date" mapstructure:"date"`
	MajorCategory *string                  `json:"major_category" mapstructure:"major_category"`
	MinorCategory *string                  `json:"minor_category" mapstructure:"minor_category"`
	IsResolved    *bool                    `json:"is_resolved" mapstructure:"is_resolved"`
}
