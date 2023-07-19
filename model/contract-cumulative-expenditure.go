package model

import "time"

type ContractCumulativeExpenditure struct {
	BasicModel
	//连接其他表的id
	ContractID int64 //合同ID
	//连接dictionary_item表的id
	//日期
	Date *time.Time `gorm:"type:date"`
	//数字
	TotalPlannedExpenditure       *float64 //计划付款总额
	TotalActualExpenditure        *float64 //实际付款总额
	TotalForecastedExpenditure    *float64 //预测付款总额
	PlannedExpenditureProgress    *float64 //计划付款进度
	ActualExpenditureProgress     *float64 //实际付款进度
	ForecastedExpenditureProgress *float64 //预测付款进度
}

func (*ContractCumulativeExpenditure) TableName() string {
	return "contract_cumulative_expenditure"
}
