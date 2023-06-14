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

type CumulativeContractIncomeAndExpenditureUpdate struct {
	Creator      int64
	LastModifier int64
	//连接关联表的id
	ContractID int64 `json:"contract_id,omitempty"`
	//连接dictionary_item表的id

	//日期
	//数字

	//字符串
}

type CumulativeContractIncomeAndExpenditureGetList struct {
	list.Input
	ContractID int64  `json:"contract_id,omitempty"`
	DateGte    string `json:"date_gte,omitempty"`
	DateLte    string `json:"date_lte,omitempty"`
}

//以下为出参

type CumulativeContractIncomeAndExpenditureOutput struct {
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

	TotalPlannedExpenditure       *float64 `json:"total_planned_expenditure"`       //计划付款总额
	TotalActualExpenditure        *float64 `json:"total_actual_expenditure"`        //实际付款总额
	TotalForecastedExpenditure    *float64 `json:"total_forecasted_expenditure"`    //预测付款总额
	PlannedExpenditureProgress    *float64 `json:"planned_expenditure_progress"`    //计划付款进度
	ActualExpenditureProgress     *float64 `json:"actual_expenditure_progress"`     //实际付款进度
	ForecastedExpenditureProgress *float64 `json:"forecasted_expenditure_progress"` //预测付款进度
	TotalPlannedIncome            *float64 `json:"total_planned_income"`            //计划收款总额
	TotalActualIncome             *float64 `json:"total_actual_income"`             //实际收款总额
	TotalForecastedIncome         *float64 `json:"total_forecasted_income"`         //预测收款总额
	PlannedIncomeProgress         *float64 `json:"planned_income_progress"`         //计划收款进度
	ActualIncomeProgress          *float64 `json:"actual_income_progress"`          //实际收款进度
	ForecastedIncomeProgress      *float64 `json:"forecasted_income_progress"`      //预测收款进度

	//其他属性

}

func (i *CumulativeContractIncomeAndExpenditureUpdate) Update() response.Common {
	//连接关联表的id
	{
		if i.ContractID > 0 {
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

			var expenditure int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("name = '付款'").Select("id").First(&expenditure).Error
			if err != nil {
				return response.Failure(util.ErrorFailToUpdateRecord)
			}
			fmt.Println("expenditure id:", expenditure)

			//计算合同的总金额
			var totalAmountOfContract float64
			err = global.DB.Model(&model.Contract{}).
				Where("id = ?", i.ContractID).
				Select("coalesce(amount * exchange_rate,0)").
				Find(&totalAmountOfContract).Error
			fmt.Println("合同总金额：", totalAmountOfContract)

			global.DB.Where("contract_id = ?", i.ContractID).
				Delete(&model.CumulativeContractIncomeAndExpenditure{})

			var datesForIncome []time.Time
			global.DB.Model(&model.IncomeAndExpenditure{}).
				Where("contract_id = ?", i.ContractID).
				Where("fund_direction = ?", income).
				Select("date").
				Distinct("date").Order("date desc").Find(&datesForIncome)
			fmt.Println("收款日期：", datesForIncome)

			var datesForExpenditure []time.Time
			global.DB.Model(&model.IncomeAndExpenditure{}).
				Where("contract_id = ?", i.ContractID).
				Where("fund_direction = ?", expenditure).
				Select("date").
				Distinct("date").Order("date desc").Find(&datesForExpenditure)
			fmt.Println("付款日期：", datesForExpenditure)

			var records []model.CumulativeContractIncomeAndExpenditure

			for j := range datesForIncome {
				var totalPlannedIncome float64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("contract_id = ?", i.ContractID).
					Where("kind = ?", planned).
					Where("fund_direction = ?", income).
					Where("date <= ?", datesForIncome[j]).
					Select("coalesce(sum(amount * exchange_rate),0)").
					Find(&totalPlannedIncome)
				fmt.Println("计划收款总额：", totalPlannedIncome)

				var totalActualIncome float64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("contract_id = ?", i.ContractID).
					Where("kind = ?", actual).
					Where("fund_direction = ?", income).
					Where("date <= ?", datesForIncome[j]).
					Select("coalesce(sum(amount * exchange_rate),0)").
					Find(&totalActualIncome)
				fmt.Println("实际收款总额：", totalActualIncome)

				var totalForecastedIncome float64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("contract_id = ?", i.ContractID).
					Where("kind = ?", forecasted).
					Where("fund_direction = ?", income).
					Where("date <= ?", datesForIncome[j]).
					Select("coalesce(sum(amount * exchange_rate),0)").
					Find(&totalForecastedIncome)
				fmt.Println("预测收款总额：", totalForecastedIncome)

				var record model.CumulativeContractIncomeAndExpenditure
				record.Creator = &i.Creator
				record.LastModifier = &i.LastModifier
				record.ContractID = i.ContractID
				record.Date = &datesForIncome[j]

				record.TotalPlannedIncome = totalPlannedIncome
				record.TotalActualIncome = totalActualIncome
				record.TotalForecastedIncome = totalForecastedIncome

				if totalAmountOfContract == 0 {
					record.PlannedIncomeProgress = 0
					record.ActualIncomeProgress = 0
					record.ForecastedIncomeProgress = 0
				} else {
					record.PlannedIncomeProgress = totalPlannedIncome / totalAmountOfContract
					record.ActualIncomeProgress = totalActualIncome / totalAmountOfContract
					record.ForecastedIncomeProgress = totalForecastedIncome / totalAmountOfContract
				}

				fmt.Println("------------------------")
				records = append(records, record)
			}

			for j := range datesForExpenditure {
				var totalPlannedExpenditure float64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("contract_id = ?", i.ContractID).
					Where("kind = ?", planned).
					Where("fund_direction = ?", expenditure).
					Where("date <= ?", datesForExpenditure[j]).
					Select("coalesce(sum(amount * exchange_rate),0)").
					Find(&totalPlannedExpenditure)
				fmt.Println("计划付款总额：", totalPlannedExpenditure)

				var totalActualExpenditure float64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("contract_id = ?", i.ContractID).
					Where("kind = ?", actual).
					Where("fund_direction = ?", expenditure).
					Where("date <= ?", datesForExpenditure[j]).
					Select("coalesce(sum(amount * exchange_rate),0)").
					Find(&totalActualExpenditure)
				fmt.Println("实际付款总额：", totalActualExpenditure)

				var totalForecastedExpenditure float64
				global.DB.Model(&model.IncomeAndExpenditure{}).
					Where("contract_id = ?", i.ContractID).
					Where("kind = ?", forecasted).
					Where("fund_direction = ?", expenditure).
					Where("date <= ?", datesForExpenditure[j]).
					Select("coalesce(sum(amount * exchange_rate),0)").
					Find(&totalForecastedExpenditure)
				fmt.Println("预测付款总额：", totalForecastedExpenditure)

				var record model.CumulativeContractIncomeAndExpenditure
				record.Creator = &i.Creator
				record.LastModifier = &i.LastModifier
				record.ContractID = i.ContractID
				record.Date = &datesForExpenditure[j]

				record.TotalPlannedExpenditure = totalPlannedExpenditure
				record.TotalActualExpenditure = totalActualExpenditure
				record.TotalForecastedExpenditure = totalForecastedExpenditure

				if totalAmountOfContract == 0 {
					record.PlannedExpenditureProgress = 0
					record.ActualExpenditureProgress = 0
					record.ForecastedExpenditureProgress = 0
				} else {
					record.PlannedExpenditureProgress = totalPlannedExpenditure / totalAmountOfContract
					record.ActualExpenditureProgress = totalActualExpenditure / totalAmountOfContract
					record.ForecastedExpenditureProgress = totalForecastedExpenditure / totalAmountOfContract
				}

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

func (i *CumulativeContractIncomeAndExpenditureGetList) GetList() response.List {
	db := global.DB.Model(&model.CumulativeContractIncomeAndExpenditure{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if i.ContractID > 0 {
		db = db.Where("contract_id = ?", i.ContractID)
	}

	if i.DateGte != "" {
		db = db.Where("date >= ?", i.DateGte)
	}

	if i.DateLte != "" {
		db = db.Where("date <= ?", i.DateLte)
	}

	//count
	var count int64
	db.Count(&count)

	//order
	orderBy := i.SortingInput.OrderBy
	desc := i.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.CumulativeContractIncomeAndExpenditure{}, orderBy)
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
	if i.PagingInput.Page > 0 {
		page = i.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if i.PagingInput.PageSize != nil && *i.PagingInput.PageSize >= 0 &&
		*i.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *i.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []CumulativeContractIncomeAndExpenditureOutput
	db.Model(&model.CumulativeContractIncomeAndExpenditure{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	for i := range data {
		//查询关联表的详情
		{
			//查项目信息
			//if data[i].ProjectID != nil {
			//	var record ProjectOutput
			//	res := global.DB.Model(&model.Project{}).
			//		Where("id = ?", *data[i].ProjectID).Limit(1).Find(&record)
			//	if res.RowsAffected > 0 {
			//		data[i].ProjectExternal = &record
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