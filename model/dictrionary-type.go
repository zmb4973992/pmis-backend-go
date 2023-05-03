package model

import (
	"gorm.io/gorm"
	"pmis-backend-go/global"
)

type DictionaryType struct {
	BasicModel
	SnowID  uint64  `gorm:"not null;uniqueIndex;"`
	Name    string  //名称
	Sort    *int    //排序
	Status  *bool   //状态
	Remarks *string //备注
}

// TableName 修改数据库的表名
func (*DictionaryType) TableName() string {
	return "dictionary_type"
}

func (d *DictionaryType) BeforeDelete(tx *gorm.DB) error {
	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []DictionaryDetail
	err = tx.Where("dictionary_type_id = ?", d.ID).
		Find(&records).Delete(&records).Error
	if err != nil {
		return err
	}
	return nil
}

type dictionaryTypeFormat struct {
	SnowID uint64
	Name   string
	Sort   int
}

var dictionaryTypes = []dictionaryTypeFormat{
	{
		SnowID: 458788161124827656,
		Name:   "国家",
		Sort:   1,
	},
	{
		SnowID: 458789350495224328,
		Name:   "项目类型",
		Sort:   2,
	},
	{
		SnowID: 458789373362569736,
		Name:   "项目状态",
		Sort:   3,
	},
	{
		SnowID: 458789387858084360,
		Name:   "合同类型",
		Sort:   4,
	},
	{
		SnowID: 458789410926756360,
		Name:   "敏感词",
		Sort:   5,
	},
	{
		SnowID: 458789428911931912,
		Name:   "敏感词",
		Sort:   6,
	},
	{
		SnowID: 458789444648960520,
		Name:   "收付款方式",
		Sort:   7,
	},
	{
		SnowID: 458789457382867464,
		Name:   "币种",
		Sort:   8,
	},
	{
		SnowID: 458789470754308616,
		Name:   "进度类型",
		Sort:   9,
	},
	{
		SnowID: 458789483035230728,
		Name:   "银行名称",
		Sort:   10,
	},
	{
		SnowID: 458789495819469320,
		Name:   "进度的数据来源",
		Sort:   11,
	},
	{
		SnowID: 458789515432034824,
		Name:   "合同资金方向",
		Sort:   12,
	},
	{
		SnowID: 458789526672769544,
		Name:   "我方签约主体",
		Sort:   13,
	},
	{
		SnowID: 458789526672769444,
		Name:   "省份",
		Sort:   14,
	},
	{
		SnowID: 458822687075074568,
		Name:   "数据范围",
		Sort:   15,
	},
}

func generateDictionaryType() (err error) {
	for _, dictionaryType := range dictionaryTypes {
		err = global.DB.FirstOrCreate(&DictionaryType{}, dictionaryType).Error
		if err != nil {
			return err
		}
	}
	return nil
}
