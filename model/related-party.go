package model

type RelatedParty struct {
	BaseModel
	ChineseName             *string `json:"chinese_name" binding:"required"`               //中文名称
	EnglishName             *string `json:"english_name" binding:"required"`               //英文名称
	Address                 *string `json:"address" binding:"required"`                    //地址
	UniformSocialCreditCode *string `json:"uniform_social_credit_code" binding:"required"` //统一社会信用代码
	Telephone               *string `json:"telephone" binding:"required"`                  //电话
	//外键
}

// TableName 修改数据库的表名
func (*RelatedParty) TableName() string {
	return "related_party"
}
