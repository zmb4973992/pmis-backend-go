package model

import (
	"gorm.io/gorm"
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
