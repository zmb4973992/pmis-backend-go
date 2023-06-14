package service

import (
	"fmt"
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

type ProjectCumulativeIncomeUpdate struct {
	Creator      int64
	LastModifier int64
	//连接关联表的id
	ProjectID int64 `json:"project_id,omitempty"`
	//连接dictionary_item表的id

	//日期
	//数字

	//字符串
}

type ProjectCumulativeIncomeGetList struct {
	list.Input
	ProjectID int64  `json:"project_id,omitempty"`
	DateGte   string `json:"date_gte,omitempty"`
	DateLte   string `json:"date_lte,omitempty"`
}

//以下为出参

type ProjectCumulativeIncomeOutput struct {
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

	TotalPlannedIncome       *float64 `json:"total_planned_income"`       //计划收款总额
	TotalActualIncome        *float64 `json:"total_actual_income"`        //实际收款总额
	TotalForecastedIncome    *float64 `json:"total_forecasted_income"`    //预测收款总额
	PlannedIncomeProgress    *float64 `json:"planned_income_progress"`    //计划收款进度
	ActualIncomeProgress     *float64 `json:"actual_income_progress"`     //实际收款进度
	ForecastedIncomeProgress *float64 `json:"forecasted_income_progress"` //预测收款进度

	//其他属性

}

func (p *ProjectCumulativeIncomeUpdate) Update() response.Common {
	//连接关联表的id
	{
		if p.ProjectID > 0 {
			var planned int64
			err := global.DB.Model(&model.DictionaryDetail{}).
				Where("name = '计划'").Select("id").First(&planned).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}
			fmt.Println("planned id:", planned)

			var actual int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("name = '实际'").Select("id").First(&actual).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}
			fmt.Println("actual id:", actual)

			var forecasted int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("name = '预测'").Select("id").First(&forecasted).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}
			fmt.Println("forecasted id:", forecasted)

			var income int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("name = '收款'").Select("id").First(&income).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}
			fmt.Println("income id:", income)

			//计算收款合同的总金额
			var totalAmountOfIncomeContract float64
			err = global.DB.Model(&model.Contract{}).
				Where("project_id = ?", p.ProjectID).
				Where("fund_direction = ?", income).
				Select("coalesce(sum(amount * exchange_rate),0)").
				Find(&totalAmountOfIncomeContract).Error
			fmt.Println("收款合同总金额：", totalAmountOfIncomeContract)
			fmt.Println("*********************************")

			global.DB.Where("project_id = ?", p.ProjectID).
				Delete(&model.ProjectCumulativeIncome{})

			var dates []time.Time
			global.DB.Model(&model.IncomeAndExpenditure{}).
				Where("project_id = ?", p.ProjectID).
				Where("fund_direction = ?", income).
				Select("date").Distinct("date").Order("date desc").
				Find(&dates)
			fmt.Println("日期：", dates)

			var records []model.ProjectCumulativeIncome

			for j := range dates {
				fmt.Println("日期：", dates[j].Format("2006-01-02")[:10])
				var record model.ProjectCumulativeIncome

				var totalPlannedIncome float64
				var countOfPlannedIncome int64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("project_id = ?", p.ProjectID).
					Where("kind = ?", planned).
					Where("fund_direction = ?", income).
					Where("date <= ?", dates[j]).
					Count(&countOfPlannedIncome)
				if countOfPlannedIncome > 0 {
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectID).
						Where("kind = ?", planned).
						Where("fund_direction = ?", income).
						Where("date <= ?", dates[j]).
						Select("coalesce(sum(amount * exchange_rate),0)").
						Find(&totalPlannedIncome)
					record.TotalPlannedIncome = &totalPlannedIncome
					fmt.Println("计划收款总额：", totalPlannedIncome)
					if totalAmountOfIncomeContract > 0 {
						var plannedIncomeProgress = totalPlannedIncome / totalAmountOfIncomeContract
						record.PlannedIncomeProgress = &plannedIncomeProgress
						fmt.Println("计划收款进度：", plannedIncomeProgress)
					}
				}

				var totalActualIncome float64
				var countOfActualIncome int64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("project_id = ?", p.ProjectID).
					Where("kind = ?", actual).
					Where("fund_direction = ?", income).
					Where("date <= ?", dates[j]).
					Count(&countOfActualIncome)
				if countOfActualIncome > 0 {
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectID).
						Where("kind = ?", actual).
						Where("fund_direction = ?", income).
						Where("date <= ?", dates[j]).
						Select("coalesce(sum(amount * exchange_rate),0)").
						Find(&totalActualIncome)
					record.TotalActualIncome = &totalActualIncome
					fmt.Println("实际收款总额：", totalActualIncome)
					if totalAmountOfIncomeContract > 0 {
						var actualIncomeProgress = totalActualIncome / totalAmountOfIncomeContract
						record.ActualIncomeProgress = &actualIncomeProgress
						fmt.Println("实际收款进度：", actualIncomeProgress)
					}
				}

				var totalForecastedIncome float64
				var countOfForecastedIncome int64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("project_id = ?", p.ProjectID).
					Where("kind = ?", forecasted).
					Where("fund_direction = ?", income).
					Where("date <= ?", dates[j]).
					Count(&countOfForecastedIncome)
				if countOfForecastedIncome > 0 {
					global.DB.Model(&model.IncomeAndExpenditure{}).
						Where("project_id = ?", p.ProjectID).
						Where("kind = ?", forecasted).
						Where("fund_direction = ?", income).
						Where("date <= ?", dates[j]).
						Select("coalesce(sum(amount * exchange_rate),0)").
						Find(&totalForecastedIncome)
					record.TotalForecastedIncome = &totalForecastedIncome
					fmt.Println("预测收款总额：", totalForecastedIncome)
					if totalAmountOfIncomeContract > 0 {
						var forecastedIncomeProgress = totalForecastedIncome / totalAmountOfIncomeContract
						record.ForecastedIncomeProgress = &forecastedIncomeProgress
						fmt.Println("预测收款进度：", forecastedIncomeProgress)
					}
				}

				record.Creator = &p.Creator
				record.LastModifier = &p.LastModifier
				record.ProjectID = p.ProjectID
				record.Date = &dates[j]

				fmt.Println("------------------------")

				records = append(records, record)
			}

			if len(records) == 0 {
				return response.Success()
			}

			err = global.DB.Create(&records).Error
			if err != nil {
				return response.Failure(util.ErrorFailToCreateRecord)
			}
		}
	}

	return response.Success()
}

func (p *ProjectCumulativeIncomeGetList) GetList() response.List {
	db := global.DB.Model(&model.ProjectCumulativeIncome{})
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
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.ProjectCumulativeIncome{}, orderBy)
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
	var data []ProjectCumulativeIncomeOutput
	db.Model(&model.ProjectCumulativeIncome{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	for i := range data {
		//查询关联表的详情
		{
			//查项目信息
			//if data[p].ProjectID != nil {
			//	var record ProjectOutput
			//	res := global.DB.Model(&model.Project{}).
			//		Where("id = ?", *data[p].ProjectID).Limit(1).Find(&record)
			//	if res.RowsAffected > 0 {
			//		data[p].ProjectExternal = &record
			//	}
			//}
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