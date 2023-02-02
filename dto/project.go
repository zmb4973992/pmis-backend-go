package dto

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type ProjectCreate struct {
	Creator          int
	LastModifier     int
	ProjectCode      string   `json:"project_code,omitempty"`
	ProjectFullName  string   `json:"project_full_name,omitempty"`
	ProjectShortName string   `json:"project_short_name,omitempty"`
	Country          string   `json:"country,omitempty"`
	Province         string   `json:"province,omitempty"`
	ProjectType      string   `json:"project_type,omitempty"`
	Amount           *float64 `json:"amount"`
	Currency         string   `json:"currency,omitempty"`
	ExchangeRate     *float64 `json:"exchange_rate,omitempty"`
	DepartmentID     int      `json:"department_id,omitempty"`
	RelatedPartyID   int      `json:"related_party_id,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ProjectUpdate struct {
	LastModifier     int
	ID               int
	ProjectCode      *string  `json:"project_code"`
	ProjectFullName  *string  `json:"project_full_name"`
	ProjectShortName *string  `json:"project_short_name"`
	Country          *string  `json:"country"`
	Province         *string  `json:"province"`
	ProjectType      *string  `json:"project_type"`
	Amount           *float64 `json:"amount"`
	Currency         *string  `json:"currency"`
	ExchangeRate     *float64 `json:"exchange_rate"`
	DepartmentID     *int     `json:"department_id"`
	RelatedPartyID   *int     `json:"related_party_id"`
}

type ProjectDelete struct {
	Deleter int
	ID      int
}

type ProjectList struct {
	ListInput
	AuthInput
	ProjectNameLike    string `json:"project_name_like,omitempty"` //包含项目全称和项目简称
	DepartmentNameLike string `json:"department_name_like,omitempty"`
	DepartmentIDIn     []int  `json:"department_id_in"`
}

//以下为出参

type ProjectOutput struct {
	Creator      *int `json:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" gorm:"id"`

	ProjectCode      *string           `json:"project_code" gorm:"project_code"`
	ProjectFullName  *string           `json:"project_full_name" gorm:"project_full_name"`
	ProjectShortName *string           `json:"project_short_name" gorm:"project_short_name"`
	Country          *string           `json:"country" gorm:"country"`
	Province         *string           `json:"province" gorm:"province"`
	ProjectType      *string           `json:"project_type" gorm:"project_type"`
	Amount           *float64          `json:"amount" gorm:"amount"`
	Currency         *string           `json:"currency" gorm:"currency"`
	ExchangeRate     *float64          `json:"exchange_rate" gorm:"exchange_rate"`
	RelatedPartyID   *int              `json:"related_party_id" gorm:"related_party_id"`
	DepartmentID     *int              `json:"-" gorm:"department_id"`
	Department       *DepartmentOutput `json:"department" gorm:"-"` //gorm -  要不要删除？
}
