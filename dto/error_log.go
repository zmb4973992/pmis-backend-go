package dto

type ErrorLogGetDTO struct {
	Detail        *string `json:"detail" mapstructure:"detail"`
	Date          *string `json:"date" mapstructure:"date"`
	MajorCategory *string `json:"major_category" mapstructure:"major_category"`
	MinorCategory *string `json:"minor_category" mapstructure:"minor_category"`
	IsResolved    *bool   `json:"is_resolved" mapstructure:"is_resolved"`
}

// ErrorLogCreateOrUpdateDTO
// 除id外，所有字段都设置为必须绑定
type ErrorLogCreateOrUpdateDTO struct {
	BaseDTO
	ID            int     `json:"id"`
	Detail        *string `json:"detail" mapstructure:"detail" binding:"required"`
	Date          *string `json:"date" mapstructure:"date" binding:"required"`
	MajorCategory *string `json:"major_category" mapstructure:"major_category" binding:"required"`
	MinorCategory *string `json:"minor_category" mapstructure:"minor_category" binding:"required"`
	IsResolved    *bool   `json:"is_resolved" mapstructure:"is_resolved" binding:"required"`
}

// ErrorLogListDTO 是list查询的过滤器
// 在dto传递给sqlCondition时，空值会被忽略
type ErrorLogListDTO struct {
	ListDTO
	Detail        *string `json:"detail" mapstructure:"detail"`
	Date          *string `json:"date" mapstructure:"date"`
	MajorCategory *string `json:"major_category" mapstructure:"major_category"`
	MinorCategory *string `json:"minor_category" mapstructure:"minor_category"`
	IsResolved    *bool   `json:"is_resolved" mapstructure:"is_resolved"`
}
