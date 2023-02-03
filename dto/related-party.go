package dto

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type RelatedPartyCreate struct {
	Creator      int
	LastModifier int

	ChineseName             string `json:"chinese_name,omitempty"`
	EnglishName             string `json:"english_name,omitempty"`
	Address                 string `json:"address,omitempty"`
	UniformSocialCreditCode string `json:"uniform_social_credit_code,omitempty"` //统一社会信用代码
	Telephone               string `json:"telephone,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type RelatedPartyUpdate struct {
	LastModifier int
	ID           int

	ChineseName             *string `json:"chinese_name"`
	EnglishName             *string `json:"english_name"`
	Address                 *string `json:"address"`
	UniformSocialCreditCode *string `json:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string `json:"telephone"`
}

type RelatedPartyDelete struct {
	Deleter int
	ID      int
}

type RelatedPartyList struct {
	ListInput

	ChineseNameInclude string `json:"chinese_name_include,omitempty"`
	EnglishNameInclude string `json:"english_name_include,omitempty"`
}

//以下为出参

type RelatedPartyOutput struct {
	Creator      *int    `json:"creator" gorm:"creator"`
	LastModifier *int    `json:"last_modifier" gorm:"last_modifier"`
	ID           int     `json:"id" gorm:"id"`
	Name         string  `json:"name" gorm:"name"`       //名称
	Sort         *int    `json:"sort" gorm:"sort"`       //顺序值
	Remarks      *string `json:"remarks" gorm:"remarks"` //备注
}
