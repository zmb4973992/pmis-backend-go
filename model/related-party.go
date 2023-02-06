package model

type RelatedParty struct {
	BaseModel
	ChineseName             *string `json:"chinese_name"`               //中文名称
	EnglishName             *string `json:"english_name"`               //英文名称
	Address                 *string `json:"address"`                    //地址
	UniformSocialCreditCode *string `json:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string `json:"telephone"`                  //电话
	Remarks                 *string `json:"remarks"`                    //备注
}

// TableName 修改数据库的表名
func (*RelatedParty) TableName() string {
	return "related_party"
}
