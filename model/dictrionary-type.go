package model

import (
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
func (d *DictionaryType) TableName() string {
	return "dictionary_type"
}

func (d *DictionaryType) BeforeDelete(tx *gorm.DB) error {
	if d.Id == 0 {
		return nil
	}

	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []DictionaryDetail
	err = tx.Where(DictionaryDetail{DictionaryTypeId: d.Id}).
		Find(&records).Delete(&records).Error
	if err != nil {
		return err
	}

	return nil
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
		Name:    "合同的资金方向",
		Sort:    IntToPointer(12),
		Remarks: stringToPointer("如：收款合同、付款合同"),
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
	{
		Name:    "收付款的资金方向",
		Sort:    IntToPointer(17),
		Remarks: stringToPointer("如：收款、付款、不涉及收付款"),
	},
	{
		Name:    "收付款的种类",
		Sort:    IntToPointer(18),
		Remarks: stringToPointer("如：计划、实际、预测"),
	},
	{
		Name:    "款项类型", //预付款、发货款、尾款等
		Sort:    IntToPointer(19),
		Remarks: stringToPointer("如：预付款、进度款、尾款等"),
	},
	{
		Name: "tabFukuan视图中不要导入的记录",
		Sort: IntToPointer(20),
	},
	{Name: "操作类型",
		Sort:    IntToPointer(21),
		Remarks: stringToPointer("这是操作日志的操作类型，如：添加、修改、删除、查看、查看列表等"),
	},
	{Name: "收款的数据来源",
		Sort:    IntToPointer(22),
		Remarks: stringToPointer("标识收款记录来自率敏的哪个视图，如收款、收汇、收票等"),
	},
}

func generateDictionaryType() (err error) {
	for i := range dictionaryTypes {
		err = global.DB.Where("name = ?", dictionaryTypes[i].Name).
			Where("sort = ?", dictionaryTypes[i].Sort).
			Attrs(DictionaryType{
				Status: BoolToPointer(true),
			}).
			FirstOrCreate(&dictionaryTypes[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}
