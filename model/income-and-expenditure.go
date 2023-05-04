package model

import "time"

type IncomeAndExpenditure struct {
	BasicModel
	//连接其他表的id
	ProjectSnowID  *uint64 //项目SnowID
	ContractSnowID *uint64 //合同SnowID
	//连接dictionary_item表的id
	FundDirection *int //资金方向(收款、付款)
	Currency      *int //币种
	Kind          *int //款项种类(计划、实际、预测等)
	//日期
	Date *time.Time `gorm:"type:date"`
	//数字
	Amount       *float64 //金额
	ExchangeRate *float64 //汇率
	//字符串
	Type       *string //款项类型(预付款、进度款、尾款等)
	Condition  *string //条件
	Term       *string //条款、方式
	Remarks    *string //备注
	Attachment *string //附件
}

func (*IncomeAndExpenditure) TableName() string {
	return "income_and_expenditure"
}
