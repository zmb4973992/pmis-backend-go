package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ContractDailyAndCumulativeIncomeUpdate struct {
	UserID int64
	//连接关联表的id
	ContractID int64 `json:"contract_id,omitempty"`
	//连接dictionary_item表的id

	//日期
	//数字

	//字符串
}

type ContractDailyAndCumulativeIncomeGetList struct {
	list.Input
	ContractID int64  `json:"contract_id,omitempty"`
	DateGte    string `json:"date_gte,omitempty"`
	DateLte    string `json:"date_lte,omitempty"`
}

//以下为出参

type ContractDailyAndCumulativeIncomeOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示

	//日期
	Date *string `json:"date"`

	//关联表的详情，不需要gorm查询，需要在json中显示
	//dictionary_item表的详情，不需要gorm查询，需要在json中显示

	DailyActualIncome        *float64 `json:"daily_actual_income"`        //当日实际收款金额
	TotalPlannedIncome       *float64 `json:"total_planned_income"`       //计划收款总额
	TotalActualIncome        *float64 `json:"total_actual_income"`        //实际收款总额
	TotalForecastedIncome    *float64 `json:"total_forecasted_income"`    //预测收款总额
	PlannedIncomeProgress    *float64 `json:"planned_income_progress"`    //计划收款进度
	ActualIncomeProgress     *float64 `json:"actual_income_progress"`     //实际收款进度
	ForecastedIncomeProgress *float64 `json:"forecasted_income_progress"` //预测收款进度

	//其他属性

}

func (c *ContractDailyAndCumulativeIncomeUpdate) Update() (errCode int) {
	//连接关联表的id
	{
		if c.ContractID > 0 {
			var typeOfIncomeAndExpenditure int64
			err := global.DB.Model(&model.DictionaryType{}).
				Where("name = '收付款的种类'").
				Select("id").
				First(&typeOfIncomeAndExpenditure).Error
			if err != nil {
				return util.ErrorFailToUpdateRecord
			}

			var planned int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_id = ?", typeOfIncomeAndExpenditure).
				Where("name = '计划'").
				Select("id").
				First(&planned).Error
			if err != nil {
				return util.ErrorFailToUpdateRecord
			}

			var actual int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_id = ?", typeOfIncomeAndExpenditure).
				Where("name = '实际'").
				Select("id").
				First(&actual).Error
			if err != nil {
				return util.ErrorFailToUpdateRecord
			}

			var forecasted int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_id = ?", typeOfIncomeAndExpenditure).
				Where("name = '预测'").
				Select("id").
				First(&forecasted).Error
			if err != nil {
				return util.ErrorFailToUpdateRecord
			}

			var contract model.Contract
			err = global.DB.Where("id = ?", c.ContractID).
				First(&contract).Error
			if err != nil {
				return util.ErrorFailToUpdateRecord
			}

			global.DB.Where("contract_id = ?", c.ContractID).
				Delete(&model.ContractDailyAndCumulativeIncome{})

			var fundDirectionOfIncomeAndExpenditure int64
			err = global.DB.Model(&model.DictionaryType{}).
				Where("name = '收付款的资金方向'").
				Select("id").
				First(&fundDirectionOfIncomeAndExpenditure).Error
			if err != nil {
				return util.ErrorFailToUpdateRecord
			}

			var income int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_id = ?", fundDirectionOfIncomeAndExpenditure).
				Where("name = '收款'").
				Select("id").
				First(&income).Error
			if err != nil {
				return util.ErrorFailToUpdateRecord
			}

			var dates []time.Time
			global.DB.Model(&model.IncomeAndExpenditure{}).
				Where("contract_id = ?", c.ContractID).
				Where("fund_direction = ?", income).
				Distinct("date").
				Order("date desc").
				Find(&dates)
			//fmt.Println("日期：", dates)

			records := make(chan model.ContractDailyAndCumulativeIncome, 10)

			for i := range dates {
				j := i
				go func() {
					var record model.ContractDailyAndCumulativeIncome

					var totalPlannedIncome float64
					var countForPlanned int64
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("contract_id = ?", c.ContractID).
						Where("kind = ?", planned).
						Where("fund_direction = ?", income).
						Where("date = ?", dates[j]).
						Count(&countForPlanned)
					if countForPlanned > 0 {
						global.DB.Model(&model.IncomeAndExpenditure{}).
							Where("contract_id = ?", c.ContractID).
							Where("kind = ?", planned).
							Where("fund_direction = ?", income).
							Where("date <= ?", dates[j]).
							Select("coalesce(sum(amount * exchange_rate),0)").
							Find(&totalPlannedIncome)
						record.TotalPlannedIncome = &totalPlannedIncome
						//fmt.Println("计划收款总额：", totalPlannedIncome)
						if contract.Amount != nil && *contract.Amount > 0 {
							var plannedIncomeProgress = totalPlannedIncome / *contract.Amount
							record.PlannedIncomeProgress = &plannedIncomeProgress
							//fmt.Println("计划收款进度：", plannedIncomeProgress)
						}
					}

					var totalActualIncome float64
					var countForActual int64
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("contract_id = ?", c.ContractID).
						Where("kind = ?", actual).
						Where("fund_direction = ?", income).
						Where("date = ?", dates[j]).
						Count(&countForActual)
					if countForActual > 0 {
						global.DB.Model(&model.IncomeAndExpenditure{}).
							Where("contract_id = ?", c.ContractID).
							Where("kind = ?", actual).
							Where("fund_direction = ?", income).
							Where("date <= ?", dates[j]).
							Select("coalesce(sum(amount * exchange_rate),0)").
							Find(&totalActualIncome)
						record.TotalActualIncome = &totalActualIncome
						//fmt.Println("实际收款总额：", totalActualIncome)
						if contract.Amount != nil && *contract.Amount > 0 {
							var actualIncomeProgress = totalActualIncome / *contract.Amount
							record.ActualIncomeProgress = &actualIncomeProgress
							//fmt.Println("实际收款进度：", actualIncomeProgress)
						}
					}

					var totalForecastedIncome float64
					var countForForecasted int64
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("contract_id = ?", c.ContractID).
						Where("kind = ?", forecasted).
						Where("fund_direction = ?", income).
						Where("date = ?", dates[j]).
						Count(&countForForecasted)
					if countForForecasted > 0 {
						global.DB.Model(&model.IncomeAndExpenditure{}).
							Where("contract_id = ?", c.ContractID).
							Where("kind = ?", forecasted).
							Where("fund_direction = ?", income).
							Where("date <= ?", dates[j]).
							Select("coalesce(sum(amount * exchange_rate),0)").
							Find(&totalForecastedIncome)
						record.TotalForecastedIncome = &totalForecastedIncome
						//fmt.Println("预测收款总额：", totalForecastedIncome)
						if contract.Amount != nil && *contract.Amount > 0 {
							var forecastedIncomeProgress = totalForecastedIncome / *contract.Amount
							record.ForecastedIncomeProgress = &forecastedIncomeProgress
							//fmt.Println("预测收款进度：", forecastedIncomeProgress)
						}
					}

					var dailyActualIncome float64
					var countForDailyActual int64
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("contract_id = ?", c.ContractID).
						Where("kind = ?", actual).
						Where("fund_direction = ?", income).
						Where("date = ?", dates[j]).
						Count(&countForDailyActual)
					if countForDailyActual > 0 {
						global.DB.Model(&model.IncomeAndExpenditure{}).
							Where("contract_id = ?", c.ContractID).
							Where("kind = ?", actual).
							Where("fund_direction = ?", income).
							Where("date <= ?", dates[j]).
							Select("coalesce(sum(amount * exchange_rate),0)").
							Find(&dailyActualIncome)
						record.DailyActualIncome = &dailyActualIncome
						//fmt.Println("当日实际收款金额：", dailyActualIncome)
					}

					record.Creator = &c.UserID
					record.ContractID = c.ContractID
					record.Date = &dates[j]

					records <- record
					//fmt.Println("------------------------")
				}()
			}

			go func() {
				for {
					select {
					case record := <-records:
						global.DB.Create(&record)
					}
				}
			}()
		}
	}

	return util.Success
}

func (c *ContractDailyAndCumulativeIncomeGetList) GetList() (
	outputs []ContractDailyAndCumulativeIncomeOutput,
	errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.ContractDailyAndCumulativeIncome{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where
	if c.ContractID > 0 {
		db = db.Where("contract_id = ?", c.ContractID)
	}

	if c.DateGte != "" {
		db = db.Where("date >= ?", c.DateGte)
	}

	if c.DateLte != "" {
		db = db.Where("date <= ?", c.DateLte)
	}

	//count
	var count int64
	db.Count(&count)

	//order
	orderBy := c.SortingInput.OrderBy
	desc := c.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("date desc")
		} else {
			db = db.Order("date")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.ContractDailyAndCumulativeIncome{}, orderBy)
		if !exists {
			return nil, util.ErrorSortingFieldDoesNotExist, nil
		}
		//如果要求降序排列
		if desc == true {
			db = db.Order(orderBy + " desc")
		} else { //如果没有要求排序方式
			db = db.Order(orderBy)
		}
	}

	//limit
	page := 1
	if c.PagingInput.Page > 0 {
		page = c.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if c.PagingInput.PageSize != nil && *c.PagingInput.PageSize >= 0 &&
		*c.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *c.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.ContractDailyAndCumulativeIncome{}).
		Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	for i := range outputs {
		//处理float64精度问题
		if outputs[i].TotalPlannedIncome != nil {
			temp := util.Round(*outputs[i].TotalPlannedIncome, 2)
			outputs[i].TotalPlannedIncome = &temp
		}
		if outputs[i].TotalActualIncome != nil {
			temp := util.Round(*outputs[i].TotalActualIncome, 2)
			outputs[i].TotalActualIncome = &temp
		}
		if outputs[i].TotalForecastedIncome != nil {
			temp := util.Round(*outputs[i].TotalForecastedIncome, 2)
			outputs[i].TotalForecastedIncome = &temp
		}
		if outputs[i].PlannedIncomeProgress != nil {
			temp := util.Round(*outputs[i].PlannedIncomeProgress, 3)
			outputs[i].PlannedIncomeProgress = &temp
		}
		if outputs[i].ActualIncomeProgress != nil {
			temp := util.Round(*outputs[i].ActualIncomeProgress, 3)
			outputs[i].ActualIncomeProgress = &temp
		}
		if outputs[i].ForecastedIncomeProgress != nil {
			temp := util.Round(*outputs[i].ForecastedIncomeProgress, 3)
			outputs[i].ForecastedIncomeProgress = &temp
		}

		//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
		//需要取年月日(即前9位)
		{
			if outputs[i].Date != nil {
				temp := *outputs[i].Date
				*outputs[i].Date = temp[:10]
			}
		}
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return outputs,
		util.Success,
		&list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		}
}
