package model

import (
	"gorm.io/gorm"
	"time"
)

type Project struct {
	BasicModel
	//连接其他表的id
	OrganizationSnowID *int64 //见organization
	RelatedPartySnowID *int64 //见related_party
	//连接dictionary_item表的id
	Country      *int64
	Type         *int64
	DetailedType *int64 //细分的项目类型
	Currency     *int64
	Status       *int64
	OurSignatory *int64 //我方签约主体
	//日期
	SigningDate       *time.Time `gorm:"type:date"` //签约日期
	EffectiveDate     *time.Time `gorm:"type:date"` //生效日期
	CommissioningDate *time.Time `gorm:"type:date"` //调试日期
	//数字(允许为0、nil)
	Amount             *float64
	ExchangeRate       *float64
	ConstructionPeriod *int //工期，天
	//数字(不允许为0、nil，必须有值)，暂无

	//字符串(允许为null)
	Code    *string
	Name    *string
	Content *string //工作内容
	//字符串(不允许为nil，必须有值)，暂无

}

// TableName 将表名改为project
func (*Project) TableName() string {
	return "project"
}

func (d *Project) BeforeDelete(tx *gorm.DB) error {
	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []Disassembly
	err = tx.Where(&Disassembly{ProjectSnowID: &d.SnowID}).
		Find(&records).Delete(&records).Error
	if err != nil {
		return err
	}

	var records1 []Contract
	err = tx.Where(&Disassembly{ProjectSnowID: &d.SnowID}).
		Find(&records1).Delete(&records1).Error
	if err != nil {
		return err
	}

	var records2 []IncomeAndExpenditure
	err = tx.Where(&Disassembly{ProjectSnowID: &d.SnowID}).
		Find(&records2).Delete(&records2).Error
	if err != nil {
		return err
	}
	return nil
}
