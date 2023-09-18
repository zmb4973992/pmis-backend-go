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

type IncomeAndExpenditureGet struct {
	Id int64
}

type IncomeAndExpenditureCreate struct {
	UserId int64
	//连接关联表的id
	ProjectId  int64 `json:"project_id,omitempty"`
	ContractId int64 `json:"contract_id,omitempty"`
	//连接dictionary_item表的id
	FundDirection string `json:"fund_direction,omitempty"`
	Currency      int64  `json:"currency,omitempty"`
	Kind          string `json:"kind,omitempty"`
	//日期
	Date string `json:"date,omitempty"`
	//数字
	Amount       *float64 `json:"amount"`
	ExchangeRate *float64 `json:"exchange_rate"`
	//字符串
	Type               string `json:"type,omitempty"`
	Term               string `json:"term,omitempty"`
	Remarks            string `json:"remarks,omitempty"`
	Attachment         string `json:"attachment,omitempty"`
	ImportedApprovalId string `json:"-"`
	DataSource         string `json:"-"`

	IgnoreUpdatingCumulativeIncomeAndExpenditure bool `json:"-"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type IncomeAndExpenditureUpdate struct {
	UserId int64
	Id     int64
	//连接关联表的id
	ProjectId  *int64 `json:"project_id"`
	ContractId *int64 `json:"contract_id"`
	//连接dictionary_item表的id
	FundDirection *string `json:"fund_direction"`
	Currency      *int64  `json:"currency"`
	Kind          *string `json:"kind"`
	//日期
	Date *string `json:"date"`
	//数字
	Amount       *float64 `json:"amount"`
	ExchangeRate *float64 `json:"exchange_rate"`
	//字符串
	Type       *string `json:"type"`
	Term       *string `json:"term"`
	Remarks    *string `json:"remarks"`
	Attachment *string `json:"attachment"`
}

type IncomeAndExpenditureDelete struct {
	Id int64
}

type IncomeAndExpenditureGetList struct {
	list.Input
	UserId        int64  `json:"-"`
	ProjectId     int64  `json:"project_id,omitempty"`
	Kind          string `json:"kind,omitempty"`
	FundDirection string `json:"fund_direction,omitempty"`
	DateGte       string `json:"date_gte,omitempty"`
	DateLte       string `json:"date_lte,omitempty"`
}

//以下为出参

type IncomeAndExpenditureOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	ProjectId  *int64 `json:"-"`
	ContractId *int64 `json:"-"`
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示
	FundDirection *int64 `json:"-"`
	Currency      *int64 `json:"-"`
	Kind          *int64 `json:"-"`
	DataSource    *int64 `json:"-"`
	//关联表的详情，不需要gorm查询，需要在json中显示
	ProjectExternal  *ProjectOutput  `json:"project" gorm:"-"`
	ContractExternal *ContractOutput `json:"contract" gorm:"-"`
	//dictionary_item表的详情，不需要gorm查询，需要在json中显示
	FundDirectionExternal *DictionaryDetailOutput `json:"fund_direction" gorm:"-"`
	CurrencyExternal      *DictionaryDetailOutput `json:"currency" gorm:"-"`
	KindExternal          *DictionaryDetailOutput `json:"kind" gorm:"-"`
	DataSourceExternal    *DictionaryDetailOutput `json:"data_source" gorm:"-"`
	//其他属性
	Date *string `json:"date"`

	Amount       *float64 `json:"amount"`
	ExchangeRate *float64 `json:"exchange_rate"`

	Type       *string `json:"type"`
	Term       *string `json:"term"`
	Remarks    *string `json:"remarks"`
	Attachment *string `json:"attachment"`
}

func (i *IncomeAndExpenditureGet) Get() (output *IncomeAndExpenditureOutput, errCode int) {
	err := global.DB.Model(model.IncomeAndExpenditure{}).
		Where("id = ?", i.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}
	//查询关联表的详情
	{
		//查项目信息
		if output.ProjectId != nil {
			var record ProjectOutput
			res := global.DB.Model(&model.Project{}).
				Where("id = ?", *output.ProjectId).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.ProjectExternal = &record
			}
		}
		//查合同信息
		if output.ContractId != nil {
			var record ContractOutput
			res := global.DB.Model(&model.Contract{}).
				Where("id = ?", *output.ContractId).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.ContractExternal = &record
			}
		}
	}

	//查询dictionary_item表的详情
	{
		if output.FundDirection != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.FundDirection).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.FundDirectionExternal = &record
			}
		}
		if output.Currency != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.Currency).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.CurrencyExternal = &record
			}
		}
		if output.Kind != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.Kind).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.KindExternal = &record
			}
		}
		if output.DataSource != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.DataSource).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.DataSourceExternal = &record
			}
		}
	}

	//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
	//需要取年月日(即前9位)
	{
		if output.Date != nil {
			temp := *output.Date
			*output.Date = temp[:10]
		}
	}

	return output, util.Success
}

func (i *IncomeAndExpenditureCreate) Create() (errCode int) {
	var paramOut model.IncomeAndExpenditure

	if i.UserId > 0 {
		paramOut.Creator = &i.UserId
	}

	//连接关联表的id
	{
		if i.ProjectId > 0 {
			paramOut.ProjectId = &i.ProjectId
		}
		if i.ContractId > 0 {
			paramOut.ContractId = &i.ContractId
		}
	}

	//连接dictionary_item表的id
	{
		if i.FundDirection != "" {
			var fundDirection model.DictionaryDetail
			err := global.DB.Where("name = ?", i.FundDirection).
				First(&fundDirection).Error
			if err != nil {
				return util.ErrorFailToCreateRecord
			}
			paramOut.FundDirection = &fundDirection.Id
		}
		if i.Currency > 0 {
			paramOut.Currency = &i.Currency
		}
		if i.Kind != "" {
			var kind model.DictionaryDetail
			err := global.DB.Where("name = ?", i.Kind).
				First(&kind).Error
			if err != nil {
				return util.ErrorFailToCreateRecord
			}
			paramOut.Kind = &kind.Id
		}
	}

	//日期
	{
		if i.Date != "" {
			date, err := time.Parse("2006-01-02", i.Date)
			if err != nil {
				return util.ErrorInvalidDateFormat
			}
			paramOut.Date = &date
		}
	}

	//数字
	{
		if i.Amount != nil {
			paramOut.Amount = i.Amount
		}
		if i.ExchangeRate != nil {
			paramOut.ExchangeRate = i.ExchangeRate
		}
	}

	//字符串
	{
		if i.Term != "" {
			paramOut.Term = &i.Term
		}
		if i.Type != "" {
			paramOut.Type = &i.Type
		}
		if i.Remarks != "" {
			paramOut.Remarks = &i.Remarks
		}
		if i.Attachment != "" {
			paramOut.Attachment = &i.Attachment
		}
		if i.ImportedApprovalId != "" {
			paramOut.ImportedApprovalId = &i.ImportedApprovalId
		}
		if i.DataSource != "" {
			var dataSource model.DictionaryType
			err := global.DB.Where("name = ?", "收款的数据来源").
				First(&dataSource).Error
			if err != nil {
				return util.ErrorFailToCreateRecord
			}

			var detailedDataSource model.DictionaryDetail
			err = global.DB.Where("dictionary_type_id = ?", dataSource.Id).
				Where("name = ?", i.DataSource).
				First(&detailedDataSource).Error
			if err != nil {
				return util.ErrorFailToCreateRecord
			}

			paramOut.DataSource = &detailedDataSource.Id
		}
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}

	if i.IgnoreUpdatingCumulativeIncomeAndExpenditure == false {
		//更新项目的累计收付款
		if i.ProjectId > 0 {
			temp1 := ProjectDailyAndCumulativeExpenditureUpdate{
				UserId:    i.UserId,
				ProjectId: i.ProjectId,
			}
			temp1.Update()
			temp2 := ProjectDailyAndCumulativeIncomeUpdate{
				UserId:    i.UserId,
				ProjectId: i.ProjectId,
			}
			temp2.Update()
		}

		//更新合同的累计收付款
		if i.ContractId > 0 {
			temp3 := ContractDailyAndCumulativeExpenditureUpdate{
				UserId:     i.UserId,
				ContractId: i.ContractId,
			}
			temp3.Update()
			temp4 := ContractDailyAndCumulativeIncomeUpdate{
				UserId:     i.UserId,
				ContractId: i.ContractId,
			}
			temp4.Update()
		}
	}

	return util.Success
}

func (i *IncomeAndExpenditureUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	if i.UserId > 0 {
		paramOut["last_modifier"] = i.UserId
	}

	//连接关联表的id
	{
		if i.ProjectId != nil {
			if *i.ProjectId > 0 {
				paramOut["project_id"] = *i.ProjectId
			}
		}
		if i.ContractId != nil {
			if *i.ContractId > 0 {
				paramOut["contract_id"] = i.ContractId
			} else if *i.ContractId == -1 {
				paramOut["contract_id"] = nil
			}
		}
	}

	//连接dictionary_item表的id
	{
		if i.FundDirection != nil {
			if *i.FundDirection != "" {
				var fundDirection model.DictionaryDetail
				err := global.DB.Where("name = ?", i.FundDirection).
					First(&fundDirection).Error
				if err != nil {
					return util.ErrorFailToUpdateRecord
				}
				paramOut["fund_direction"] = fundDirection.Id
			} else if *i.FundDirection == "" {
				paramOut["fund_direction"] = nil
			}
		}
		if i.Currency != nil {
			if *i.Currency > 0 {
				paramOut["currency"] = i.Currency
			} else if *i.Currency == -1 {
				paramOut["currency"] = nil
			}
		}
		if i.Kind != nil {
			if *i.Kind != "" {
				var kind model.DictionaryDetail
				err := global.DB.Where("name = ?", i.Kind).
					First(&kind).Error
				if err != nil {
					return util.ErrorFailToUpdateRecord
				}
				paramOut["kind"] = kind.Id
			} else if *i.Kind == "" {
				paramOut["kind"] = nil
			}
		}
	}

	//日期
	{
		if i.Date != nil {
			if *i.Date != "" {
				var err error
				paramOut["date"], err = time.Parse("2006-01-02", *i.Date)
				if err != nil {
					return util.ErrorInvalidJSONParameters
				}
			} else {
				paramOut["date"] = nil
			}
		}
	}

	//数字
	{
		if i.Amount != nil {
			if *i.Amount != -1 {
				paramOut["amount"] = i.Amount
			} else {
				paramOut["amount"] = nil
			}
		}
		if i.ExchangeRate != nil {
			if *i.ExchangeRate != -1 {
				paramOut["exchange_rate"] = i.ExchangeRate
			} else {
				paramOut["exchange_rate"] = nil
			}
		}
	}

	//字符串
	{
		if i.Type != nil {
			if *i.Type != "" {
				paramOut["type"] = *i.Type
			} else {
				paramOut["type"] = nil
			}
		}
		if i.Term != nil {
			if *i.Term != "" {
				paramOut["term"] = *i.Term
			} else {
				paramOut["term"] = nil
			}
		}
		if i.Remarks != nil {
			if *i.Remarks != "" {
				paramOut["remarks"] = *i.Remarks
			} else {
				paramOut["remarks"] = nil
			}
		}
		if i.Attachment != nil {
			if *i.Attachment != "" {
				paramOut["attachment"] = *i.Attachment
			} else {
				paramOut["attachment"] = nil
			}
		}
	}

	err := global.DB.Model(&model.IncomeAndExpenditure{}).
		Where("id = ?", i.Id).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	//更新项目的累计收付款
	if i.ProjectId == nil {
		//如果没有修改projectId,就用原纪录的projectId
		var record model.IncomeAndExpenditure
		err = global.DB.Where("id = ?", i.Id).
			First(&record).Error
		if err == nil && record.ProjectId != nil {
			temp1 := ProjectDailyAndCumulativeExpenditureUpdate{ProjectId: *record.ProjectId}
			temp1.Update()
			temp2 := ProjectDailyAndCumulativeIncomeUpdate{ProjectId: *record.ProjectId}
			temp2.Update()
		}
	} else {
		temp1 := ProjectDailyAndCumulativeExpenditureUpdate{ProjectId: *i.ProjectId}
		temp1.Update()
		temp2 := ProjectDailyAndCumulativeIncomeUpdate{ProjectId: *i.ProjectId}
		temp2.Update()
	}

	//更新合同的累计收付款
	if i.ContractId == nil {
		//如果没有修改contractId,就用原纪录的contractId
		var record model.IncomeAndExpenditure
		err = global.DB.Where("id = ?", i.Id).
			First(&record).Error
		if err == nil && record.ContractId != nil {
			temp3 := ContractDailyAndCumulativeExpenditureUpdate{ContractId: *record.ContractId}
			temp3.Update()
			temp4 := ContractDailyAndCumulativeIncomeUpdate{ContractId: *record.ContractId}
			temp4.Update()
		}
	} else {
		temp3 := ContractDailyAndCumulativeExpenditureUpdate{ContractId: *i.ContractId}
		temp3.Update()
		temp4 := ContractDailyAndCumulativeIncomeUpdate{ContractId: *i.ContractId}
		temp4.Update()
	}

	return util.Success
}

func (i *IncomeAndExpenditureDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.IncomeAndExpenditure
	err := global.DB.Where("id = ?", i.Id).
		Find(&record).
		Delete(&record).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	//更新项目的累计收付款
	if record.ProjectId != nil {
		temp1 := ProjectDailyAndCumulativeExpenditureUpdate{ProjectId: *record.ProjectId}
		temp1.Update()
		temp2 := ProjectDailyAndCumulativeIncomeUpdate{ProjectId: *record.ProjectId}
		temp2.Update()
	}

	//更新合同的累计收付款
	if record.ProjectId != nil {
		temp3 := ContractDailyAndCumulativeExpenditureUpdate{ContractId: *record.ContractId}
		temp3.Update()
		temp4 := ContractDailyAndCumulativeIncomeUpdate{ContractId: *record.ContractId}
		temp4.Update()
	}

	return util.Success
}

func (i *IncomeAndExpenditureGetList) GetList() (
	outputs []IncomeAndExpenditureOutput, errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.IncomeAndExpenditure{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where
	if i.ProjectId > 0 {
		db = db.Where("project_id = ?", i.ProjectId)
	}

	if i.Kind != "" {
		var dictionaryDetail model.DictionaryDetail
		err := global.DB.Model(&model.DictionaryDetail{}).
			Where("name = ?", i.Kind).
			First(&dictionaryDetail).Error
		if err != nil {
			return nil, util.ErrorRecordNotFound, nil
		}
		db = db.Where("kind = ?", dictionaryDetail.Id)
	}

	if i.FundDirection != "" {
		var dictionaryDetail model.DictionaryDetail
		err := global.DB.Model(&model.DictionaryDetail{}).
			Where("name = ?", i.FundDirection).
			First(&dictionaryDetail).Error
		if err != nil {
			return nil, util.ErrorRecordNotFound, nil
		}
		db = db.Where("fund_direction = ?", dictionaryDetail.Id)
	}

	if i.DateGte != "" {
		db = db.Where("date >= ?", i.DateGte)
	}

	if i.DateLte != "" {
		db = db.Where("date <= ?", i.DateLte)
	}

	//用来确定组织的数据范围
	organizationIds := util.GetOrganizationIdsForDataAuthority(i.UserId)

	db = db.Joins("join (select distinct income_and_expenditure.id as income_and_expenditure_id from income_and_expenditure join (select distinct contract.id as contract_id from contract join (select distinct project.id as project_id from project where organization_id in ?) as temp1 on contract.project_id = temp1.project_id) as temp2  on income_and_expenditure.contract_id = temp2.contract_id union select distinct income_and_expenditure.id as income_and_expenditure_id from income_and_expenditure join (select distinct project.id as project_id from project where organization_id in ?) as temp2 on income_and_expenditure.project_id = temp2.project_id) as temp3 on income_and_expenditure.id = temp3.income_and_expenditure_id", organizationIds, organizationIds)
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
		exists := util.FieldIsInModel(&model.IncomeAndExpenditure{}, orderBy)
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
	if i.PagingInput.Page > 0 {
		page = i.PagingInput.Page
	}
	pageSize := global.Config.Paging.DefaultPageSize
	if i.PagingInput.PageSize != nil && *i.PagingInput.PageSize >= 0 &&
		*i.PagingInput.PageSize <= global.Config.Paging.MaxPageSize {
		pageSize = *i.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.IncomeAndExpenditure{}).Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	for i := range outputs {
		//查询关联表的详情
		{
			//查项目信息
			if outputs[i].ProjectId != nil {
				var record ProjectOutput
				res := global.DB.Model(&model.Project{}).
					Where("id = ?", *outputs[i].ProjectId).
					Limit(1).
					Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].ProjectExternal = &record
				}
			}
			//查合同信息
			if outputs[i].ContractId != nil {
				var record ContractOutput
				res := global.DB.Model(&model.Contract{}).
					Where("id = ?", *outputs[i].ContractId).
					Limit(1).
					Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].ContractExternal = &record
				}
			}
		}

		//查dictionary_item表的详情
		{
			if outputs[i].FundDirection != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *outputs[i].FundDirection).
					Limit(1).
					Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].FundDirectionExternal = &record
				}
			}
			if outputs[i].Currency != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *outputs[i].Currency).
					Limit(1).
					Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].CurrencyExternal = &record
				}
			}
			if outputs[i].Kind != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *outputs[i].Kind).
					Limit(1).
					Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].KindExternal = &record
				}
			}
			if outputs[i].DataSource != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *outputs[i].DataSource).
					Limit(1).
					Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].DataSourceExternal = &record
				}
			}
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
