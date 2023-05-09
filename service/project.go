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

type ProjectGet struct {
	SnowID int64
}

type ProjectCreate struct {
	Creator      int64
	LastModifier int64
	//连接其他表的id
	OrganizationSnowID int64 `json:"organization_snow_id,omitempty"`
	RelatedPartySnowID int64 `json:"related_party_snow_id,omitempty"`
	//连接dictionary_item表的id
	Country      int64 `json:"country,omitempty"`
	Type         int64 `json:"type,omitempty"`
	DetailedType int64 `json:"detailed_type,omitempty"` //细分的项目类型
	Currency     int64 `json:"currency,omitempty"`
	Status       int64 `json:"status,omitempty"`
	OurSignatory int64 `json:"our_signatory,omitempty"`
	//日期
	SigningDate       string `json:"signing_date,omitempty"`
	EffectiveDate     string `json:"effective_date,omitempty"`
	CommissioningDate string `json:"commissioning_date,omitempty"`
	//数字(允许为0)
	Amount             *float64 `json:"amount"`
	ExchangeRate       *float64 `json:"exchange_rate"`
	ConstructionPeriod *int     `json:"construction_period"`
	//字符串(允许为null)
	Code    string `json:"code,omitempty"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ProjectUpdate struct {
	LastModifier int64
	SnowID       int64
	//连接其他表的id
	OrganizationSnowID *int64 `json:"organization_snow_id"`
	RelatedPartySnowID *int64 `json:"related_party_snow_id"`
	//连接dictionary_item表的id
	Country      *int64 `json:"country"`
	Type         *int64 `json:"type"`
	DetailedType *int64 `json:"detailed_type"`
	Currency     *int64 `json:"currency"`
	Status       *int64 `json:"status"`
	OurSignatory *int64 `json:"our_signatory"`
	//日期
	SigningDate       *string `json:"signing_date"`
	EffectiveDate     *string `json:"effective_date"`
	CommissioningDate *string `json:"commissioning_date"`
	//数字(允许为0)
	Amount             *float64 `json:"amount"`
	ExchangeRate       *float64 `json:"exchange_rate"`
	ConstructionPeriod *int     `json:"construction_period"`
	//字符串(允许为null)
	Code    *string `json:"code"`
	Name    *string `json:"name"`
	Content *string `json:"content"`
}

type ProjectDelete struct {
	SnowID int64
}

type ProjectGetList struct {
	list.Input
	list.DataScopeInput
	NameInclude             string  `json:"name_include,omitempty"`
	OrganizationNameInclude string  `json:"organization_name_include,omitempty"`
	OrganizationSnowIDIn    []int64 `json:"organization_snow_id_in"`
}

//以下为出参

type ProjectOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	SnowID       int64  `json:"snow_id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	OrganizationSnowID *int64 `json:"-"`
	RelatedPartySnowID *int64 `json:"-"`
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示
	Country      *int64 `json:"-"`
	Type         *int64 `json:"-"`
	DetailedType *int64 `json:"-"`
	Currency     *int64 `json:"-"`
	Status       *int64 `json:"-"`
	OurSignatory *int64 `json:"-"`
	//关联表的详情，不需要gorm查询，需要在json中显示
	OrganizationExternal *OrganizationOutput `json:"organization" gorm:"-"`
	RelatedPartyExternal *RelatedPartyOutput `json:"related_party_external" gorm:"-"`
	//dictionary_item表的详情，不需要gorm查询，需要在json中显示
	CountryExternal      *DictionaryDetailOutput `json:"country" gorm:"-"`
	TypeExternal         *DictionaryDetailOutput `json:"type" gorm:"-"`
	DetailedTypeExternal *DictionaryDetailOutput `json:"detailed_type" gorm:"-"`
	CurrencyExternal     *DictionaryDetailOutput `json:"currency" gorm:"-"`
	StatusExternal       *DictionaryDetailOutput `json:"status" gorm:"-"`
	OurSignatoryExternal *DictionaryDetailOutput `json:"our_signatory" gorm:"-"`
	//其他属性
	SigningDate       *string `json:"signing_date"`
	EffectiveDate     *string `json:"effective_date"`
	CommissioningDate *string `json:"commissioning_date"`

	Amount             *float64 `json:"amount"`
	ExchangeRate       *float64 `json:"exchange_rate"`
	ConstructionPeriod *int     `json:"construction_period"`

	Code    *string `json:"code"`
	Name    *string `json:"name"`
	Content *string `json:"content"`
}

func (p *ProjectGet) Get() response.Common {
	var result ProjectOutput
	err := global.DB.Model(model.Project{}).
		Where("snow_id = ?", p.SnowID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	//默认格式为这样的string：2019-11-01T00:00:00Z，需要取年月日(前9位)
	if result.SigningDate != nil {
		temp := *result.SigningDate
		*result.SigningDate = temp[:10]
	}

	//默认格式为这样的string：2019-11-01T00:00:00Z，需要取年月日(前9位)
	if result.EffectiveDate != nil {
		temp := *result.EffectiveDate
		*result.EffectiveDate = temp[:10]
	}

	//查部门信息
	if result.OrganizationSnowID != nil {
		var record OrganizationOutput
		res := global.DB.Model(&model.Organization{}).
			Where("snow_id = ?", *result.OrganizationSnowID).Limit(1).Find(&record)
		if res.RowsAffected > 0 {
			result.OrganizationExternal = &record
		}
	}

	//查dictionary_item表
	{
		if result.Country != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("snow_id = ?", *result.Country).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.CountryExternal = &record
			}
		}

		if result.Type != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("snow_id = ?", *result.Type).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.TypeExternal = &record
			}
		}

		if result.Currency != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("snow_id = ?", *result.Currency).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.CurrencyExternal = &record
			}
		}

		if result.Status != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("snow_id = ?", *result.Status).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.StatusExternal = &record
			}
		}

		if result.OurSignatory != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("snow_id = ?", *result.OurSignatory).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.OurSignatoryExternal = &record
			}
		}
	}

	return response.SuccessWithData(result)
}

func (p *ProjectCreate) Create() response.Common {
	var paramOut model.Project
	if p.Creator > 0 {
		paramOut.Creator = &p.Creator
	}
	if p.LastModifier > 0 {
		paramOut.LastModifier = &p.LastModifier
	}
	paramOut.SnowID = idgen.NextId()

	//连接其他表的id
	{
		if p.OrganizationSnowID != 0 {
			paramOut.OrganizationSnowID = &p.OrganizationSnowID
		}
		if p.RelatedPartySnowID != 0 {
			paramOut.RelatedPartySnowID = &p.RelatedPartySnowID
		}
	}
	//连接dictionary_item表的id
	{
		if p.Country > 0 {
			paramOut.Country = &p.Country
		}
		if p.Type > 0 {
			paramOut.Type = &p.Type
		}
		if p.DetailedType > 0 {
			paramOut.DetailedType = &p.DetailedType
		}
		if p.Currency > 0 {
			paramOut.Currency = &p.Currency
		}
		if p.Status > 0 {
			paramOut.Status = &p.Status
		}
		if p.OurSignatory > 0 {
			paramOut.OurSignatory = &p.OurSignatory
		}
	}
	//日期
	{
		if p.SigningDate != "" {
			signingDate, err := time.Parse("2006-01-02", p.SigningDate)
			if err != nil {
				return response.Failure(util.ErrorInvalidDateFormat)
			}
			paramOut.SigningDate = &signingDate
		}
		if p.EffectiveDate != "" {
			effectiveDate, err := time.Parse("2006-01-02", p.EffectiveDate)
			if err != nil {
				return response.Failure(util.ErrorInvalidDateFormat)
			}
			paramOut.EffectiveDate = &effectiveDate
		}
		if p.CommissioningDate != "" {
			commissioningDate, err := time.Parse("2006-01-02", p.CommissioningDate)
			if err != nil {
				return response.Failure(util.ErrorInvalidDateFormat)
			}
			paramOut.CommissioningDate = &commissioningDate
		}
	}
	//数字(允许为0)
	{
		if p.Amount != nil {
			paramOut.Amount = p.Amount
		}
		if p.ExchangeRate != nil {
			paramOut.ExchangeRate = p.ExchangeRate
		}
		if p.ConstructionPeriod != nil {
			paramOut.ConstructionPeriod = p.ConstructionPeriod
		}
	}
	//字符串(允许为null)
	{
		if p.Code != "" {
			paramOut.Code = &p.Code
		}
		if p.Name != "" {
			paramOut.Name = &p.Name
		}
		if p.Content != "" {
			paramOut.Content = &p.Content
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "CreateAt", "UpdatedAt", "SnowID")

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

func (p *ProjectUpdate) Update() response.Common {
	paramOut := make(map[string]any)
	if p.LastModifier > 0 {
		paramOut["last_modifier"] = p.LastModifier
	}
	//连接其他表的id
	{
		if p.OrganizationSnowID != nil {
			if *p.OrganizationSnowID > 0 {
				paramOut["organization_snow_id"] = p.OrganizationSnowID
			} else if *p.OrganizationSnowID == -1 {
				paramOut["organization_snow_id"] = nil
			}
		}
		if p.RelatedPartySnowID != nil {
			if *p.RelatedPartySnowID > 0 {
				paramOut["related_party_snow_id"] = p.RelatedPartySnowID
			} else if *p.RelatedPartySnowID == -1 {
				paramOut["related_party_snow_id"] = nil
			}
		}
	}
	//连接dictionary_item表的id
	{
		if p.Country != nil {
			if *p.Country > 0 {
				paramOut["country"] = p.Country
			} else if *p.Country == -1 {
				paramOut["country"] = nil
			}
		}
		if p.Type != nil {
			if *p.Type > 0 {
				paramOut["type"] = p.Type
			} else if *p.Type == -1 {
				paramOut["type"] = nil
			}
		}
		if p.DetailedType != nil {
			if *p.DetailedType > 0 {
				paramOut["detailed_type"] = p.DetailedType
			} else if *p.DetailedType == -1 {
				paramOut["detailed_type"] = nil
			}
		}
		if p.Currency != nil {
			if *p.Currency > 0 {
				paramOut["currency"] = p.Currency
			} else if *p.Currency == -1 {
				paramOut["currency"] = nil
			}
		}
		if p.Status != nil {
			if *p.Status > 0 {
				paramOut["status"] = p.Status
			} else if *p.Status == -1 {
				paramOut["status"] = nil
			}
		}
		if p.OurSignatory != nil {
			if *p.OurSignatory > 0 {
				paramOut["our_signatory"] = p.OurSignatory
			} else if *p.OurSignatory == -1 {
				paramOut["our_signatory"] = nil
			}
		}
	}
	//日期
	{
		if p.SigningDate != nil {
			if *p.SigningDate != "" {
				var err error
				paramOut["signing_date"], err = time.Parse("2006-01-02", *p.SigningDate)
				if err != nil {
					return response.Failure(util.ErrorInvalidJSONParameters)
				}
			} else {
				paramOut["signing_date"] = nil
			}
		}
		if p.EffectiveDate != nil {
			if *p.EffectiveDate != "" {
				var err error
				paramOut["effective_date"], err = time.Parse("2006-01-02", *p.EffectiveDate)
				if err != nil {
					return response.Failure(util.ErrorInvalidJSONParameters)
				}
			} else {
				paramOut["effective_date"] = nil
			}
		}
		if p.CommissioningDate != nil {
			if *p.CommissioningDate != "" {
				var err error
				paramOut["commissioning_date"], err = time.Parse("2006-01-02", *p.CommissioningDate)
				if err != nil {
					return response.Failure(util.ErrorInvalidJSONParameters)
				}
			} else {
				paramOut["commissioning_date"] = nil
			}
		}
	}
	//数字(允许为0)
	{
		if p.Amount != nil {
			if *p.Amount != -1 {
				paramOut["amount"] = p.Amount
			} else {
				paramOut["amount"] = nil
			}
		}
		if p.ExchangeRate != nil {
			if *p.ExchangeRate != -1 {
				paramOut["exchange_rate"] = p.ExchangeRate
			} else {
				paramOut["exchange_rate"] = nil
			}
		}
		if p.ConstructionPeriod != nil {
			if *p.ConstructionPeriod != -1 {
				paramOut["construction_period"] = p.ConstructionPeriod
			} else {
				paramOut["construction_period"] = nil
			}
		}
	}
	//字符串(允许为null)
	{
		if p.Code != nil {
			if *p.Code != "" {
				paramOut["code"] = p.Code
			} else {
				paramOut["code"] = nil
			}
		}
		if p.Name != nil {
			if *p.Name != "" {
				paramOut["name"] = p.Name
			} else {
				paramOut["name"] = nil
			}
		}
		if p.Content != nil {
			if *p.Content != "" {
				paramOut["content"] = p.Content
			} else {
				paramOut["content"] = nil
			}
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "Deleter", "CreateAt", "UpdatedAt", "DeletedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Project{}).Where("snow_id = ?", p.SnowID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (p *ProjectDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.Project
	global.DB.Where("snow_id = ?", p.SnowID).Find(&record)
	err := global.DB.Where("snow_id = ?", p.SnowID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (p *ProjectGetList) GetList() response.List {
	db := global.DB.Model(&model.Project{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if p.NameInclude != "" {
		db = db.Where("name like ?", "%"+p.NameInclude+"%")
	}

	if p.OrganizationNameInclude != "" {
		var organizationSnowIDs []int64
		global.DB.Model(&model.Organization{}).Where("name like ?", "%"+p.OrganizationNameInclude+"%").
			Select("snow_id").Find(&organizationSnowIDs)
		if len(organizationSnowIDs) > 0 {
			db = db.Where("organization_snow_id in ?", organizationSnowIDs)
		}
	}

	if len(p.OrganizationSnowIDIn) > 0 {
		db = db.Where("organization_snow_id in ?", p.OrganizationSnowIDIn)
	}

	//用来确定数据范围
	organizationIDsForDataScope := util.GetOrganizationSnowIDsInDataScope(p.UserSnowID)
	db = db.Where("organization_snow_id in ?", organizationIDsForDataScope)

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
			db = db.Order("snow_id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Project{}, orderBy)
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
	var data []ProjectOutput
	db.Model(&model.Project{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	for i := range data {
		//查询关联表的详情
		{
			//查部门信息
			if data[i].OrganizationSnowID != nil {
				var record OrganizationOutput
				res := global.DB.Model(&model.Organization{}).
					Where("snow_id = ?", *data[i].OrganizationSnowID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].OrganizationExternal = &record
				}
			}
			//查相关方信息
			if data[i].RelatedPartySnowID != nil {
				var record RelatedPartyOutput
				res := global.DB.Model(&model.RelatedParty{}).
					Where("snow_id = ?", *data[i].RelatedPartySnowID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].RelatedPartyExternal = &record
				}
			}
		}

		//查dictionary_item表的详情
		{
			if data[i].Country != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("snow_id = ?", *data[i].Country).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].CountryExternal = &record
				}
			}
			if data[i].Type != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("snow_id = ?", *data[i].Type).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].TypeExternal = &record
				}
			}
			if data[i].DetailedType != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("snow_id = ?", *data[i].DetailedType).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].DetailedTypeExternal = &record
				}
			}
			if data[i].Currency != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("snow_id = ?", *data[i].Currency).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].CurrencyExternal = &record
				}
			}
			if data[i].Status != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("snow_id = ?", *data[i].Status).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].StatusExternal = &record
				}
			}
			if data[i].OurSignatory != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("snow_id = ?", *data[i].OurSignatory).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].OurSignatoryExternal = &record
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
