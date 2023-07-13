package model

import (
	"gorm.io/gorm"
	"time"
)

type Project struct {
	BasicModel
	//连接其他表的id
	OrganizationID *int64 //见organization
	RelatedPartyID *int64 //见related_party
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

func (p *Project) AfterCreate(tx *gorm.DB) error {
	var disassembly Disassembly
	level := 1
	disassembly.Creator = p.Creator
	disassembly.LastModifier = p.LastModifier
	disassembly.ProjectID = &p.ID
	disassembly.Name = p.Name
	disassembly.Level = &level

	err = tx.Create(&disassembly).Error
	return err
}

func (p *Project) BeforeDelete(tx *gorm.DB) error {
	if p.ID == 0 {
		return nil
	}

	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var disassemblies []Disassembly
	err = tx.Where(&Disassembly{ProjectID: &p.ID}).
		Find(&disassemblies).Delete(&disassemblies).Error
	if err != nil {
		return err
	}

	var contracts []Contract
	err = tx.Where(&Contract{ProjectID: &p.ID}).
		Find(&contracts).Delete(&contracts).Error
	if err != nil {
		return err
	}

	var incomeAndExpenditures []IncomeAndExpenditure
	err = tx.Where(&IncomeAndExpenditure{ProjectID: &p.ID}).
		Find(&incomeAndExpenditures).Delete(&incomeAndExpenditures).Error
	if err != nil {
		return err
	}

	var projectCumulativeIncomes []ProjectCumulativeIncome
	err = tx.Where(&ProjectCumulativeIncome{ProjectID: p.ID}).
		Find(&projectCumulativeIncomes).Delete(&projectCumulativeIncomes).Error
	if err != nil {
		return err
	}

	var projectCumulativeExpenditures []ProjectCumulativeExpenditure
	err = tx.Where(&ProjectCumulativeExpenditure{ProjectID: p.ID}).
		Find(&projectCumulativeExpenditures).Delete(&projectCumulativeExpenditures).Error
	if err != nil {
		return err
	}

	return nil
}
