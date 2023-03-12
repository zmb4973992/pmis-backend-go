package model

import (
	"gorm.io/gorm"
	"pmis-backend-go/global"
)

type DictionaryType struct {
	BasicModel
	Name               string  //名称
	Sequence           *int    //排序
	IsValidForFrontend *bool   //是否在前端展现
	Remarks            *string //备注
}

// TableName 修改数据库的表名
func (*DictionaryType) TableName() string {
	return "dictionary_type"
}

func (d *DictionaryType) BeforeDelete(tx *gorm.DB) error {
	if d.ID > 0 {
		//如果有删除人的id，则记录下来
		//if d.Deleter != nil && *d.Deleter > 0 {
		//	err := tx.Model(&DictionaryType{}).Where("id = ?", d.ID).
		//		Update("deleter", d.Deleter).Error
		//	if err != nil {
		//		return err
		//	}
		//}
		//删除相关的子表记录
		//err = tx.Model(&DictionaryItem{}).Where("dictionary_type_id = ?", d.ID).
		//	Updates(map[string]any{
		//		"deleted_at": time.Now(),
		//		"deleter":    d.Deleter,
		//	}).Error
		//if err != nil {
		//	return err
		//}
	}
	return nil
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
		{Name: "LDAP允许的OU"},
		{Name: "进度类型"},
	}
	for _, dictionaryType := range dictionaryTypes {
		err := global.DB.FirstOrCreate(&DictionaryType{}, dictionaryType).Error
		if err != nil {
			return err
		}
	}
	return nil
}
