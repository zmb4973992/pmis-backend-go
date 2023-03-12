package model

type PredictedReceiptAndPayment struct {
	BasicModel
	ProjectID                *int     //项目id
	ContractID               *int     //合同id
	FundDirection            *string  //资金方向，收款还是付款
	NameOfTheOtherParty      *string  //对方名称
	TypeOfReceiptAndPayment  *string  //款项类型
	Condition                *string  //条件
	Date                     *string  `gorm:"type:date"` //日期
	Amount                   *float64 //金额
	Currency                 *string  //币种
	ExchangeRate             *float64 //汇率
	TermsOfReceiptAndPayment *string  //收付款方式
	Remark                   *string  //备注
	Attachment               *string  //附件
}

func (*PredictedReceiptAndPayment) TableName() string {
	return "predicted_receipt_and_payment"
}
