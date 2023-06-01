package model

import "time"

type CumulativeIncomeAndExpenditure struct {
	BasicModel
	//连接其他表的id
	ProjectID int64 //项目ID
	//ProjectName string //项目名称
	//连接dictionary_item表的id
	//日期
	Date *time.Time `gorm:"type:date"`
	//数字
	TotalPlannedExpenditure       float64 //计划付款总额
	TotalActualExpenditure        float64 //实际付款总额
	TotalForecastedExpenditure    float64 //预测付款总额
	PlannedExpenditureProgress    float64 //计划付款进度
	ActualExpenditureProgress     float64 //实际付款进度
	ForecastedExpenditureProgress float64 //预测付款进度
	TotalPlannedIncome            float64 //计划收款总额
	TotalActualIncome             float64 //实际收款总额
	TotalForecastedIncome         float64 //预测收款总额
	PlannedIncomeProgress         float64 //计划收款进度
	ActualIncomeProgress          float64 //实际收款进度
	ForecastedIncomeProgress      float64 //预测收款进度
}

func (*CumulativeIncomeAndExpenditure) TableName() string {
	return "cumulative_income_and_expenditure"
}
