package model

import (
	"pmis-backend-go/global"
)

type DictionaryItem struct {
	BasicModel
	DictionaryTypeID   int     //字典类型的id
	Name               string  //名称
	Sequence           *int    //用于排序的值
	IsValidForFrontend *bool   //是否在前端展现
	Remarks            *string //备注
}

// TableName 修改数据库的表名
func (*DictionaryItem) TableName() string {
	return "dictionary_item"
}

type dictionaryItemFormat struct {
	TypeName  string
	ItemNames []string
}

var initialDictionary = []dictionaryItemFormat{
	{
		TypeName:  "省份",
		ItemNames: []string{"上海", "北京", "山东", "河南"},
	},
	{
		TypeName:  "收付款方式",
		ItemNames: []string{"现金", "汇票", "信用证"},
	},
	{
		TypeName:  "进度类型",
		ItemNames: []string{"计划进度", "实际进度", "预测进度"},
	},
	{
		TypeName:  "币种",
		ItemNames: []string{"人民币", "美元", "欧元", "港币"},
	},
	{
		TypeName:  "合同类型",
		ItemNames: []string{"采购", "销售", "代理"},
	},
	{
		TypeName:  "项目类型",
		ItemNames: []string{"EPC", "EP"},
	},
	{
		TypeName:  "进度的数据来源",
		ItemNames: []string{"系统计算", "人工填写"},
	},
}

func generateDictionary() (err error) {
	var dictionaryItems []DictionaryItem
	for i := range initialDictionary {
		var dictionaryTypeRecord DictionaryType
		global.DB.FirstOrCreate(&dictionaryTypeRecord, DictionaryType{Name: initialDictionary[i].TypeName})
		for j := range initialDictionary[i].ItemNames {
			dictionaryItems = append(dictionaryItems, DictionaryItem{
				DictionaryTypeID: dictionaryTypeRecord.ID,
				Name:             initialDictionary[i].ItemNames[j],
			})
		}
	}

	for _, dictionaryItem := range dictionaryItems {
		err = global.DB.FirstOrCreate(&DictionaryItem{}, dictionaryItem).Error
		if err != nil {
			return err
		}
	}
	return nil
}
