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
	ID     int64
	UserID int64
}

type ProjectCreate struct {
	Creator      int64
	LastModifier int64
	//连接其他表的id
	OrganizationID int64 `json:"organization_id,omitempty"`
	RelatedPartyID int64 `json:"related_party_id,omitempty"`
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
	Sort    int    `json:"sort,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ProjectUpdate struct {
	LastModifier int64
	ID           int64
	//连接其他表的id
	OrganizationID *int64 `json:"organization_id"`
	RelatedPartyID *int64 `json:"related_party_id"`
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
	Sort    *int    `json:"sort"`
}

type ProjectDelete struct {
	ID int64
}

type ProjectGetList struct {
	list.Input
	list.DataScopeInput
	NameInclude             string  `json:"name_include,omitempty"`
	RelatedPartyID          int64   `json:"related_party_id,omitempty"`
	OrganizationNameInclude string  `json:"organization_name_include,omitempty"`
	OrganizationIDIn        []int64 `json:"organization_id_in"`

	//是否忽略数据范围的限制，用于请求数据范围外的全部数据
	IgnoreDataScope bool `json:"ignore_data_scope"`
}

//以下为出参

type ProjectOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	OrganizationID *int64 `json:"-"`
	RelatedPartyID *int64 `json:"-"`
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示
	Country      *int64 `json:"-"`
	Type         *int64 `json:"-"`
	DetailedType *int64 `json:"-"`
	Currency     *int64 `json:"-"`
	Status       *int64 `json:"-"`
	OurSignatory *int64 `json:"-"`
	//关联表的详情，不需要gorm查询，需要在json中显示
	OrganizationExternal *OrganizationOutput `json:"organization" gorm:"-"`
	RelatedPartyExternal *RelatedPartyOutput `json:"related_party" gorm:"-"`
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
	Sort    *int    `json:"sort"`

	//用来告诉前端，该记录是否为数据范围内，用来判定是否能访问详情、是否需要做跳转等
	Authorized bool `json:"authorized" gorm:"-"`
}

type projectAuthorize struct {
	ID     int64
	UserID int64
}

func (p *ProjectGet) Get() response.Common {
	var result ProjectOutput
	err := global.DB.Model(model.Project{}).
		Where("id = ?", p.ID).First(&result).Error
	if err != nil {
		return response.Failure(util.ErrorRecordNotFound)
	}

	var authorize projectAuthorize
	authorize.ID = p.ID
	authorize.UserID = p.UserID
	authorizationResult := authorize.authorizedOrNot()

	if !authorizationResult {
		return response.Failure(util.ErrorUnauthorized)
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

	//查dictionary_item表
	{
		if result.Country != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.Country).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.CountryExternal = &record
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

		if result.Currency != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.Currency).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.CurrencyExternal = &record
			}
		}

		if result.Status != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.Status).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.StatusExternal = &record
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

	//连接其他表的id
	{
		if p.OrganizationID != 0 {
			paramOut.OrganizationID = &p.OrganizationID
		}
		if p.RelatedPartyID != 0 {
			paramOut.RelatedPartyID = &p.RelatedPartyID
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
		if p.Sort > 0 {
			paramOut.Sort = &p.Sort
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

	return response.Success()
}

func (p *ProjectUpdate) Update() response.Common {
	var result ProjectOutput
	err := global.DB.Model(model.Project{}).
		Where("id = ?", p.ID).First(&result).Error
	if err != nil {
		return response.Failure(util.ErrorRecordNotFound)
	}

	var authorize projectAuthorize
	authorize.ID = p.ID
	authorize.UserID = p.LastModifier
	authorizationResult := authorize.authorizedOrNot()

	if !authorizationResult {
		return response.Failure(util.ErrorUnauthorized)
	}

	paramOut := make(map[string]any)
	if p.LastModifier > 0 {
		paramOut["last_modifier"] = p.LastModifier
	}
	//连接其他表的id
	{
		if p.OrganizationID != nil {
			if *p.OrganizationID > 0 {
				paramOut["organization_id"] = p.OrganizationID
			} else if *p.OrganizationID == -1 {
				paramOut["organization_id"] = nil
			}
		}
		if p.RelatedPartyID != nil {
			if *p.RelatedPartyID > 0 {
				paramOut["related_party_id"] = p.RelatedPartyID
			} else if *p.RelatedPartyID == -1 {
				paramOut["related_party_id"] = nil
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
		if p.Sort != nil {
			if *p.Sort > 0 {
				paramOut["sort"] = p.Sort
			} else {
				paramOut["sort"] = nil
			}
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err = global.DB.Model(&model.Project{}).Where("id = ?", p.ID).
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
	global.DB.Where("id = ?", p.ID).Find(&record)
	err := global.DB.Where("id = ?", p.ID).Delete(&record).Error

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
		var organizationIDs []int64
		global.DB.Model(&model.Organization{}).Where("name like ?", "%"+p.OrganizationNameInclude+"%").
			Select("id").Find(&organizationIDs)
		if len(organizationIDs) > 0 {
			db = db.Where("organization_id in ?", organizationIDs)
		}
	}

	//将临时表的字段名改为temp，是为了防止临时表的字段和主表发生重复，影响后面的查询
	if p.RelatedPartyID > 0 {
		db = db.Joins("join (select distinct project_id as temp_project_id from contract where contract.related_party_id = ?) as temp on project.id = temp.temp_project_id ", p.RelatedPartyID)
	}

	if len(p.OrganizationIDIn) > 0 {
		db = db.Where("organization_id in ?", p.OrganizationIDIn)
	}

	//用来确定数据范围
	if p.IgnoreDataScope == false {
		organizationIDs := util.GetOrganizationIDs(p.UserID)
		db = db.Where("organization_id in ?", organizationIDs)
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
		exists := util.FieldIsInModel(&model.Project{}, orderBy)
		if !exists {
			return response.FailureForList(util.ErrorSortingFieldDoesNotExist)
		}
		//orderBy = "project." + orderBy
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
				res := global.DB.Debug().Model(&model.RelatedParty{}).
					Where("id = ?", *data[i].RelatedPartyID).Limit(1).Find(&record)
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
					Where("id = ?", *data[i].Country).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].CountryExternal = &record
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
			if data[i].DetailedType != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *data[i].DetailedType).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].DetailedTypeExternal = &record
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
			if data[i].Status != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *data[i].Status).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].StatusExternal = &record
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

		if p.IgnoreDataScope == true {
			var authorize projectAuthorize
			authorize.ID = data[i].ID
			authorize.UserID = p.UserID
			data[i].Authorized = authorize.authorizedOrNot()
		} else {
			data[i].Authorized = true
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

// 该方法一定要在确定记录存在后再调用
func (p *projectAuthorize) authorizedOrNot() bool {
	//用来确定数据范围内的组织id
	organizationIDs := util.GetOrganizationIDs(p.UserID)
	if len(organizationIDs) == 0 {
		return false
	}

	batchID := idgen.NextId()

	var temps []model.Temp
	for i := range organizationIDs {
		var temp model.Temp
		temp.OrganizationID = &organizationIDs[i]
		temp.BatchID = batchID
		temps = append(temps, temp)
	}
	global.DB.CreateInBatches(&temps, 100)

	//找出数据范围内的项目id
	var projectIDs []int64
	global.DB.Model(&model.Project{}).
		Joins("join temp on project.organization_id = temp.organization_id").
		Where("batch_id = ?", batchID).
		Select("project.id").Find(&projectIDs)

	if len(projectIDs) == 0 {
		return false
	}

	for i := range projectIDs {
		var temp model.Temp
		temp.ProjectID = &projectIDs[i]
		temp.BatchID = batchID
		temps = append(temps, temp)
	}
	global.DB.CreateInBatches(&temps, 100)

	//看看在数据范围内是否有该记录
	var count int64
	global.DB.Model(model.Project{}).
		Joins("join (select distinct project_id from temp where batch_id = ?) as temp2 on project.id = temp2.project_id", batchID).
		Where("id = ?", p.ID).
		Count(&count)
	if count > 0 {
		return true
	}

	return false
}
