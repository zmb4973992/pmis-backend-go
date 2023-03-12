package model

import "time"

type ActualReceiptAndPayment struct {
	BasicModel
	ProjectID                *int       `gorm:"comment:'项目ID'"` //项目id
	ContractID               *int       //合同id
	FundDirection            *string    //资金方向，收款还是付款
	NameOfTheOtherParty      *string    //对方名称
	TypeOfReceiptAndPayment  *string    //款项类型
	Date                     *time.Time `gorm:"type:date;"` //日期
	Amount                   *float64   //金额
	Currency                 *string    //币种
	ExchangeRate             *float64   //汇率
	TermsOfReceiptAndPayment *string    //收付款方式
	Remark                   *string    //备注
	Attachment               *string    //附件
}

func (*ActualReceiptAndPayment) TableName() string {
	return "actual_receipt_and_payment"
}
