package model

type PredictedIncomeAndExpenditure struct {
	BasicModel
	ProjectID           *int     //项目id
	ContractID          *int     //合同id
	FundDirection       *string  //资金方向，收款还是付款
	NameOfTheOtherParty *string  //对方名称
	Type                *string  //款项类型
	Condition           *string  //条件
	Date                *string  `gorm:"type:date"` //日期
	Amount              *float64 //金额
	Currency            *string  //币种
	ExchangeRate        *float64 //汇率
	Term                *string  //条款、方式
	Remark              *string  //备注
	Attachment          *string  //附件
}

func (*PredictedIncomeAndExpenditure) TableName() string {
	return "predicted_income_and_expenditure"
}
