package model

type Project struct {
	BaseModel
	ProjectCode      *string
	ProjectFullName  *string
	ProjectShortName *string
	Country          *string
	Province         *string
	ProjectType      *string
	Department       *string
	Amount           *float64
	Currency         *string
	ExchangeRate     *float64
	RelatedPartyID   *int
	//外键
	Disassemblies               []Disassembly                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActualReceiptAndPayments    []ActualReceiptAndPayment    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PlannedReceiptAndPayments   []PlannedReceiptAndPayment   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PredictedReceiptAndPayments []PredictedReceiptAndPayment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProjectAndUsers             []ProjectAndUser             `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OperationRecords            []OperationRecord            `gorm:"constraint:OnUpdate:CASCADE;"`
}

// TableName 将表名改为project
func (Project) TableName() string {
	return "project"
}
