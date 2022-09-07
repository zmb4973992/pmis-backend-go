package dto

//RelatedPartyGetDTO dto只接收、发送数据，并不会对数据库进行任何操作
//所有操作数据库的工作会经过dao层的清洗后都交给model来完成
//这里必须是指针类型，因为只有指针才能向前端传递nil
//以后如果接收和推送用不同的dto，可以考虑不用指针
type RelatedPartyGetDTO struct {
	ChineseName             *string `json:"chinese_name" mapstructure:"chinese_name"`
	EnglishName             *string `json:"english_name" mapstructure:"english_name"`
	Address                 *string `json:"address" mapstructure:"address"`
	UniformSocialCreditCode *string `json:"uniform_social_credit_code" mapstructure:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string `json:"telephone" mapstructure:"telephone"`
}

type RelatedPartyCreateOrUpdateDTO struct {
	BaseDTO
	ID                      int     `json:"id"`
	ChineseName             *string `json:"chinese_name" binding:"required"`
	EnglishName             *string `json:"english_name" binding:"required"`
	Address                 *string `json:"address" binding:"required"`
	UniformSocialCreditCode *string `json:"uniform_social_credit_code" binding:"required"` //统一社会信用代码
	Telephone               *string `json:"telephone" binding:"required"`
}

// RelatedPartyListDTO 是list查询的过滤器
// 在dto传递给sqlCondition时，空值会被忽略
type RelatedPartyListDTO struct {
	ListDTO
	ID    int  `form:"id"`
	IDGte *int `form:"id_gte"`
	IDLte *int `form:"id_lte"`

	ChineseName        *string `form:"chinese_name"`
	ChineseNameInclude *string `form:"chinese_name_include"`
}
