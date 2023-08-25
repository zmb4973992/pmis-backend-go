package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ProjectDailyAndCumulativeIncomeUpdate struct {
	UserID int64
	//连接关联表的id
	ProjectID int64 `json:"project_id,omitempty"`
	//连接dictionary_item表的id

	//日期
	//数字

	//字符串
}

type ProjectDailyAndCumulativeIncomeGetList struct {
	list.Input
	ProjectID int64  `json:"project_id,omitempty"`
	DateGte   string `json:"date_gte,omitempty"`
	DateLte   string `json:"date_lte,omitempty"`
}

//以下为出参

type ProjectDailyAndCumulativeIncomeOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	//ProjectID *int64 `json:"-"`
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示

	//日期
	Date *string `json:"date"`

	//关联表的详情，不需要gorm查询，需要在json中显示
	//ProjectExternal *ProjectOutput `json:"project" gorm:"-"`
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

func (p *ProjectDailyAndCumulativeIncomeUpdate) Update() response.Common {
	//连接关联表的id
	{
		if p.ProjectID > 0 {
			var typeOfIncomeAndExpenditure int64
			err := global.DB.Model(&model.DictionaryType{}).
				Where("name = '收付款的种类'").Select("id").First(&typeOfIncomeAndExpenditure).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}

			var planned int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("name = '计划'").Select("id").First(&planned).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}

			var actual int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("name = '实际'").Select("id").First(&actual).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}

			var forecasted int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("name = '预测'").Select("id").First(&forecasted).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}

			var fundDirectionOfContract int64
			err = global.DB.Model(&model.DictionaryType{}).
				Where("name = '合同的资金方向'").Select("id").First(&fundDirectionOfContract).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}

			var incomeContract int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("name = '收款合同'").Select("id").First(&incomeContract).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}

			//计算收款合同的总金额
			var totalAmountOfIncomeContract float64
			err = global.DB.Model(&model.Contract{}).
				Where("project_id = ?", p.ProjectID).
				Where("fund_direction = ?", incomeContract).
				Select("coalesce(sum(amount * exchange_rate),0)").
				Find(&totalAmountOfIncomeContract).Error
			//fmt.Println("收款合同总金额：", totalAmountOfIncomeContract)
			//fmt.Println("*********************************")

			global.DB.Where("project_id = ?", p.ProjectID).
				Delete(&model.ProjectDailyAndCumulativeIncome{})

			var fundDirectionOfIncomeAndExpenditure int64
			err = global.DB.Model(&model.DictionaryType{}).
				Where("name = '收付款的资金方向'").Select("id").First(&fundDirectionOfIncomeAndExpenditure).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}

			var income int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_id = ?", fundDirectionOfIncomeAndExpenditure).
				Where("name = '收款'").Select("id").First(&income).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}

			var dates []time.Time
			global.DB.Model(&model.IncomeAndExpenditure{}).
				Where("project_id = ?", p.ProjectID).
				Where("fund_direction = ?", income).
				Distinct("date").
				Order("date desc").
				Find(&dates)
			//fmt.Println("日期：", dates)

			records := make(chan model.ProjectDailyAndCumulativeIncome, 10)

			for i := range dates {
				j := i
				go func() {
					//fmt.Println("日期：", dates[j].Format("2006-01-02")[:10])
					var record model.ProjectDailyAndCumulativeIncome

					var totalPlannedIncome float64
					var countForPlanned int64
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectID).
						Where("kind = ?", planned).
						Where("fund_direction = ?", income).
						Where("date = ?", dates[j]).
						Count(&countForPlanned)
					if countForPlanned > 0 {
						global.DB.Model(&model.IncomeAndExpenditure{}).
							Where("project_id = ?", p.ProjectID).
							Where("kind = ?", planned).
							Where("fund_direction = ?", income).
							Where("date <= ?", dates[j]).
							Select("coalesce(sum(amount * exchange_rate),0)").
							Find(&totalPlannedIncome)
						record.TotalPlannedIncome = &totalPlannedIncome
						//fmt.Println("计划收款总额：", totalPlannedIncome)
						if totalAmountOfIncomeContract > 0 {
							var plannedIncomeProgress = totalPlannedIncome / totalAmountOfIncomeContract
							record.PlannedIncomeProgress = &plannedIncomeProgress
							//fmt.Println("计划收款进度：", plannedIncomeProgress)
						}
					}

					var totalActualIncome float64
					var countForActual int64
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectID).
						Where("kind = ?", actual).
						Where("fund_direction = ?", income).
						Where("date = ?", dates[j]).
						Count(&countForActual)
					if countForActual > 0 {
						global.DB.Model(&model.IncomeAndExpenditure{}).
							Where("project_id = ?", p.ProjectID).
							Where("kind = ?", actual).
							Where("fund_direction = ?", income).
							Where("date <= ?", dates[j]).
							Select("coalesce(sum(amount * exchange_rate),0)").
							Find(&totalActualIncome)
						record.TotalActualIncome = &totalActualIncome
						//fmt.Println("实际收款总额：", totalActualIncome)
						if totalAmountOfIncomeContract > 0 {
							var actualIncomeProgress = totalActualIncome / totalAmountOfIncomeContract
							record.ActualIncomeProgress = &actualIncomeProgress
							//fmt.Println("实际收款进度：", actualIncomeProgress)
						}
					}

					var totalForecastedIncome float64
					var countForForecasted int64
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectID).
						Where("kind = ?", forecasted).
						Where("fund_direction = ?", income).
						Where("date = ?", dates[j]).
						Count(&countForForecasted)
					if countForForecasted > 0 {
						global.DB.Model(&model.IncomeAndExpenditure{}).
							Where("project_id = ?", p.ProjectID).
							Where("kind = ?", forecasted).
							Where("fund_direction = ?", income).
							Where("date <= ?", dates[j]).
							Select("coalesce(sum(amount * exchange_rate),0)").
							Find(&totalForecastedIncome)
						record.TotalForecastedIncome = &totalForecastedIncome
						//fmt.Println("预测收款总额：", totalForecastedIncome)
						if totalAmountOfIncomeContract > 0 {
							var forecastedIncomeProgress = totalForecastedIncome / totalAmountOfIncomeContract
							record.ForecastedIncomeProgress = &forecastedIncomeProgress
							//fmt.Println("预测收款进度：", forecastedIncomeProgress)
						}
					}

					var dailyActualIncome float64
					var countForDailyActual int64
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectID).
						Where("kind = ?", actual).
						Where("fund_direction = ?", income).
						Where("date = ?", dates[j]).
						Count(&countForDailyActual)
					if countForDailyActual > 0 {
						global.DB.Model(&model.IncomeAndExpenditure{}).
							Where("project_id = ?", p.ProjectID).
							Where("kind = ?", actual).
							Where("fund_direction = ?", income).
							Where("date = ?", dates[j]).
							Select("coalesce(sum(amount * exchange_rate),0)").
							Find(&dailyActualIncome)
						record.DailyActualIncome = &dailyActualIncome
						//fmt.Println("当日收款金额：", dailyActualIncome)
					}

					record.Creator = &p.UserID
					record.ProjectID = p.ProjectID
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

	return response.Success()
}

func (p *ProjectDailyAndCumulativeIncomeGetList) GetList() response.List {
	db := global.DB.Model(&model.ProjectDailyAndCumulativeIncome{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if p.ProjectID > 0 {
		db = db.Where("project_id = ?", p.ProjectID)
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
		exists := util.FieldIsInModel(&model.ProjectDailyAndCumulativeIncome{}, orderBy)
		if !exists {
			return response.FailureForList(util.ErrorSortingFieldDoesNotExist)
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

	//data
	var data []ProjectDailyAndCumulativeIncomeOutput
	db.Model(&model.ProjectDailyAndCumulativeIncome{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	for i := range data {
		//处理float64精度问题
		if data[i].TotalPlannedIncome != nil {
			temp := util.Round(*data[i].TotalPlannedIncome, 2)
			data[i].TotalPlannedIncome = &temp
		}
		if data[i].TotalActualIncome != nil {
			temp := util.Round(*data[i].TotalActualIncome, 2)
			data[i].TotalActualIncome = &temp
		}
		if data[i].TotalForecastedIncome != nil {
			temp := util.Round(*data[i].TotalForecastedIncome, 2)
			data[i].TotalForecastedIncome = &temp
		}
		if data[i].PlannedIncomeProgress != nil {
			temp := util.Round(*data[i].PlannedIncomeProgress, 3)
			data[i].PlannedIncomeProgress = &temp
		}
		if data[i].ActualIncomeProgress != nil {
			temp := util.Round(*data[i].ActualIncomeProgress, 3)
			data[i].ActualIncomeProgress = &temp
		}
		if data[i].ForecastedIncomeProgress != nil {
			temp := util.Round(*data[i].ForecastedIncomeProgress, 3)
			data[i].ForecastedIncomeProgress = &temp
		}
		if data[i].DailyActualIncome != nil {
			temp := util.Round(*data[i].DailyActualIncome, 2)
			data[i].DailyActualIncome = &temp
		}

		//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
		//需要取年月日(即前9位)
		{
			if data[i].Date != nil {
				temp := *data[i].Date
				*data[i].Date = temp[:10]
			}
		}
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return response.List{
		Data: data,
		Paging: &list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
