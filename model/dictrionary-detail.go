package model

import (
	"pmis-backend-go/global"
)

type DictionaryDetail struct {
	BasicModel
	DictionaryTypeID   int     //字典类型的id
	Name               string  //名称
	Sequence           *int    //用于排序的值
	IsValidForFrontend *bool   //是否在前端展现
	Remarks            *string //备注
}

// TableName 修改数据库的表名
func (*DictionaryDetail) TableName() string {
	return "dictionary_detail"
}

type dictionaryDetailFormat struct {
	TypeName    string
	DetailNames []string
}

var initialDictionary = []dictionaryDetailFormat{
	{
		TypeName:    "省份",
		DetailNames: []string{"上海", "北京", "山东", "河南"},
	},
	{
		TypeName:    "收付款方式",
		DetailNames: []string{"现金", "汇票", "信用证"},
	},
	{
		TypeName:    "进度类型",
		DetailNames: []string{"计划进度", "实际进度", "预测进度"},
	},
	{
		TypeName:    "币种",
		DetailNames: []string{"人民币", "美元", "欧元", "港币"},
	},
	{
		TypeName:    "合同类型",
		DetailNames: []string{"采购", "销售", "代理"},
	},
	{
		TypeName:    "项目类型",
		DetailNames: []string{"EPC", "EP"},
	},
	{
		TypeName:    "进度的数据来源",
		DetailNames: []string{"系统计算", "人工填写"},
	},
	{
		TypeName:    "数据范围",
		DetailNames: []string{"用户所在部门", "用户所在部门和子部门", "所有部门", "自定义部门"},
	},
}

func generateDictionaryDetail() (err error) {
	var dictionaryDetails []DictionaryDetail
	for i := range initialDictionary {
		var dictionaryTypeRecord DictionaryType
		global.DB.FirstOrCreate(&dictionaryTypeRecord, DictionaryType{Name: initialDictionary[i].TypeName})
		for j := range initialDictionary[i].DetailNames {
			dictionaryDetails = append(dictionaryDetails, DictionaryDetail{
				DictionaryTypeID: dictionaryTypeRecord.ID,
				Name:             initialDictionary[i].DetailNames[j],
			})
		}
	}

	for _, dictionaryDetail := range dictionaryDetails {
		err = global.DB.FirstOrCreate(&DictionaryDetail{}, dictionaryDetail).Error
		if err != nil {
			return err
		}
	}
	return nil
}
