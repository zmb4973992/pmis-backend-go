package dto

type ProjectGetDTO struct {
	BaseDTO          `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	ProjectCode      *string                  `json:"project_code" mapstructure:"project_code"`
	ProjectFullName  *string                  `json:"project_full_name" mapstructure:"project_full_name"`
	ProjectShortName *string                  `json:"project_short_name" mapstructure:"project_short_name"`
	Country          *string                  `json:"country" mapstructure:"country"`
	Province         *string                  `json:"province" mapstructure:"province"`
	ProjectType      *string                  `json:"project_type" mapstructure:"project_type"`
	Amount           *float64                 `json:"amount" mapstructure:"amount"`
	Currency         *string                  `json:"currency" mapstructure:"currency"`
	ExchangeRate     *float64                 `json:"exchange_rate" mapstructure:"exchange_rate"`
	RelatedPartyID   *int                     `json:"related_party_id" mapstructure:"related_party_id"`
	DepartmentID     *int                     `json:"-" mapstructure:"department_id"`
	Department       *DepartmentGetDTO        `json:"department" gorm:"-"`
}

// ProjectCreateOrUpdateDTO
// 除id外，所有字段都设置为必须绑定
type ProjectCreateOrUpdateDTO struct {
	BaseDTO
	ProjectCode      *string  `json:"project_code" binding:"required"`
	ProjectFullName  *string  `json:"project_full_name" binding:"required"`
	ProjectShortName *string  `json:"project_short_name" binding:"required"`
	Country          *string  `json:"country" binding:"required"`
	Province         *string  `json:"province" binding:"required"`
	ProjectType      *string  `json:"project_type" binding:"required"`
	Amount           *float64 `json:"amount" binding:"required"`
	Currency         *string  `json:"currency" binding:"required"`
	ExchangeRate     *float64 `json:"exchange_rate" binding:"required"`
	DepartmentID     *int     `json:"department_id" binding:"required"`
	RelatedPartyID   *int     `json:"related_party_id" binding:"required"`
}

// ProjectListDTO 是list查询的过滤器
// 在dto传递给sqlCondition时，空值会被忽略
type ProjectListDTO struct {
	ListDTO
	ProjectNameLike    *string `json:"project_name_like"` //包含项目全称和项目简称
	DepartmentNameLike *string `json:"department_name_like"`
	DepartmentIDIn     []int   `json:"department_id_in"`
}
