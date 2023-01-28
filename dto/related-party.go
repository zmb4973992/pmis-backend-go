package dto

//以下为入参

type RelatedPartyCreateOrUpdate struct {
	Base
	ChineseName             *string `json:"chinese_name" binding:"required"`
	EnglishName             *string `json:"english_name" binding:"required"`
	Address                 *string `json:"address" binding:"required"`
	UniformSocialCreditCode *string `json:"uniform_social_credit_code" binding:"required"` //统一社会信用代码
	Telephone               *string `json:"telephone" binding:"required"`
}

type RelatedPartyList struct {
	ListInput
	ID              int     `form:"id"`
	IDGte           *int    `form:"id_gte"`
	IDLte           *int    `form:"id_lte"`
	ChineseName     *string `form:"chinese_name"`
	ChineseNameLike *string `form:"chinese_name_like"`
	EnglishNameLike *string `form:"english_name_like"`
}

//以下为出参

type RelatedPartyOutput struct {
	Base                    `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	ChineseName             *string                  `json:"chinese_name" mapstructure:"chinese_name"`
	EnglishName             *string                  `json:"english_name" mapstructure:"english_name"`
	Address                 *string                  `json:"address" mapstructure:"address"`
	UniformSocialCreditCode *string                  `json:"uniform_social_credit_code" mapstructure:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string                  `json:"telephone" mapstructure:"telephone"`
}
