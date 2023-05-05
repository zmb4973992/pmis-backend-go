package model

import (
	"github.com/yitter/idgenerator-go/idgen"
	"pmis-backend-go/global"
)

type DictionaryDetail struct {
	BasicModel
	DictionaryTypeSnowID int64   //字典类型的SnowID
	Name                 string  //名称
	Sequence             *int    //用于排序的值
	Status               *bool   //是否启用
	Remarks              *string //备注
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
		//先找到字典类型的记录
		var dictionaryTypeInfo DictionaryType
		err = global.DB.Where("name = ?", initialDictionary[i].TypeName).
			First(&dictionaryTypeInfo).Error
		if err != nil {
			return err
		}

		for j := range initialDictionary[i].DetailNames {
			dictionaryDetails = append(dictionaryDetails, DictionaryDetail{
				DictionaryTypeSnowID: dictionaryTypeInfo.SnowID,
				Name:                 initialDictionary[i].DetailNames[j],
			})
		}
	}

	for _, dictionaryDetail := range dictionaryDetails {
		err = global.DB.Where("name = ?", dictionaryDetail.Name).
			Attrs(&DictionaryDetail{
				BasicModel: BasicModel{
					SnowID: idgen.NextId(),
				},
				Status: BoolToPointer(true),
			}).
			FirstOrCreate(&dictionaryDetail).Error
		if err != nil {
			return err
		}
	}
	return nil
}
