package model

import "time"

type Project struct {
	BaseModel
	ProjectCode          *string
	ProjectFullName      *string
	ProjectShortName     *string
	Country              *string
	Province             *string
	ProjectType          *string
	SegmentedProjectType *string //细分的项目类型
	Amount               *float64
	Currency             *string
	ExchangeRate         *float64
	ProjectStatus        *string
	OurSignatory         *string    //我方签约主体
	ConstructionPeriod   *int       //工期，天
	SigningDate          *time.Time `gorm:"type:date"` //签约日期
	EffectiveDate        *time.Time `gorm:"type:date"` //生效日期
	CommissioningDate    *time.Time `gorm:"type:date"` //调试日期
	DepartmentID         *int
	RelatedPartyID       *int
	Task                 *string //工作内容
	//外键
	Disassemblies               []Disassembly                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActualReceiptAndPayments    []ActualReceiptAndPayment    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PlannedReceiptAndPayments   []PlannedReceiptAndPayment   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PredictedReceiptAndPayments []PredictedReceiptAndPayment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProjectAndUsers             []ProjectAndUser             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 将表名改为project
func (Project) TableName() string {
	return "project"
}
