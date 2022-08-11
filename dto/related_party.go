package dto

//RelatedPartyDTO dto只接收、发送数据，并不会对数据库进行任何操作
//所有操作数据库的工作会经过dao层的清洗后都交给model来完成
//这里必须是指针类型，因为只有指针才能向前端传递nil
//以后如果接收和推送用不同的dto，可以考虑不用指针
type RelatedPartyDTO struct {
	ID                      int     `form:"id" json:"id"`
	ChineseName             *string `json:"chinese_name"`
	EnglishName             *string `json:"english_name"`
	Address                 *string `json:"address"`
	UniformSocialCreditCode *string `json:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string `json:"telephone"`
}

// RelatedPartyListDTO 是list查询的过滤器
// 在dto传递给sqlCondition时，空值会被忽略
type RelatedPartyListDTO struct {
	ID    int  `json:"id"`
	IDGte *int `json:"id_gte"`
	IDLte *int `json:"id_lte"`

	ChineseName        *string `json:"chinese_name"`
	ChineseNameInclude *string `json:"chinese_name_include"`

	ListDTO
}
