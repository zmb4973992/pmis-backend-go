package model

type IncomeAndExpenditure struct {
	BasicModel
	ProjectID           *int     //项目id
	ContractID          *int     //合同id
	FundDirection       *int     //资金方向(收款、付款)，见dictionary_item
	NameOfTheOtherParty *string  //对方名称
	Kind                *int     //款项种类(计划、实际、预测等)，见dictionary_item
	Type                *string  //款项类型(预付款、进度款、尾款等)
	Condition           *string  //条件
	Date                *string  `gorm:"type:date"` //日期
	Amount              *float64 //金额
	Currency            *int     //币种，见dictionary_item
	ExchangeRate        *float64 //汇率
	Term                *string  //条款、方式
	Remark              *string  //备注
	Attachment          *string  //附件
}

func (*IncomeAndExpenditure) TableName() string {
	return "income_and_expenditure"
}
