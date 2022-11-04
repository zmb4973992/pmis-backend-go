package dto

type ProjectGetDTO struct {
	BaseDTO
	ProjectCode      *string  `json:"project_code" mapstructure:"project_code"`
	ProjectFullName  *string  `json:"project_full_name" mapstructure:"project_full_name"`
	ProjectShortName *string  `json:"project_short_name" mapstructure:"project_short_name"`
	Country          *string  `json:"country" mapstructure:"country"`
	Province         *string  `json:"province" mapstructure:"province"`
	ProjectType      *string  `json:"project_type" mapstructure:"project_type"`
	Department       *string  `json:"department" mapstructure:"department"`
	Amount           *float64 `json:"amount" mapstructure:"amount"`
	Currency         *string  `json:"currency" mapstructure:"currency"`
	ExchangeRate     *float64 `json:"exchange_rate" mapstructure:"exchange_rate"`
	RelatedPartyID   *int     `json:"related_party_id" mapstructure:"related_party_id"`
}

// ProjectCreateOrUpdateDTO
// 除id外，所有字段都设置为必须绑定
type ProjectCreateOrUpdateDTO struct {
	BaseDTO
	ProjectCode      *string  `json:"project_code" binding:"required"`
	ProjectFullName  *string  `json:"project_full_name" mapstructure:"project_full_name"`
	ProjectShortName *string  `json:"project_short_name" binding:"required"`
	Country          *string  `json:"country" binding:"required"`
	Province         *string  `json:"province" binding:"required"`
	ProjectType      *string  `json:"project_type" binding:"required"`
	Department       *string  `json:"department" binding:"required"`
	Amount           *float64 `json:"amount" binding:"required"`
	Currency         *string  `json:"currency" binding:"required"`
	ExchangeRate     *float64 `json:"exchange_rate" binding:"required"`
	RelatedPartyID   *int     `json:"related_party_id" binding:"required"`
}

// ProjectListDTO 是list查询的过滤器
// 在dto传递给sqlCondition时，空值会被忽略
type ProjectListDTO struct {
	ListDTO

	ProjectID               *int    `form:"project_id"`
	RelatedPartyNameInclude *string `form:"related_party_name_include"`
	DepartmentInclude       *string `form:"department_include"`
	ProjectFullNameInclude  *string `form:"project_full_name_include"`
}
