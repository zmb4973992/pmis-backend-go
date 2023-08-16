package model

type RelatedParty struct {
	BasicModel
	Name                    *string //名称
	EnglishName             *string //英文名称
	Address                 *string //地址
	UniformSocialCreditCode *string //统一社会信用代码
	Telephone               *string //电话
	Remarks                 *string //备注
	FileIDs                 *string //附件
	ImportedOriginalName    *string //导入的原名称
}

// TableName 修改数据库的表名
func (r *RelatedParty) TableName() string {
	return "related_party"
}
