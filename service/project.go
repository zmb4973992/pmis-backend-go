package service

import (
	"github.com/gin-gonic/gin"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
	"sync"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type ProjectGet struct {
	Id     int64
	UserId int64
}

type ProjectCreate struct {
	UserId       int64
	LastModifier int64
	//连接其他表的id
	OrganizationId int64 `json:"organization_id,omitempty"`
	RelatedPartyId int64 `json:"related_party_id,omitempty"`
	//连接dictionary_item表的id
	Country      int64 `json:"country,omitempty"`
	Type         int64 `json:"type,omitempty"`
	DetailedType int64 `json:"detailed_type,omitempty"` //细分的项目类型
	Currency     int64 `json:"currency,omitempty"`
	Status       int64 `json:"status,omitempty"`
	OurSignatory int64 `json:"our_signatory,omitempty"`
	//日期
	ApprovalDate      string `json:"approval_date,omitempty"`
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
	UserId int64
	Id     int64
	//连接其他表的id
	OrganizationId *int64 `json:"organization_id"`
	RelatedPartyId *int64 `json:"related_party_id"`
	//连接dictionary_item表的id
	Country      *int64 `json:"country"`
	Type         *int64 `json:"type"`
	DetailedType *int64 `json:"detailed_type"`
	Currency     *int64 `json:"currency"`
	Status       *int64 `json:"status"`
	OurSignatory *int64 `json:"our_signatory"`
	//日期
	ApprovalDate      *string `json:"approval_date"`
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

	IgnoreDataAuthority bool `json:"-"`
}

type ProjectDelete struct {
	Id int64
}

type ProjectGetList struct {
	list.Input
	UserId                  int64   `json:"-"`
	NameInclude             string  `json:"name_include,omitempty"`
	RelatedPartyId          int64   `json:"related_party_id,omitempty"`
	OrganizationNameInclude string  `json:"organization_name_include,omitempty"`
	OrganizationIdIn        []int64 `json:"organization_id_in"`
	Country                 int64   `json:"country,omitempty"`
	ApprovalDateGte         string  `json:"approval_date_gte,omitempty"`
	ApprovalDateLte         string  `json:"approval_date_lte,omitempty"`
	//是否忽略数据权限的限制，用于请求数据范围外的全部数据
	IgnoreDataAuthority bool `json:"ignore_data_authority"`
}

type ProjectGetCount struct {
	UserId          int64  `json:"-"`
	ApprovalDateGte string `json:"approval_date_gte,omitempty"`
	ApprovalDateLte string `json:"approval_date_lte,omitempty"`
}

//获取简化的列表（），用于下拉框选项

type ProjectGetSimplifiedList struct {
	UserId int64 `json:"-"`
	//是否忽略数据权限的限制，用于请求数据范围外的全部数据
	IgnoreDataAuthority bool `json:"ignore_data_authority"`
}

//以下为出参

type ProjectOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	OrganizationId *int64 `json:"-"`
	RelatedPartyId *int64 `json:"-"`
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
	ApprovalDate      *string `json:"approval_date"`
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

type ProjectSimplifiedOutput struct {
	Id   int64   `json:"id"`
	Name *string `json:"name"`
}

type projectCheckAuth struct {
	ProjectId int64
	UserId    int64
}

func (p *ProjectGet) Get() (output *ProjectOutput, errCode int) {
	err := global.DB.Model(model.Project{}).
		Where("id = ?", p.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	var authorize projectCheckAuth
	authorize.ProjectId = p.Id
	authorize.UserId = p.UserId
	authorized := authorize.checkAuth()

	if !authorized {
		return nil, util.ErrorUnauthorized
	}

	//默认格式为这样的string：2019-11-01T00:00:00Z，需要取年月日(前9位)
	if output.ApprovalDate != nil {
		temp := *output.ApprovalDate
		*output.ApprovalDate = temp[:10]
	}

	//默认格式为这样的string：2019-11-01T00:00:00Z，需要取年月日(前9位)
	if output.SigningDate != nil {
		temp := *output.SigningDate
		*output.SigningDate = temp[:10]
	}

	//默认格式为这样的string：2019-11-01T00:00:00Z，需要取年月日(前9位)
	if output.EffectiveDate != nil {
		temp := *output.EffectiveDate
		*output.EffectiveDate = temp[:10]
	}

	//查部门信息
	if output.OrganizationId != nil {
		var record OrganizationOutput
		res := global.DB.Model(&model.Organization{}).
			Where("id = ?", *output.OrganizationId).
			Limit(1).
			Find(&record)
		if res.RowsAffected > 0 {
			output.OrganizationExternal = &record
		}
	}

	//查相关方信息
	if output.RelatedPartyId != nil {
		var record RelatedPartyOutput
		res := global.DB.Model(&model.RelatedParty{}).
			Where("id = ?", *output.RelatedPartyId).
			Limit(1).
			Find(&record)
		if res.RowsAffected > 0 {
			output.RelatedPartyExternal = &record
		}
	}

	//查dictionary_item表
	{
		if output.Country != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.Country).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.CountryExternal = &record
			}
		}

		if output.Type != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.Type).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.TypeExternal = &record
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

		if output.Status != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.Status).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.StatusExternal = &record
			}
		}

		if output.OurSignatory != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.OurSignatory).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.OurSignatoryExternal = &record
			}
		}
	}

	return output, util.Success
}

func (p *ProjectCreate) Create() (errCode int) {
	var paramOut model.Project
	if p.UserId > 0 {
		paramOut.Creator = &p.UserId
	}
	if p.LastModifier > 0 {
		paramOut.LastModifier = &p.LastModifier
	}

	//连接其他表的id
	{
		if p.OrganizationId != 0 {
			paramOut.OrganizationId = &p.OrganizationId
		}
		if p.RelatedPartyId != 0 {
			paramOut.RelatedPartyId = &p.RelatedPartyId
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
		if p.ApprovalDate != "" {
			approvalDate, err := time.Parse("2006-01-02", p.ApprovalDate)
			if err != nil {
				return util.ErrorInvalidDateFormat
			}
			paramOut.ApprovalDate = &approvalDate
		}
		if p.SigningDate != "" {
			signingDate, err := time.Parse("2006-01-02", p.SigningDate)
			if err != nil {
				return util.ErrorInvalidDateFormat
			}
			paramOut.SigningDate = &signingDate
		}
		if p.EffectiveDate != "" {
			effectiveDate, err := time.Parse("2006-01-02", p.EffectiveDate)
			if err != nil {
				return util.ErrorInvalidDateFormat
			}
			paramOut.EffectiveDate = &effectiveDate
		}
		if p.CommissioningDate != "" {
			commissioningDate, err := time.Parse("2006-01-02", p.CommissioningDate)
			if err != nil {
				return util.ErrorInvalidDateFormat
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

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}

	return util.Success
}

func (p *ProjectUpdate) Update() (errCode int) {
	var result ProjectOutput
	err := global.DB.Model(model.Project{}).
		Where("id = ?", p.Id).
		First(&result).Error
	if err != nil {
		return util.ErrorRecordNotFound
	}

	if p.IgnoreDataAuthority == false {
		var authorization projectCheckAuth
		authorization.ProjectId = p.Id
		authorization.UserId = p.UserId
		authorized := authorization.checkAuth()
		if !authorized {
			return util.ErrorUnauthorized
		}
	}

	paramOut := make(map[string]any)
	if p.UserId > 0 {
		paramOut["last_modifier"] = p.UserId
	}
	//连接其他表的id
	{
		if p.OrganizationId != nil {
			if *p.OrganizationId > 0 {
				paramOut["organization_id"] = p.OrganizationId
			} else if *p.OrganizationId == -1 {
				paramOut["organization_id"] = nil
			}
		}
		if p.RelatedPartyId != nil {
			if *p.RelatedPartyId > 0 {
				paramOut["related_party_id"] = p.RelatedPartyId
			} else if *p.RelatedPartyId == -1 {
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
		if p.ApprovalDate != nil {
			if *p.ApprovalDate != "" {
				paramOut["approval_date"], err = time.Parse("2006-01-02", *p.ApprovalDate)
				if err != nil {
					return util.ErrorInvalidJSONParameters
				}
			} else {
				paramOut["approval_date"] = nil
			}
		}
		if p.SigningDate != nil {
			if *p.SigningDate != "" {
				paramOut["signing_date"], err = time.Parse("2006-01-02", *p.SigningDate)
				if err != nil {
					return util.ErrorInvalidJSONParameters
				}
			} else {
				paramOut["signing_date"] = nil
			}
		}
		if p.EffectiveDate != nil {
			if *p.EffectiveDate != "" {
				paramOut["effective_date"], err = time.Parse("2006-01-02", *p.EffectiveDate)
				if err != nil {
					return util.ErrorInvalidJSONParameters
				}
			} else {
				paramOut["effective_date"] = nil
			}
		}
		if p.CommissioningDate != nil {
			if *p.CommissioningDate != "" {
				paramOut["commissioning_date"], err = time.Parse("2006-01-02", *p.CommissioningDate)
				if err != nil {
					return util.ErrorInvalidJSONParameters
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

	err = global.DB.Model(&model.Project{}).
		Where("id = ?", p.Id).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}

func (p *ProjectDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.Project
	err := global.DB.Where("id = ?", p.Id).
		Find(&record).
		Delete(&record).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}
	return util.Success
}

func (p *ProjectGetList) GetList() (
	outputs []ProjectOutput, errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.Project{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where
	if p.NameInclude != "" {
		db = db.Where("name like ?", "%"+p.NameInclude+"%")
	}

	if p.OrganizationNameInclude != "" {
		var organizationIds []int64
		global.DB.Model(&model.Organization{}).
			Where("name like ?", "%"+p.OrganizationNameInclude+"%").
			Select("id").
			Find(&organizationIds)
		if len(organizationIds) > 0 {
			db = db.Where("organization_id in ?", organizationIds)
		}
	}

	//将临时表的字段名改为temp，是为了防止临时表的字段和主表发生重复，影响后面的查询
	if p.RelatedPartyId > 0 {
		db = db.Joins("join (select distinct project_id as temp_project_id from contract where contract.related_party_id = ?) as temp on project.id = temp.temp_project_id ", p.RelatedPartyId)
	}

	if len(p.OrganizationIdIn) > 0 {
		db = db.Where("organization_id in ?", p.OrganizationIdIn)
	}

	if p.Country > 0 {
		db = db.Where("country = ?", p.Country)
	}

	if p.ApprovalDateGte != "" {
		db = db.Where("approval_date >= ?", p.ApprovalDateGte)
	}

	if p.ApprovalDateLte != "" {
		db = db.Where("approval_date <= ?", p.ApprovalDateLte)
	}

	//用来确定数据范围
	if p.IgnoreDataAuthority == false {
		organizationIds := util.GetOrganizationIdsForDataAuthority(p.UserId)
		db = db.Where("organization_id in ?", organizationIds)
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
			return nil, util.ErrorSortingFieldDoesNotExist, nil
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
	pageSize := global.Config.Paging.DefaultPageSize
	if p.PagingInput.PageSize != nil && *p.PagingInput.PageSize >= 0 &&
		*p.PagingInput.PageSize <= global.Config.Paging.MaxPageSize {
		pageSize = *p.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.Project{}).Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	var wg sync.WaitGroup

	for j := range outputs {
		i := j
		wg.Add(1)
		go func() {
			defer wg.Done()
			//查询关联表的详情
			{
				//查部门信息
				if outputs[i].OrganizationId != nil {
					var record OrganizationOutput
					res := global.DB.Model(&model.Organization{}).
						Where("id = ?", *outputs[i].OrganizationId).Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].OrganizationExternal = &record
					}
				}
				//查相关方信息
				if outputs[i].RelatedPartyId != nil {
					var record RelatedPartyOutput
					res := global.DB.Model(&model.RelatedParty{}).
						Where("id = ?", *outputs[i].RelatedPartyId).Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].RelatedPartyExternal = &record
					}
				}
			}

			//查dictionary_item表的详情
			{
				if outputs[i].Country != nil {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *outputs[i].Country).Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].CountryExternal = &record
					}
				}
				if outputs[i].Type != nil {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *outputs[i].Type).Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].TypeExternal = &record
					}
				}
				if outputs[i].DetailedType != nil {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *outputs[i].DetailedType).Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].DetailedTypeExternal = &record
					}
				}
				if outputs[i].Currency != nil {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *outputs[i].Currency).Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].CurrencyExternal = &record
					}
				}
				if outputs[i].Status != nil {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *outputs[i].Status).Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].StatusExternal = &record
					}
				}
				if outputs[i].OurSignatory != nil {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *outputs[i].OurSignatory).Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].OurSignatoryExternal = &record
					}
				}
			}
			//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
			//需要取年月日(即前9位)
			{
				if outputs[i].ApprovalDate != nil {
					temp := *outputs[i].ApprovalDate
					*outputs[i].ApprovalDate = temp[:10]
				}
				if outputs[i].SigningDate != nil {
					temp := *outputs[i].SigningDate
					*outputs[i].SigningDate = temp[:10]
				}
				if outputs[i].EffectiveDate != nil {
					temp := *outputs[i].EffectiveDate
					*outputs[i].EffectiveDate = temp[:10]
				}
				if outputs[i].CommissioningDate != nil {
					temp := *outputs[i].CommissioningDate
					*outputs[i].CommissioningDate = temp[:10]
				}
			}

			if p.IgnoreDataAuthority == true {
				var authorize projectCheckAuth
				authorize.ProjectId = outputs[i].Id
				authorize.UserId = p.UserId
				outputs[i].Authorized = authorize.checkAuth()
			} else {
				outputs[i].Authorized = true
			}
		}()
	}

	wg.Wait()

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

func (p *ProjectGetSimplifiedList) GetSimplifiedList() (
	outputs []ProjectSimplifiedOutput, errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.Project{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where

	//用来确定数据范围
	if p.IgnoreDataAuthority == false {
		organizationIds := util.GetOrganizationIdsForDataAuthority(p.UserId)
		db = db.Where("organization_id in ?", organizationIds)
	}

	//count
	var count int64
	db.Count(&count)

	//outputs
	db.Model(&model.Project{}).
		Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	numberOfRecords := int(count)

	return outputs,
		util.Success,
		&list.PagingOutput{
			Page:            1,
			PageSize:        0,
			NumberOfPages:   1,
			NumberOfRecords: numberOfRecords,
		}
}

func (p *ProjectGetCount) GetCount() (output any, errCode int) {
	db := global.DB.Model(&model.Project{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if p.ApprovalDateGte != "" {
		db = db.Where("approval_date >= ?", p.ApprovalDateGte)
	}

	if p.ApprovalDateLte != "" {
		db = db.Where("approval_date <= ?", p.ApprovalDateLte)
	}

	//count
	var count int64
	db.Count(&count)

	return gin.H{"count": int(count)}, util.Success
}

func (p *projectCheckAuth) checkAuth() (authorized bool) {
	//用来确定数据范围内的组织id
	organizationIds := util.GetOrganizationIdsForDataAuthority(p.UserId)
	if len(organizationIds) == 0 {
		return false
	}

	//看看在数据范围内是否有该记录
	var count int64
	global.DB.Model(model.Project{}).
		Where("organization_id in ?", organizationIds).
		Where("id = ?", p.ProjectId).
		Count(&count)

	if count > 0 {
		return true
	}

	return false
}
