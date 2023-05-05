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
	if d.SnowID == 0 {
		return nil
	}

	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []DictionaryDetail
	err = tx.Where(DictionaryDetail{DictionaryTypeSnowID: d.SnowID}).
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

var dictionaryTypes = []DictionaryType{
	{
		Name: "国家",
		Sort: IntToPointer(1),
	},
	{
		Name: "项目类型",
		Sort: IntToPointer(2),
	},
	{
		Name: "项目状态",
		Sort: IntToPointer(3),
	},
	{
		Name: "合同类型",
		Sort: IntToPointer(4),
	},
	{
		Name: "敏感词",
		Sort: IntToPointer(5),
	},
	{
		Name: "敏感词",
		Sort: IntToPointer(6),
	},
	{
		Name: "收付款方式",
		Sort: IntToPointer(7),
	},
	{
		Name: "币种",
		Sort: IntToPointer(8),
	},
	{
		Name: "进度类型",
		Sort: IntToPointer(9),
	},
	{
		Name: "银行名称",
		Sort: IntToPointer(10),
	},
	{
		Name: "进度的数据来源",
		Sort: IntToPointer(11),
	},
	{
		Name: "合同资金方向",
		Sort: IntToPointer(12),
	},
	{
		Name: "我方签约主体",
		Sort: IntToPointer(13),
	},
	{
		Name: "省份",
		Sort: IntToPointer(14),
	},
	{
		Name: "数据范围",
		Sort: IntToPointer(15),
	},
	{
		Name: "LDAP允许的OU",
		Sort: IntToPointer(16),
	},
}

func generateDictionaryType() (err error) {
	for i := range dictionaryTypes {
		err = global.DB.Where("name = ?", dictionaryTypes[i].Name).
			Where("sort = ?", dictionaryTypes[i].Sort).
			Attrs(DictionaryType{
				BasicModel: BasicModel{
					SnowID: idgen.NextId(),
				},
				Status: BoolToPointer(true),
			}).
			FirstOrCreate(&dictionaryTypes[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}
