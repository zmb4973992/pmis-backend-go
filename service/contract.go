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

type ContractGet struct {
	ID int64
}

type ContractCreate struct {
	Creator      int64
	LastModifier int64
	//连接关联表的id
	ProjectID      int64 `json:"project_id,omitempty"`
	OrganizationID int64 `json:"organization_id,omitempty"`
	RelatedPartyID int64 `json:"related_party_id,omitempty"`
	//连接dictionary_item表的id
	FundDirection int64 `json:"fund_direction,omitempty"`
	OurSignatory  int64 `json:"our_signatory,omitempty"`
	Currency      int64 `json:"currency,omitempty"`
	Type          int64 `json:"type,omitempty"`
	//日期
	SigningDate       string `json:"signing_date,omitempty"`
	EffectiveDate     string `json:"effective_date,omitempty"`
	CommissioningDate string `json:"commissioning_date,omitempty"`
	CompletionDate    string `json:"completion_date,omitempty"`
	//数字(允许为0、nil)
	Amount             *float64 `json:"amount"`
	ExchangeRate       *float64 `json:"exchange_rate"`
	ConstructionPeriod *int     `json:"construction_period"`
	//字符串(允许为nil)
	Name        string `json:"name,omitempty"`
	Code        string `json:"code,omitempty"`
	Content     string `json:"content,omitempty"`
	Deliverable string `json:"deliverable,omitempty"`
	PenaltyRule string `json:"penalty_rule,omitempty"`
	Attachment  string `json:"attachment,omitempty"`
	Operator    string `json:"operator,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ContractUpdate struct {
	LastModifier int64
	ID           int64
	//连接关联表的id
	ProjectID      *int64 `json:"project_id"`
	OrganizationID *int64 `json:"organization_id"`
	RelatedPartyID *int64 `json:"related_party_id"`
	//连接dictionary_item表的id
	FundDirection *int64 `json:"fund_direction"`
	OurSignatory  *int64 `json:"our_signatory"`
	Currency      *int64 `json:"currency"`
	Type          *int64 `json:"type"`
	//日期
	SigningDate       *string `json:"signing_date"`
	EffectiveDate     *string `json:"effective_date"`
	CommissioningDate *string `json:"commissioning_date"`
	CompletionDate    *string `json:"completion_date"`
	//允许为0的数字
	Amount             *float64 `json:"amount"`
	ExchangeRate       *float64 `json:"exchange_rate"`
	ConstructionPeriod *int     `json:"construction_period"`
	//允许为null的字符串
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Content     *string `json:"content"`
	Deliverable *string `json:"deliverable"`
	PenaltyRule *string `json:"penalty_rule"`
	Attachment  *string `json:"attachment"`
	Operator    *string `json:"operator"`
}

type ContractDelete struct {
	ID int64
}

type ContractGetList struct {
	list.Input
	list.DataScopeInput
	ProjectID      int64  `json:"project_id,omitempty"`
	RelatedPartyID int64  `json:"related_party_id,omitempty"`
	FundDirection  int64  `json:"fund_direction,omitempty"`
	NameInclude    string `json:"name_include,omitempty"`
}

//以下为出参

type ContractOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	ProjectID      *int64 `json:"-"`
	OrganizationID *int64 `json:"-"`
	RelatedPartyID *int64 `json:"-"`
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示
	FundDirection *int64 `json:"-"`
	OurSignatory  *int64 `json:"-"`
	Currency      *int64 `json:"-"`
	Type          *int64 `json:"-"`
	//关联表的详情，不需要gorm查询，需要在json中显示
	ProjectExternal      *ProjectOutput      `json:"project" gorm:"-"`
	OrganizationExternal *OrganizationOutput `json:"organization" gorm:"-"`
	RelatedPartyExternal *RelatedPartyOutput `json:"related_party" gorm:"-"`
	//dictionary_item表的详情，不需要gorm查询，需要在json中显示
	FundDirectionExternal *DictionaryDetailOutput `json:"fund_direction" gorm:"-"`
	OurSignatoryExternal  *DictionaryDetailOutput `json:"our_signatory" gorm:"-"`
	CurrencyExternal      *DictionaryDetailOutput `json:"currency" gorm:"-"`
	TypeExternal          *DictionaryDetailOutput `json:"type" gorm:"-"`
	//其他属性
	SigningDate       *string `json:"signing_date"`
	EffectiveDate     *string `json:"effective_date"`
	CommissioningDate *string `json:"commissioning_date"`
	CompletionDate    *string `json:"completion_date"`

	Amount             *float64 `json:"amount"`
	ExchangeRate       *float64 `json:"exchange_rate"`
	ConstructionPeriod *int     `json:"construction_period"`

	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Content     *string `json:"content"`
	Deliverable *string `json:"deliverable"`
	PenaltyRule *string `json:"penalty_rule"`
	Attachment  *string `json:"attachment"`
	Operator    *string `json:"operator"`
}

func (c *ContractGet) Get() response.Common {
	var result ContractOutput
	err := global.DB.Model(model.Contract{}).
		Where("id = ?", c.ID).First(&result).Error
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
		//查部门信息
		if result.OrganizationID != nil {
			var record OrganizationOutput
			res := global.DB.Model(&model.Organization{}).
				Where("id = ?", *result.OrganizationID).Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.OrganizationExternal = &record
			}
		}
		//查相关方信息
		if result.RelatedPartyID != nil {
			var record RelatedPartyOutput
			res := global.DB.Model(&model.RelatedParty{}).
				Where("id = ?", *result.RelatedPartyID).Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.RelatedPartyExternal = &record
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
		if result.OurSignatory != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.OurSignatory).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.OurSignatoryExternal = &record
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
		if result.SigningDate != nil {
			temp := *result.SigningDate
			*result.SigningDate = temp[:10]
		}
		if result.EffectiveDate != nil {
			temp := *result.EffectiveDate
			*result.EffectiveDate = temp[:10]
		}
		if result.CommissioningDate != nil {
			temp := *result.CommissioningDate
			*result.CommissioningDate = temp[:10]
		}
		if result.CompletionDate != nil {
			temp := *result.CompletionDate
			*result.CompletionDate = temp[:10]
		}
	}

	return response.SuccessWithData(result)
}

func (c *ContractCreate) Create() response.Common {
	var paramOut model.Contract

	if c.Creator > 0 {
		paramOut.Creator = &c.Creator
	}
	if c.LastModifier > 0 {
		paramOut.LastModifier = &c.LastModifier
	}

	//连接关联表的id
	{
		if c.ProjectID > 0 {
			paramOut.ProjectID = &c.ProjectID
		}
		if c.OrganizationID > 0 {
			paramOut.OrganizationID = &c.OrganizationID
		}
		if c.RelatedPartyID > 0 {
			paramOut.RelatedPartyID = &c.RelatedPartyID
		}
	}

	//连接dictionary_item表的id
	{
		if c.FundDirection > 0 {
			paramOut.FundDirection = &c.FundDirection
		}
		if c.OurSignatory > 0 {
			paramOut.OurSignatory = &c.OurSignatory
		}
		if c.Currency > 0 {
			paramOut.Currency = &c.Currency
		}
		if c.Type > 0 {
			paramOut.Type = &c.Type
		}
	}

	//日期
	{
		if c.SigningDate != "" {
			signingDate, err := time.Parse("2006-01-02", c.SigningDate)
			if err != nil {
				return response.Failure(util.ErrorInvalidDateFormat)
			}
			paramOut.SigningDate = &signingDate
		}

		if c.EffectiveDate != "" {
			effectiveDate, err := time.Parse("2006-01-02", c.EffectiveDate)
			if err != nil {
				return response.Failure(util.ErrorInvalidDateFormat)
			}
			paramOut.EffectiveDate = &effectiveDate
		}

		if c.CommissioningDate != "" {
			commissioningDate, err := time.Parse("2006-01-02", c.CommissioningDate)
			if err != nil {
				return response.Failure(util.ErrorInvalidDateFormat)
			}
			paramOut.CommissioningDate = &commissioningDate
		}

		if c.CompletionDate != "" {
			completionDate, err := time.Parse("2006-01-02", c.CompletionDate)
			if err != nil {
				return response.Failure(util.ErrorInvalidDateFormat)
			}
			paramOut.CompletionDate = &completionDate
		}
	}

	//允许为0的数字
	{
		if c.Amount != nil {
			paramOut.Amount = c.Amount
		}
		if c.ExchangeRate != nil {
			paramOut.ExchangeRate = c.ExchangeRate
		}
		if c.ConstructionPeriod != nil {
			paramOut.ConstructionPeriod = c.ConstructionPeriod
		}
	}

	//允许为null的字符串
	{
		if c.Name != "" {
			paramOut.Name = &c.Name
		}

		if c.Code != "" {
			paramOut.Code = &c.Code
		}

		if c.Content != "" {
			paramOut.Content = &c.Content
		}

		if c.Deliverable != "" {
			paramOut.Deliverable = &c.Deliverable
		}

		if c.PenaltyRule != "" {
			paramOut.PenaltyRule = &c.PenaltyRule
		}

		if c.Attachment != "" {
			paramOut.Attachment = &c.Attachment
		}

		if c.Operator != "" {
			paramOut.Operator = &c.Operator
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "CreateAt", "UpdatedAt", "ID")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (c *ContractUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if c.LastModifier > 0 {
		paramOut["last_modifier"] = c.LastModifier
	}

	//连接关联表的id
	{
		if c.ProjectID != nil {
			if *c.ProjectID > 0 {
				paramOut["project_id"] = *c.ProjectID
			}
		}
		if c.OrganizationID != nil {
			if *c.OrganizationID > 0 {
				paramOut["organization_id"] = c.OrganizationID
			} else if *c.OrganizationID == -1 {
				paramOut["organization_id"] = nil
			}
		}
		if c.RelatedPartyID != nil {
			if *c.RelatedPartyID > 0 {
				paramOut["related_party_id"] = c.RelatedPartyID
			} else if *c.RelatedPartyID == -1 {
				paramOut["related_party_id"] = nil
			}
		}
	}

	//连接dictionary_item表的id
	{
		if c.FundDirection != nil {
			if *c.FundDirection > 0 {
				paramOut["fund_direction"] = c.FundDirection
			} else if *c.FundDirection == -1 {
				paramOut["fund_direction"] = nil
			}
		}
		if c.OurSignatory != nil {
			if *c.OurSignatory > 0 {
				paramOut["our_signatory"] = c.OurSignatory
			} else if *c.OurSignatory == -1 {
				paramOut["our_signatory"] = nil
			}
		}
		if c.Currency != nil {
			if *c.Currency > 0 {
				paramOut["currency"] = c.Currency
			} else if *c.Currency == -1 {
				paramOut["currency"] = nil
			}
		}
		if c.Type != nil {
			if *c.Type > 0 {
				paramOut["type"] = c.Type
			} else if *c.Type == -1 {
				paramOut["type"] = nil
			}
		}
	}

	//日期
	{
		if c.SigningDate != nil {
			if *c.SigningDate != "" {
				var err error
				paramOut["signing_date"], err = time.Parse("2006-01-02", *c.SigningDate)
				if err != nil {
					return response.Failure(util.ErrorInvalidJSONParameters)
				}
			} else {
				paramOut["signing_date"] = nil
			}
		}
		if c.EffectiveDate != nil {
			if *c.EffectiveDate != "" {
				var err error
				paramOut["effective_date"], err = time.Parse("2006-01-02", *c.EffectiveDate)
				if err != nil {
					return response.Failure(util.ErrorInvalidJSONParameters)
				}
			} else {
				paramOut["effective_date"] = nil
			}
		}
		if c.CommissioningDate != nil {
			if *c.CommissioningDate != "" {
				var err error
				paramOut["commissioning_date"], err = time.Parse("2006-01-02", *c.CommissioningDate)
				if err != nil {
					return response.Failure(util.ErrorInvalidJSONParameters)
				}
			} else {
				paramOut["commissioning_date"] = nil
			}
		}
		if c.CompletionDate != nil {
			if *c.CompletionDate != "" {
				var err error
				paramOut["completion_date"], err = time.Parse("2006-01-02", *c.CompletionDate)
				if err != nil {
					return response.Failure(util.ErrorInvalidJSONParameters)
				}
			} else {
				paramOut["completion_date"] = nil
			}
		}
	}

	//允许为0的数字
	{
		if c.Amount != nil {
			if *c.Amount != -1 {
				paramOut["amount"] = c.Amount
			} else {
				paramOut["amount"] = nil
			}
		}
		if c.ExchangeRate != nil {
			if *c.ExchangeRate != -1 {
				paramOut["exchange_rate"] = c.ExchangeRate
			} else {
				paramOut["exchange_rate"] = nil
			}
		}
		if c.ConstructionPeriod != nil {
			if *c.ConstructionPeriod != -1 {
				paramOut["construction_period"] = c.ConstructionPeriod
			} else {
				paramOut["construction_period"] = nil
			}
		}
	}

	//允许为null的字符串
	{
		if c.Name != nil {
			if *c.Name != "" {
				paramOut["name"] = c.Name
			} else {
				paramOut["name"] = nil
			}
		}
		if c.Code != nil {
			if *c.Code != "" {
				paramOut["code"] = c.Code
			} else {
				paramOut["code"] = nil
			}
		}
		if c.Content != nil {
			if *c.Content != "" {
				paramOut["content"] = c.Content
			} else {
				paramOut["content"] = nil
			}
		}
		if c.Deliverable != nil {
			if *c.Deliverable != "" {
				paramOut["deliverable"] = c.Deliverable
			} else {
				paramOut["deliverable"] = nil
			}
		}
		if c.PenaltyRule != nil {
			if *c.PenaltyRule != "" {
				paramOut["penalty_rule"] = c.PenaltyRule
			} else {
				paramOut["penalty_rule"] = nil
			}
		}
		if c.Attachment != nil {
			if *c.Attachment != "" {
				paramOut["attachment"] = c.Attachment
			} else {
				paramOut["attachment"] = nil
			}
		}
		if c.Operator != nil {
			if *c.Operator != "" {
				paramOut["operator"] = c.Operator
			} else {
				paramOut["operator"] = nil
			}
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Contract{}).Where("id = ?", c.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (c *ContractDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.Contract
	global.DB.Where("id = ?", c.ID).Find(&record)
	err := global.DB.Where("id = ?", c.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (c *ContractGetList) GetList() response.List {
	db := global.DB.Model(&model.Contract{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if c.ProjectID > 0 {
		db = db.Where("project_id = ?", c.ProjectID)
	}

	if c.RelatedPartyID > 0 {
		db = db.Where("related_party_id = ?", c.RelatedPartyID)
	}

	if c.FundDirection > 0 {
		db = db.Where("fund_direction = ?", c.FundDirection)
	}

	if c.NameInclude != "" {
		db = db.Where("name like ?", "%"+c.NameInclude+"%")
	}

	//用来确定数据范围
	organizationIDsForDataScope := util.GetOrganizationIDsInDataScope(c.UserID)
	//先找出项目的数据范围
	var projectIDs []int64
	global.DB.Model(&model.Project{}).Where("organization_id in ?", organizationIDsForDataScope).
		Select("id").Find(&projectIDs)
	//然后再加上组织的数据范围
	db = db.Where("organization_id in ? or project_id in ?",
		organizationIDsForDataScope, projectIDs)

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
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Contract{}, orderBy)
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

	//data
	var data []ContractOutput
	db.Model(&model.Contract{}).Debug().Find(&data)

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
			//查部门信息
			if data[i].OrganizationID != nil {
				var record OrganizationOutput
				res := global.DB.Model(&model.Organization{}).
					Where("id = ?", *data[i].OrganizationID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].OrganizationExternal = &record
				}
			}
			//查相关方信息
			if data[i].RelatedPartyID != nil {
				var record RelatedPartyOutput
				res := global.DB.Model(&model.RelatedParty{}).
					Where("id = ?", *data[i].RelatedPartyID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].RelatedPartyExternal = &record
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
			if data[i].OurSignatory != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *data[i].OurSignatory).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].OurSignatoryExternal = &record
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
		}

		//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
		//需要取年月日(即前9位)
		{
			if data[i].SigningDate != nil {
				temp := *data[i].SigningDate
				*data[i].SigningDate = temp[:10]
			}
			if data[i].EffectiveDate != nil {
				temp := *data[i].EffectiveDate
				*data[i].EffectiveDate = temp[:10]
			}
			if data[i].CommissioningDate != nil {
				temp := *data[i].CommissioningDate
				*data[i].CommissioningDate = temp[:10]
			}
			if data[i].CompletionDate != nil {
				temp := *data[i].CompletionDate
				*data[i].CompletionDate = temp[:10]
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
