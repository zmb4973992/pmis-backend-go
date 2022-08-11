package model

import "time"

type Contract struct {
	BaseModel
	ID                       int
	ProjectID                *int       //项目id
	Name                     *string    //合同名称
	Code                     *string    //合同编码
	RelatedPartyID           *int       //相关方id
	FundDirection            *string    //资金方向
	OurSignatory             *string    //我方签约名称
	Amount                   *float64   //金额
	Currency                 *string    //币种
	ExchangeRate             *float64   //汇率
	SigningDate              *time.Time //签约日期
	EffectiveDate            *time.Time //生效日期
	CommissioningDate        *time.Time //调试日期
	AgreedConstructionPeriod *int       //约定工期，天
	CompletionDate           *time.Time //完工日期
	JobDescription           *string    //工作内容
	Deliverables             *string    //交付物
	Penalty                  *string    //罚则
	Attachment               *string    //附件
	Operator                 *string    //经办人
	DepartmentID             *int       //所属部门

	ActualReceiptAndPayments    []ActualReceiptAndPayment    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PlannedReceiptAndPayments   []PlannedReceiptAndPayment   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PredictedReceiptAndPayments []PredictedReceiptAndPayment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 修改表名
func (Contract) TableName() string {
	return "contract"
}
