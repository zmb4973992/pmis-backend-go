package service

import (
	"github.com/yitter/idgenerator-go/idgen"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type IncomeAndExpenditureGet struct {
	ID int64
}

type IncomeAndExpenditureCreate struct {
	Creator      int64
	LastModifier int64
	//连接关联表的id
	ProjectID  int64 `json:"project_id,omitempty"`
	ContractID int64 `json:"contract_id,omitempty"`
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
	Type               int64  `json:"type,omitempty"`
	Term               int64  `json:"term,omitempty"`
	Remarks            string `json:"remarks,omitempty"`
	Attachment         string `json:"attachment,omitempty"`
	ImportedApprovalID string `json:"imported_approval_id,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type IncomeAndExpenditureUpdate struct {
	LastModifier int64
	ID           int64
	//连接关联表的id
	ProjectID  *int64 `json:"project_id"`
	ContractID *int64 `json:"contract_id"`
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
	Type       *int64  `json:"type"`
	Term       *int64  `json:"term"`
	Remarks    *string `json:"remarks"`
	Attachment *string `json:"attachment"`
}

type IncomeAndExpenditureDelete struct {
	ID int64
}

type IncomeAndExpenditureGetList struct {
	list.Input
	list.DataScopeInput
	ProjectID     int64  `json:"project_id,omitempty"`
	Kind          string `json:"kind,omitempty"`
	FundDirection string `json:"fund_direction,omitempty"`
	DateGte       string `json:"date_gte,omitempty"`
	DateLte       string `json:"date_lte,omitempty"`
}

//以下为出参

type IncomeAndExpenditureOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	ProjectID  *int64 `json:"-"`
	ContractID *int64 `json:"-"`
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示
	FundDirection *int64 `json:"-"`
	Currency      *int64 `json:"-"`
	Kind          *int64 `json:"-"`
	Type          *int64 `json:"-"`

	//关联表的详情，不需要gorm查询，需要在json中显示
	ProjectExternal  *ProjectOutput  `json:"project" gorm:"-"`
	ContractExternal *ContractOutput `json:"contract" gorm:"-"`
	//dictionary_item表的详情，不需要gorm查询，需要在json中显示
	FundDirectionExternal *DictionaryDetailOutput `json:"fund_direction" gorm:"-"`
	CurrencyExternal      *DictionaryDetailOutput `json:"currency" gorm:"-"`
	KindExternal          *DictionaryDetailOutput `json:"kind" gorm:"-"`
	TypeExternal          *DictionaryDetailOutput `json:"type" gorm:"-"`

	//其他属性
	Date *string `json:"date"`

	Amount       *float64 `json:"amount"`
	ExchangeRate *float64 `json:"exchange_rate"`

	Term       *string `json:"term"`
	Remarks    *string `json:"remarks"`
	Attachment *string `json:"attachment"`
}

func (i *IncomeAndExpenditureGet) Get() response.Common {
	var result IncomeAndExpenditureOutput
	err := global.DB.Model(model.IncomeAndExpenditure{}).
		Where("id = ?", i.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	//查询关联表的详情
	{
		//查项目信息
		if result.ProjectID != nil {
			var record ProjectOutput
			res := global.DB.Model(&model.Project{}).
				Where("id = ?", *result.ProjectID).Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.ProjectExternal = &record
			}
		}
		//查合同信息
		if result.ContractID != nil {
			var record ContractOutput
			res := global.DB.Model(&model.Contract{}).
				Where("id = ?", *result.ContractID).Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.ContractExternal = &record
			}
		}
	}

	//查询dictionary_item表的详情
	{
		if result.FundDirection != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.FundDirection).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.FundDirectionExternal = &record
			}
		}
		if result.Currency != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.Currency).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.CurrencyExternal = &record
			}
		}
		if result.Kind != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.Kind).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.KindExternal = &record
			}
		}
		if result.Type != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.Type).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.TypeExternal = &record
			}
		}
	}

	//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
	//需要取年月日(即前9位)
	{
		if result.Date != nil {
			temp := *result.Date
			*result.Date = temp[:10]
		}
	}

	return response.SuccessWithData(result)
}

func (i *IncomeAndExpenditureCreate) Create() response.Common {
	var paramOut model.IncomeAndExpenditure

	if i.Creator > 0 {
		paramOut.Creator = &i.Creator
	}
	if i.LastModifier > 0 {
		paramOut.LastModifier = &i.LastModifier
	}

	//连接关联表的id
	{
		if i.ProjectID > 0 {
			paramOut.ProjectID = &i.ProjectID
		}
		if i.ContractID > 0 {
			paramOut.ContractID = &i.ContractID
		}
	}

	//连接dictionary_item表的id
	{
		if i.FundDirection != "" {
			var fundDirection model.DictionaryDetail
			err := global.DB.Where("name = ?", i.FundDirection).
				First(&fundDirection).Error
			if err != nil {
				return response.Failure(util.ErrorFailToCreateRecord)
			}
			paramOut.FundDirection = &fundDirection.ID
		}
		if i.Currency > 0 {
			paramOut.Currency = &i.Currency
		}
		if i.Kind != "" {
			var kind model.DictionaryDetail
			err := global.DB.Where("name = ?", i.Kind).
				First(&kind).Error
			if err != nil {
				return response.Failure(util.ErrorFailToCreateRecord)
			}
			paramOut.Kind = &kind.ID
		}
	}

	//日期
	{
		if i.Date != "" {
			date, err := time.Parse("2006-01-02", i.Date)
			if err != nil {
				return response.Failure(util.ErrorInvalidDateFormat)
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
		if i.Type > 0 {
			paramOut.Type = &i.Type
		}
		if i.Term > 0 {
			paramOut.Term = &i.Term
		}
		if i.Remarks != "" {
			paramOut.Remarks = &i.Remarks
		}
		if i.Attachment != "" {
			paramOut.Attachment = &i.Attachment
		}
		if i.ImportedApprovalID != "" {
			paramOut.ImportedApprovalID = &i.ImportedApprovalID
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	//更新项目累计收付款
	var temp1 = ProjectCumulativeExpenditureUpdate{ProjectID: i.ProjectID}
	temp1.Update()
	var temp2 = ProjectCumulativeIncomeUpdate{ProjectID: i.ProjectID}
	temp2.Update()

	return response.Success()
}

func (i *IncomeAndExpenditureUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if i.LastModifier > 0 {
		paramOut["last_modifier"] = i.LastModifier
	}

	//连接关联表的id
	{
		if i.ProjectID != nil {
			if *i.ProjectID > 0 {
				paramOut["project_id"] = *i.ProjectID
			}
		}
		if i.ContractID != nil {
			if *i.ContractID > 0 {
				paramOut["contract_id"] = i.ContractID
			} else if *i.ContractID == -1 {
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
					return response.Failure(util.ErrorFailToUpdateRecord)
				}
				paramOut["fund_direction"] = fundDirection.ID
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
					return response.Failure(util.ErrorFailToUpdateRecord)
				}
				paramOut["kind"] = kind.ID
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
					return response.Failure(util.ErrorInvalidJSONParameters)
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
			if *i.Type > 0 {
				paramOut["type"] = *i.Type
			} else {
				paramOut["type"] = nil
			}
		}
		if i.Term != nil {
			if *i.Term > 0 {
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

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "Deleter", "CreateAt", "UpdatedAt", "DeletedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.IncomeAndExpenditure{}).Where("id = ?", i.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (i *IncomeAndExpenditureDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.IncomeAndExpenditure
	global.DB.Where("id = ?", i.ID).Find(&record)
	err := global.DB.Where("id = ?", i.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (i *IncomeAndExpenditureGetList) GetList() response.List {
	db := global.DB.Model(&model.IncomeAndExpenditure{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if i.ProjectID > 0 {
		db = db.Where("project_id = ?", i.ProjectID)
	}

	if i.Kind != "" {
		var dictionaryDetail model.DictionaryDetail
		err := global.DB.Model(&model.DictionaryDetail{}).
			Where("name = ?", i.Kind).First(&dictionaryDetail).Error
		if err != nil {
			return response.FailureForList(util.ErrorRecordNotFound)
		}
		db = db.Where("kind = ?", dictionaryDetail.ID)
	}

	if i.FundDirection != "" {
		var dictionaryDetail model.DictionaryDetail
		err := global.DB.Model(&model.DictionaryDetail{}).
			Where("name = ?", i.FundDirection).First(&dictionaryDetail).Error
		if err != nil {
			return response.FailureForList(util.ErrorRecordNotFound)
		}
		db = db.Where("fund_direction = ?", dictionaryDetail.ID)
	}

	if i.DateGte != "" {
		db = db.Where("date >= ?", i.DateGte)
	}

	if i.DateLte != "" {
		db = db.Where("date <= ?", i.DateLte)
	}

	batchID := idgen.NextId()

	//用来确定组织的数据范围
	organizationIDsForDataScope := util.GetOrganizationIDsInDataScope(i.UserID)
	if len(organizationIDsForDataScope) > 0 {
		var temps []model.Temp
		for j := range organizationIDsForDataScope {
			var temp model.Temp
			temp.OrganizationID = &organizationIDsForDataScope[j]
			temp.BatchID = batchID
			temps = append(temps, temp)
		}
		global.DB.CreateInBatches(&temps, 100)
	}

	//找出项目的数据范围
	var projectIDs []int64
	global.DB.Model(&model.Project{}).
		Joins("join temp on project.organization_id = temp.organization_id").
		Where("batch_id = ?", batchID).
		Select("project.id").Find(&projectIDs)
	if len(projectIDs) > 0 {
		var temps []model.Temp
		for j := range organizationIDsForDataScope {
			var temp model.Temp
			temp.ProjectID = &projectIDs[j]
			temp.BatchID = batchID
			temps = append(temps, temp)
		}
		global.DB.CreateInBatches(&temps, 100)
	}

	//然后再找出合同的数据范围
	var contractIDs []int64
	global.DB.Model(&model.Contract{}).
		Joins("join temp on contract.project_id = temp.project_id").
		Where("temp.batch_id = ?", batchID).
		Select("contract.id").Find(&contractIDs)
	if len(contractIDs) > 0 {
		var temps []model.Temp
		for j := range contractIDs {
			var temp model.Temp
			temp.ContractID = &contractIDs[j]
			temp.BatchID = batchID
			temps = append(temps, temp)
		}
		global.DB.CreateInBatches(&temps, 100)
	}

	//汇总
	db = db.Joins("join temp on income_and_expenditure.project_id = temp.project_id or income_and_expenditure.contract_id = temp.contract_id").
		Where("temp.batch_id = ?", batchID)

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
	var data []IncomeAndExpenditureOutput
	db.Model(&model.IncomeAndExpenditure{}).Find(&data)

	//temp使用完毕，需要删除数据
	global.DB.Where("batch_id = ?", batchID).Delete(&model.Temp{})

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	for i := range data {
		//查询关联表的详情
		{
			//查项目信息
			if data[i].ProjectID != nil {
				var record ProjectOutput
				res := global.DB.Model(&model.Project{}).
					Where("id = ?", *data[i].ProjectID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].ProjectExternal = &record
				}
			}
			//查合同信息
			if data[i].ContractID != nil {
				var record ContractOutput
				res := global.DB.Model(&model.Contract{}).
					Where("id = ?", *data[i].ContractID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].ContractExternal = &record
				}
			}
		}

		//查dictionary_item表的详情
		{
			if data[i].FundDirection != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *data[i].FundDirection).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].FundDirectionExternal = &record
				}
			}
			if data[i].Currency != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *data[i].Currency).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].CurrencyExternal = &record
				}
			}
			if data[i].Kind != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *data[i].Kind).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].KindExternal = &record
				}
			}
			if data[i].Type != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *data[i].Type).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].TypeExternal = &record
				}
			}
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
