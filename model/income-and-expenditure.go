package model

import "time"

type IncomeAndExpenditure struct {
	BasicModel
	//连接其他表的id
	ProjectId  *int64 //项目id
	ContractId *int64 //合同id
	//连接dictionary_item表的id
	FundDirection *int64  //资金方向(收款、付款)
	Currency      *int64  //币种
	Kind          *int64  //款项种类(计划、实际、预测等)
	Type          *string //款项类型(预付款、进度款、尾款等)
	DataSource    *int64  //数据来源，来自哪个视图(来自收款、收汇、收票等某个视图)

	Term *string //条款、方式

	//日期
	Date *time.Time `gorm:"type:date"`
	//数字
	Amount       *float64 //金额
	ExchangeRate *float64 //汇率
	//字符串
	Remarks    *string //备注
	Attachment *string //附件

	ImportedApprovalId *string //外部导入的付款审批id

}

func (i *IncomeAndExpenditure) TableName() string {
	return "income_and_expenditure"
}
