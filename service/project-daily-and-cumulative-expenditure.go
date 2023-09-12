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

type ProjectDailyAndCumulativeExpenditureUpdate struct {
	UserId int64
	//连接关联表的id
	ProjectId int64 `json:"project_id,omitempty"`
	//连接dictionary_item表的id

	//日期
	//数字

	//字符串
}

type ProjectDailyAndCumulativeExpenditureGetList struct {
	list.Input
	ProjectId int64  `json:"project_id,omitempty"`
	DateGte   string `json:"date_gte,omitempty"`
	DateLte   string `json:"date_lte,omitempty"`
}

//以下为出参

type ProjectDailyAndCumulativeExpenditureOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	//ProjectId *int64 `json:"-"`
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示

	//日期
	Date *string `json:"date"`

	//关联表的详情，不需要gorm查询，需要在json中显示
	//ProjectExternal *ProjectOutput `json:"project" gorm:"-"`
	//dictionary_item表的详情，不需要gorm查询，需要在json中显示

	DailyActualExpenditure        *float64 `json:"daily_actual_expenditure"`        //当日实际付款金额
	TotalPlannedExpenditure       *float64 `json:"total_planned_expenditure"`       //计划付款总额
	TotalActualExpenditure        *float64 `json:"total_actual_expenditure"`        //实际付款总额
	TotalForecastedExpenditure    *float64 `json:"total_forecasted_expenditure"`    //预测付款总额
	PlannedExpenditureProgress    *float64 `json:"planned_expenditure_progress"`    //计划付款进度
	ActualExpenditureProgress     *float64 `json:"actual_expenditure_progress"`     //实际付款进度
	ForecastedExpenditureProgress *float64 `json:"forecasted_expenditure_progress"` //预测付款进度
	//其他属性

}

func (p *ProjectDailyAndCumulativeExpenditureUpdate) Update() (errCode int) {
	//连接关联表的id
	if p.ProjectId > 0 {
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
			Where("name = '计划'").Select("id").First(&planned).Error
		if err != nil {
			return util.ErrorFailToUpdateRecord
		}

		var actual int64
		err = global.DB.Model(&model.DictionaryDetail{}).
			Where("dictionary_type_id = ?", typeOfIncomeAndExpenditure).
			Where("name = '实际'").Select("id").First(&actual).Error
		if err != nil {
			return util.ErrorFailToUpdateRecord
		}

		var forecasted int64
		err = global.DB.Model(&model.DictionaryDetail{}).
			Where("dictionary_type_id = ?", typeOfIncomeAndExpenditure).
			Where("name = '预测'").Select("id").First(&forecasted).Error
		if err != nil {
			return util.ErrorFailToUpdateRecord
		}

		var fundDirectionOfContract int64
		err = global.DB.Model(&model.DictionaryType{}).
			Where("name = '合同的资金方向'").Select("id").First(&fundDirectionOfContract).Error
		if err != nil {
			return util.ErrorFailToUpdateRecord
		}

		var expenditureContract int64
		err = global.DB.Model(&model.DictionaryDetail{}).
			Where("dictionary_type_id = ?", fundDirectionOfContract).
			Where("name = '付款合同'").Select("id").First(&expenditureContract).Error
		if err != nil {
			return util.ErrorFailToUpdateRecord
		}

		//计算付款合同的总金额
		var totalAmountOfExpenditureContract float64
		err = global.DB.Model(&model.Contract{}).
			Where("project_id = ?", p.ProjectId).
			Where("fund_direction = ?", expenditureContract).
			Select("coalesce(sum(amount * exchange_rate),0)").
			Find(&totalAmountOfExpenditureContract).Error
		//fmt.Println("付款合同总金额：", totalAmountOfExpenditureContract)
		//fmt.Println("*********************************")

		global.DB.Where("project_id = ?", p.ProjectId).
			Delete(&model.ProjectDailyAndCumulativeExpenditure{})

		var fundDirectionOfIncomeAndExpenditure int64
		err = global.DB.Model(&model.DictionaryType{}).
			Where("name = '收付款的资金方向'").Select("id").First(&fundDirectionOfIncomeAndExpenditure).Error
		if err != nil {
			return util.ErrorFailToUpdateRecord
		}

		var expenditure int64
		err = global.DB.Model(&model.DictionaryDetail{}).
			Where("dictionary_type_id = ?", fundDirectionOfIncomeAndExpenditure).
			Where("name = '付款'").Select("id").First(&expenditure).Error
		if err != nil {
			return util.ErrorFailToUpdateRecord
		}

		var dates []time.Time
		global.DB.Model(&model.IncomeAndExpenditure{}).
			Where("project_id = ?", p.ProjectId).
			Where("fund_direction = ?", expenditure).
			Distinct("date").
			Order("date desc").
			Find(&dates)
		//fmt.Println("日期：", dates)

		records := make(chan model.ProjectDailyAndCumulativeExpenditure, 10)

		for i := range dates {
			j := i
			go func() {
				//fmt.Println("日期：", dates[j].Format("2006-01-02")[:10])
				var record model.ProjectDailyAndCumulativeExpenditure

				var totalPlannedExpenditure float64
				var countForPlanned int64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("project_id = ?", p.ProjectId).
					Where("kind = ?", planned).
					Where("fund_direction = ?", expenditure).
					Where("date = ?", dates[j]).
					Count(&countForPlanned)
				if countForPlanned > 0 {
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectId).
						Where("kind = ?", planned).
						Where("fund_direction = ?", expenditure).
						Where("date <= ?", dates[j]).
						Select("coalesce(sum(amount * exchange_rate),0)").
						Find(&totalPlannedExpenditure)
					record.TotalPlannedExpenditure = &totalPlannedExpenditure
					//fmt.Println("计划付款总额：", totalPlannedExpenditure)
					if totalAmountOfExpenditureContract > 0 {
						var plannedExpenditureProgress = totalPlannedExpenditure / totalAmountOfExpenditureContract
						record.PlannedExpenditureProgress = &plannedExpenditureProgress
						//fmt.Println("计划付款进度：", plannedExpenditureProgress)
					}
				}

				var totalActualExpenditure float64
				var countForActual int64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("project_id = ?", p.ProjectId).
					Where("kind = ?", actual).
					Where("fund_direction = ?", expenditure).
					Where("date = ?", dates[j]).
					Count(&countForActual)
				if countForActual > 0 {
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectId).
						Where("kind = ?", actual).
						Where("fund_direction = ?", expenditure).
						Where("date <= ?", dates[j]).
						Select("coalesce(sum(amount * exchange_rate),0)").
						Find(&totalActualExpenditure)
					record.TotalActualExpenditure = &totalActualExpenditure
					//fmt.Println("实际付款总额：", totalActualExpenditure)
					if totalAmountOfExpenditureContract > 0 {
						var actualExpenditureProgress = totalActualExpenditure / totalAmountOfExpenditureContract
						record.ActualExpenditureProgress = &actualExpenditureProgress
						//fmt.Println("实际付款进度：", actualExpenditureProgress)
					}
				}

				var totalForecastedExpenditure float64
				var countForForecasted int64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("project_id = ?", p.ProjectId).
					Where("kind = ?", forecasted).
					Where("fund_direction = ?", expenditure).
					Where("date = ?", dates[j]).
					Count(&countForForecasted)
				if countForForecasted > 0 {
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectId).
						Where("kind = ?", forecasted).
						Where("fund_direction = ?", expenditure).
						Where("date <= ?", dates[j]).
						Select("coalesce(sum(amount * exchange_rate),0)").
						Find(&totalForecastedExpenditure)
					record.TotalForecastedExpenditure = &totalForecastedExpenditure
					//fmt.Println("预测付款总额：", totalForecastedExpenditure)
					if totalAmountOfExpenditureContract > 0 {
						var forecastedExpenditureProgress = totalForecastedExpenditure / totalAmountOfExpenditureContract
						record.ForecastedExpenditureProgress = &forecastedExpenditureProgress
						//fmt.Println("预测付款进度：", forecastedExpenditureProgress)
					}
				}

				var dailyActualExpenditure float64
				var countForDailyActual int64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("project_id = ?", p.ProjectId).
					Where("kind = ?", actual).
					Where("fund_direction = ?", expenditure).
					Where("date = ?", dates[j]).
					Count(&countForDailyActual)
				if countForDailyActual > 0 {
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectId).
						Where("kind = ?", actual).
						Where("fund_direction = ?", expenditure).
						Where("date = ?", dates[j]).
						Select("coalesce(sum(amount * exchange_rate),0)").
						Find(&dailyActualExpenditure)
					record.DailyActualExpenditure = &dailyActualExpenditure
					//fmt.Println("当日付款金额：", dailyActualExpenditure)
				}

				record.Creator = &p.UserId
				record.ProjectId = p.ProjectId
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

	return util.Success
}

func (p *ProjectDailyAndCumulativeExpenditureGetList) GetList() (
	outputs []ProjectDailyAndCumulativeExpenditureOutput,
	errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.ProjectDailyAndCumulativeExpenditure{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where
	if p.ProjectId > 0 {
		db = db.Where("project_id = ?", p.ProjectId)
	}

	if p.DateGte != "" {
		db = db.Where("date >= ?", p.DateGte)
	}

	if p.DateLte != "" {
		db = db.Where("date <= ?", p.DateLte)
	}

	//count
	var count int64
	db.Count(&count)

	//order
	orderBy := p.SortingInput.OrderBy
	desc := p.SortingInput.Desc
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
		exists := util.FieldIsInModel(&model.ProjectDailyAndCumulativeExpenditure{}, orderBy)
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
	if p.PagingInput.Page > 0 {
		page = p.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if p.PagingInput.PageSize != nil && *p.PagingInput.PageSize >= 0 &&
		*p.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *p.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.ProjectDailyAndCumulativeExpenditure{}).
		Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	for i := range outputs {
		//处理float64精度问题
		if outputs[i].TotalPlannedExpenditure != nil {
			temp := util.Round(*outputs[i].TotalPlannedExpenditure, 2)
			outputs[i].TotalPlannedExpenditure = &temp
		}
		if outputs[i].TotalActualExpenditure != nil {
			temp := util.Round(*outputs[i].TotalActualExpenditure, 2)
			outputs[i].TotalActualExpenditure = &temp
		}
		if outputs[i].TotalForecastedExpenditure != nil {
			temp := util.Round(*outputs[i].TotalForecastedExpenditure, 2)
			outputs[i].TotalForecastedExpenditure = &temp
		}
		if outputs[i].PlannedExpenditureProgress != nil {
			temp := util.Round(*outputs[i].PlannedExpenditureProgress, 3)
			outputs[i].PlannedExpenditureProgress = &temp
		}
		if outputs[i].ActualExpenditureProgress != nil {
			temp := util.Round(*outputs[i].ActualExpenditureProgress, 3)
			outputs[i].ActualExpenditureProgress = &temp
		}
		if outputs[i].ForecastedExpenditureProgress != nil {
			temp := util.Round(*outputs[i].ForecastedExpenditureProgress, 3)
			outputs[i].ForecastedExpenditureProgress = &temp
		}
		if outputs[i].DailyActualExpenditure != nil {
			temp := util.Round(*outputs[i].DailyActualExpenditure, 2)
			outputs[i].DailyActualExpenditure = &temp
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
