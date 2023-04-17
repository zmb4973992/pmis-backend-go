package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type ProjectGet struct {
	ID int
}

type ProjectCreate struct {
	Creator      int
	LastModifier int
	//连接其他表的id
	DepartmentID   int `json:"department_id,omitempty"`
	RelatedPartyID int `json:"related_party_id,omitempty"`
	//连接dictionary_item表的id
	Country      int `json:"country,omitempty"`
	Type         int `json:"type,omitempty"`
	DetailedType int `json:"detailed_type,omitempty"` //细分的项目类型
	Currency     int `json:"currency,omitempty"`
	Status       int `json:"status,omitempty"`
	OurSignatory int `json:"our_signatory,omitempty"`
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
	LastModifier int
	ID           int
	//连接其他表的id
	DepartmentID   *int `json:"department_id"`
	RelatedPartyID *int `json:"related_party_id"`
	//连接dictionary_item表的id
	Country      *int `json:"country"`
	Type         *int `json:"type"`
	DetailedType *int `json:"detailed_type"`
	Currency     *int `json:"currency"`
	Status       *int `json:"status"`
	OurSignatory *int `json:"our_signatory"`
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
	ID int
}

type ProjectGetList struct {
	dto.ListInput
	dto.AuthInput
	NameInclude           string `json:"name_include,omitempty"`
	DepartmentNameInclude string `json:"department_name_include,omitempty"`
	DepartmentIDIn        []int  `json:"department_id_in"`
}

//以下为出参

type ProjectOutput struct {
	Creator      *int `json:"creator"`
	LastModifier *int `json:"last_modifier"`
	ID           int  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	DepartmentID   *int `json:"-"`
	RelatedPartyID *int `json:"-"`
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示
	Country      *int `json:"-"`
	Type         *int `json:"-"`
	DetailedType *int `json:"-"`
	Currency     *int `json:"-"`
	Status       *int `json:"-"`
	OurSignatory *int `json:"-"`
	//关联表的详情，不需要gorm查询，需要在json中显示
	DepartmentExternal   *DepartmentOutput   `json:"department" gorm:"-"`
	RelatedPartyExternal *RelatedPartyOutput `json:"related_party_external" gorm:"-"`
	//dictionary_item表的详情，不需要gorm查询，需要在json中显示
	CountryExternal      *DictionaryItemOutput `json:"country" gorm:"-"`
	TypeExternal         *DictionaryItemOutput `json:"type" gorm:"-"`
	DetailedTypeExternal *DictionaryItemOutput `json:"detailed_type" gorm:"-"`
	CurrencyExternal     *DictionaryItemOutput `json:"currency" gorm:"-"`
	StatusExternal       *DictionaryItemOutput `json:"status" gorm:"-"`
	OurSignatoryExternal *DictionaryItemOutput `json:"our_signatory" gorm:"-"`
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
		Where("id = ?", p.ID).First(&result).Error
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
	if result.DepartmentID != nil {
		var record DepartmentOutput
		res := global.DB.Model(&model.Department{}).
			Where("id=?", *result.DepartmentID).Limit(1).Find(&record)
		if res.RowsAffected > 0 {
			result.DepartmentExternal = &record
		}
	}

	//查dictionary_item表
	{
		if result.Country != nil {
			var record DictionaryItemOutput
			res := global.DB.Model(&model.DictionaryItem{}).
				Where("id = ?", *result.Country).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.CountryExternal = &record
			}
		}

		if result.Type != nil {
			var record DictionaryItemOutput
			res := global.DB.Model(&model.DictionaryItem{}).
				Where("id = ?", *result.Type).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.TypeExternal = &record
			}
		}

		if result.Currency != nil {
			var record DictionaryItemOutput
			res := global.DB.Model(&model.DictionaryItem{}).
				Where("id = ?", *result.Currency).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.CurrencyExternal = &record
			}
		}

		if result.Status != nil {
			var record DictionaryItemOutput
			res := global.DB.Model(&model.DictionaryItem{}).
				Where("id = ?", *result.Status).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.StatusExternal = &record
			}
		}

		if result.OurSignatory != nil {
			var record DictionaryItemOutput
			res := global.DB.Model(&model.DictionaryItem{}).
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
		if p.DepartmentID != 0 {
			paramOut.DepartmentID = &p.DepartmentID
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
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "Deleter", "CreateAt", "UpdatedAt", "DeletedAt")

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
		if p.DepartmentID != nil {
			if *p.DepartmentID > 0 {
				paramOut["department_id"] = p.DepartmentID
			} else if *p.DepartmentID == -1 {
				paramOut["department_id"] = nil
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
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "Deleter", "CreateAt", "UpdatedAt", "DeletedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Project{}).Where("id = ?", p.ID).
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

	if p.DepartmentNameInclude != "" {
		var departmentIDs []int
		global.DB.Model(&model.Department{}).Where("name like ?", "%"+p.DepartmentNameInclude+"%").
			Select("id").Find(&departmentIDs)
		if len(departmentIDs) > 0 {
			db = db.Where("department_id in ?", departmentIDs)
		}
	}

	if len(p.DepartmentIDIn) > 0 {
		db = db.Where("department_id in ?", p.DepartmentIDIn)
	}

	if p.IsShowedByRole {
		//先获得最大角色的名称
		biggestRoleName := util.GetBiggestRoleName(p.UserID)
		if biggestRoleName == "事业部级" {
			//获取所在事业部的id数组
			businessDivisionIDs := util.GetBusinessDivisionIDs(p.UserID)
			//获取归属这些事业部的部门id数组
			var departmentIDs []int
			global.DB.Model(&model.Department{}).Where("superior_id in ?", businessDivisionIDs).
				Select("id").Find(&departmentIDs)
			//两个数组进行合并
			departmentIDs = append(departmentIDs, businessDivisionIDs...)
			//找到部门id在上面两个数组中的记录
			db = db.Where("department_id in ?", departmentIDs)
		} else if biggestRoleName == "部门级" || biggestRoleName == "项目级" {
			//获取用户所属部门的id数组
			departmentIDs := util.GetDepartmentIDs(p.UserID)
			//找到部门id在上面数组中的记录
			db = db.Where("department_id in ?", departmentIDs)
		}
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
	db = db.Limit(pageSize)

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
			if data[i].DepartmentID != nil {
				var record DepartmentOutput
				res := global.DB.Model(&model.Department{}).
					Where("id=?", *data[i].DepartmentID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].DepartmentExternal = &record
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
			if data[i].Country != nil {
				var record DictionaryItemOutput
				res := global.DB.Model(&model.DictionaryItem{}).
					Where("id = ?", *data[i].Country).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].CountryExternal = &record
				}
			}
			if data[i].Type != nil {
				var record DictionaryItemOutput
				res := global.DB.Model(&model.DictionaryItem{}).
					Where("id = ?", *data[i].Type).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].TypeExternal = &record
				}
			}
			if data[i].DetailedType != nil {
				var record DictionaryItemOutput
				res := global.DB.Model(&model.DictionaryItem{}).
					Where("id = ?", *data[i].DetailedType).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].DetailedTypeExternal = &record
				}
			}
			if data[i].Currency != nil {
				var record DictionaryItemOutput
				res := global.DB.Model(&model.DictionaryItem{}).
					Where("id = ?", *data[i].Currency).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].CurrencyExternal = &record
				}
			}
			if data[i].Status != nil {
				var record DictionaryItemOutput
				res := global.DB.Model(&model.DictionaryItem{}).
					Where("id = ?", *data[i].Status).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].StatusExternal = &record
				}
			}
			if data[i].OurSignatory != nil {
				var record DictionaryItemOutput
				res := global.DB.Model(&model.DictionaryItem{}).
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
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return response.List{
		Data: data,
		Paging: &dto.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
