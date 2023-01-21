package model

import (
	"pmis-backend-go/global"
)

type DictionaryType struct {
	BaseModel
	Name            string           //名称
	Sort            *int             //排序
	Remarks         *string          //备注
	DictionaryItems []DictionaryItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 修改数据库的表名
func (DictionaryType) TableName() string {
	return "dictionary_type"
}

func generateDictionaryTypes() error {
	dictionaryTypes := []DictionaryType{
		{Name: "国家"},
		{Name: "省份"},
		{Name: "收付款方式"},
		{Name: "币种"},
		{Name: "合同类型"},
		{Name: "项目类型"},
		{Name: "项目状态"},
		{Name: "银行名称"},
		{Name: "合同资金方向"},
		{Name: "我方签约主体"},
		{Name: "敏感词"},
	}
	for _, dictionaryType := range dictionaryTypes {
		err := global.DB.FirstOrCreate(&DictionaryType{}, dictionaryType).Error
		if err != nil {
			return err
		}
	}
	return nil
}
