package model

import "time"

type ProjectDailyAndCumulativeIncome struct {
	BasicModel
	//连接其他表的id
	ProjectID int64 //项目ID
	//连接dictionary_item表的id
	//日期
	Date *time.Time `gorm:"type:date"`
	//数字
	TotalPlannedIncome       *float64 //计划收款总额
	TotalActualIncome        *float64 //实际收款总额
	TotalForecastedIncome    *float64 //预测收款总额
	PlannedIncomeProgress    *float64 //计划收款进度
	ActualIncomeProgress     *float64 //实际收款进度
	ForecastedIncomeProgress *float64 //预测收款进度
	DailyActualIncome        *float64 //当日实际收款金额

}

func (p *ProjectDailyAndCumulativeIncome) TableName() string {
	return "project_daily_and_cumulative_income"
}
