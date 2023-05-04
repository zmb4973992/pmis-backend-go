package model

import (
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"pmis-backend-go/global"
)

type DictionaryType struct {
	BasicModel
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
	err = tx.Where("dictionary_type_snow_id = ?", d.SnowID).
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
		Name: "国家",
		Sort: 1,
	},
	{
		Name: "项目类型",
		Sort: 2,
	},
	{
		Name: "项目状态",
		Sort: 3,
	},
	{
		Name: "合同类型",
		Sort: 4,
	},
	{
		Name: "敏感词",
		Sort: 5,
	},
	{
		Name: "敏感词",
		Sort: 6,
	},
	{
		Name: "收付款方式",
		Sort: 7,
	},
	{
		Name: "币种",
		Sort: 8,
	},
	{
		Name: "进度类型",
		Sort: 9,
	},
	{
		Name: "银行名称",
		Sort: 10,
	},
	{
		Name: "进度的数据来源",
		Sort: 11,
	},
	{
		Name: "合同资金方向",
		Sort: 12,
	},
	{
		Name: "我方签约主体",
		Sort: 13,
	},
	{
		Name: "省份",
		Sort: 14,
	},
	{
		Name: "数据范围",
		Sort: 15,
	},
}

func generateDictionaryType() (err error) {
	for i := range dictionaryTypes {
		err = global.DB.Model(&DictionaryType{}).
			Where("name = ?", dictionaryTypes[i].Name).
			Where("sort = ?", dictionaryTypes[i].Sort).
			Attrs("snow_id = ?", idgen.NextId()).
			FirstOrCreate(&dictionaryTypes[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}
