package model

type RelatedParty struct {
	BasicModel
	Name                    *string `json:"name"`                       //名称
	EnglishName             *string `json:"english_name"`               //英文名称
	Address                 *string `json:"address"`                    //地址
	UniformSocialCreditCode *string `json:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string `json:"telephone"`                  //电话
	Remarks                 *string `json:"remarks"`                    //备注
	ImportedOriginalName    *string `json:"imported_original_name"`     //导入的原名称
}

// TableName 修改数据库的表名
func (*RelatedParty) TableName() string {
	return "related_party"
}
