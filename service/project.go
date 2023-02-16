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
	Creator        int
	LastModifier   int
	Code           string   `json:"code,omitempty"`
	Name           string   `json:"name,omitempty"`
	Country        int      `json:"country,omitempty"`
	Province       int      `json:"province,omitempty"`
	Type           int      `json:"type,omitempty"`
	Amount         *float64 `json:"amount"`
	Currency       int      `json:"currency,omitempty"`
	ExchangeRate   *float64 `json:"exchange_rate"`
	DepartmentID   int      `json:"department_id,omitempty"`
	RelatedPartyID int      `json:"related_party_id,omitempty"`
	SigningDate    string   `json:"signing_date,omitempty"`
	EffectiveDate  string   `json:"effective_date,omitempty"`
	Content        string   `json:"content,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ProjectUpdate struct {
	LastModifier   int
	ID             int
	Code           *string  `json:"code"`
	Name           *string  `json:"name"`
	Country        *int     `json:"country"`
	Province       *int     `json:"province"`
	Type           *int     `json:"type"`
	Amount         *float64 `json:"amount"`
	Currency       *int     `json:"currency"`
	ExchangeRate   *float64 `json:"exchange_rate"`
	DepartmentID   *int     `json:"department_id"`
	RelatedPartyID *int     `json:"related_party_id"`
	SigningDate    *string  `json:"signing_date"`
	EffectiveDate  *string  `json:"effective_date"`
	Content        *string  `json:"content"`
}

type ProjectDelete struct {
	Deleter int
	ID      int
}

type ProjectGetList struct {
	dto.ListInput
	dto.AuthInput
	NameInclude           string `json:"name_include,omitempty"`
	DepartmentNameInclude string `json:"department_name_include,omitempty"`
	DepartmentIDIn        []int  `json:"department_id_in"`
}

type ProjectGetArray struct {
	dto.ListInput
	dto.AuthInput
	NameInclude           string `json:"name_include,omitempty"` //包含项目全称和项目简称
	DepartmentNameInclude string `json:"department_name_include,omitempty"`
	DepartmentIDIn        []int  `json:"department_id_in"`
}

//以下为出参

type ProjectOutput struct {
	Creator      *int `json:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" gorm:"id"`

	Code           *string           `json:"code" gorm:"code"`
	Name           *string           `json:"name" gorm:"name"`
	Country        *int              `json:"country" gorm:"country"`
	Province       *int              `json:"province" gorm:"province"`
	Type           *int              `json:"type" gorm:"type"`
	Amount         *float64          `json:"amount" gorm:"amount"`
	Currency       *int              `json:"currency" gorm:"currency"`
	ExchangeRate   *float64          `json:"exchange_rate" gorm:"exchange_rate"`
	RelatedPartyID *int              `json:"related_party_id" gorm:"related_party_id"`
	DepartmentID   *int              `json:"-" gorm:"department_id"`
	SigningDate    *string           `json:"signing_date" gorm:"signing_date"`
	EffectiveDate  *string           `json:"effective_date" gorm:"effective_date"`
	Content        *string           `json:"content" gorm:"content"`
	Department     *DepartmentOutput `json:"department"`
}

type ProjectGetListOutput struct {
	Creator      *int `json:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" gorm:"id"`

	Code                 *string           `json:"code"`
	Name                 *string           `json:"name"`
	CountryInternal      *int              `json:"-" gorm:"column:country"` //用来接收gorm的值，不展示
	Country              *string           `json:"country" gorm:"-"`        //经过查询后，展示数据
	TypeInternal         *int              `json:"-" gorm:"column:type"`
	Type                 *string           `json:"type" gorm:"-"`
	Amount               *float64          `json:"amount"`
	CurrencyInternal     *int              `json:"-" gorm:"column:currency"`
	Currency             *string           `json:"currency" gorm:"-"`
	ExchangeRate         *float64          `json:"exchange_rate" `
	OurSignatoryInternal *int              `json:"-" gorm:"column:our_signatory"`
	OurSignatory         *string           `json:"our_signatory" gorm:"-"`
	RelatedPartyID       *int              `json:"related_party_id" `
	DepartmentID         *int              `json:"-" `
	SigningDate          *string           `json:"signing_date" `
	EffectiveDate        *string           `json:"effective_date" `
	Content              *string           `json:"content"`
	Department           *DepartmentOutput `json:"department"`
}

func (p *ProjectGet) Get() response.Common {
	var result ProjectOutput
	err := global.DB.Model(model.Project{}).
		Where("id = ?", p.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
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

	//如果有部门id，就查部门信息
	if result.DepartmentID != nil {
		err = global.DB.Model(model.Department{}).
			Where("id=?", result.DepartmentID).First(&result.Department).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			result.Department = nil
		}
	}
	return response.SucceedWithData(result)
}

func (p *ProjectCreate) Create() response.Common {
	var paramOut model.Project

	if p.Creator > 0 {
		paramOut.Creator = &p.Creator
	}

	if p.LastModifier > 0 {
		paramOut.LastModifier = &p.LastModifier
	}

	if p.Code != "" {
		paramOut.Code = &p.Code
	}

	if p.Name != "" {
		paramOut.Name = &p.Name
	}

	if p.Country > 0 {
		paramOut.Country = &p.Country
	}

	if p.Province > 0 {
		paramOut.Province = &p.Province
	}

	if p.Type > 0 {
		paramOut.Type = &p.Type
	}

	if p.Amount != nil {
		paramOut.Amount = p.Amount
	}

	if p.Currency > 0 {
		paramOut.Currency = &p.Currency
	}

	if p.ExchangeRate != nil {
		paramOut.ExchangeRate = p.ExchangeRate
	}

	if p.DepartmentID != 0 {
		paramOut.DepartmentID = &p.DepartmentID
	}

	if p.RelatedPartyID != 0 {
		paramOut.RelatedPartyID = &p.RelatedPartyID
	}

	if p.SigningDate != "" {
		signingDate, err := time.Parse("2006-01-02", p.SigningDate)
		if err != nil {
			return response.Fail(util.ErrorInvalidDateFormat)
		}
		paramOut.SigningDate = &signingDate
	}

	if p.EffectiveDate != "" {
		effectiveDate, err := time.Parse("2006-01-02", p.EffectiveDate)
		if err != nil {
			return response.Fail(util.ErrorInvalidDateFormat)
		}
		paramOut.EffectiveDate = &effectiveDate
	}

	if p.Content != "" {
		paramOut.Content = &p.Content
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "Deleter", "CreateAt", "UpdatedAt", "DeletedAt")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeCreatedNotFound)
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (p *ProjectUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if p.LastModifier > 0 {
		paramOut["last_modifier"] = p.LastModifier
	}

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

	if p.Country != nil {
		if *p.Country > 0 {
			paramOut["country"] = p.Country
		} else {
			paramOut["country"] = nil
		}
	}

	if p.Province != nil {
		if *p.Province > 0 {
			paramOut["province"] = p.Province
		} else {
			paramOut["province"] = nil
		}
	}

	if p.Type != nil {
		if *p.Type > 0 {
			paramOut["type"] = p.Type
		} else {
			paramOut["type"] = nil
		}
	}

	if p.Amount != nil {
		if *p.Amount != -1 {
			paramOut["amount"] = p.Amount
		} else {
			paramOut["amount"] = nil
		}
	}

	if p.Currency != nil {
		if *p.Currency > 0 {
			paramOut["currency"] = p.Currency
		} else {
			paramOut["currency"] = nil
		}
	}

	if p.ExchangeRate != nil {
		if *p.ExchangeRate != -1 {
			paramOut["exchange_rate"] = p.ExchangeRate
		} else {
			paramOut["exchange_rate"] = nil
		}
	}

	if p.DepartmentID != nil {
		if *p.DepartmentID != 0 {
			paramOut["department_id"] = p.DepartmentID
		} else {
			paramOut["department_id"] = nil
		}
	}

	if p.RelatedPartyID != nil {
		if *p.RelatedPartyID != 0 {
			paramOut["related_party_id"] = p.RelatedPartyID
		} else {
			paramOut["related_party_id"] = nil
		}
	}

	if p.SigningDate != nil {
		if *p.SigningDate != "" {
			var err error
			paramOut["signing_date"], err = time.Parse("2006-01-02", *p.SigningDate)
			if err != nil {
				return response.Fail(util.ErrorInvalidJSONParameters)
			}
		}
	}

	if p.EffectiveDate != nil {
		if *p.EffectiveDate != "" {
			var err error
			paramOut["effective_date"], err = time.Parse("2006-01-02", *p.EffectiveDate)
			if err != nil {
				return response.Fail(util.ErrorInvalidJSONParameters)
			}
		}
	}

	if p.Content != nil {
		if *p.Content != "" {
			paramOut["content"] = p.Content
		} else {
			paramOut["content"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "Deleter", "CreateAt", "UpdatedAt", "DeletedAt")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Project{}).Where("id = ?", p.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (p *ProjectDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.Project
	global.DB.Where("id = ?", p.ID).Find(&record)
	record.Deleter = &p.Deleter
	err := global.DB.Where("id = ?", p.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (p *ProjectGetArray) GetArray() response.Common {
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

	// count
	var count int64
	db.Count(&count)

	//Order
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
			return response.Fail(util.ErrorSortingFieldDoesNotExist)
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
	if p.PagingInput.PageSize > 0 &&
		p.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = p.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//array
	var array []string
	db.Model(&model.Project{}).Select("project_full_name").Find(&array)

	if len(array) == 0 {
		return response.Fail(util.ErrorRecordNotFound)
	}

	return response.Common{
		Data:    array,
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
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

	// count
	var count int64
	db.Count(&count)

	//Order
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
			return response.FailForList(util.ErrorSortingFieldDoesNotExist)
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
	if p.PagingInput.PageSize > 0 &&
		p.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = p.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []ProjectGetListOutput
	db.Model(&model.Project{}).Find(&data)

	if len(data) == 0 {
		return response.FailForList(util.ErrorRecordNotFound)
	}

	for i := range data {
		if data[i].DepartmentID != nil && *data[i].DepartmentID > 0 {
			departmentID := *data[i].DepartmentID
			global.DB.Model(&model.Department{}).Where("id = ?", departmentID).
				Limit(1).Find(&data[i].Department)
		}

		//查找字典里相应的值
		if data[i].CountryInternal != nil {
			var country string
			global.DB.Model(&model.DictionaryItem{}).Where("id = ?", *data[i].CountryInternal).
				Limit(1).Select("name").Find(&country)
			if country != "" {
				data[i].Country = &country
			}
		}

		if data[i].TypeInternal != nil {
			var type1 string
			global.DB.Model(&model.DictionaryItem{}).Where("id = ?", *data[i].TypeInternal).
				Limit(1).Select("name").Find(&type1)
			if type1 != "" {
				data[i].Type = &type1
			}
		}

		if data[i].CurrencyInternal != nil {
			var currency string
			global.DB.Model(&model.DictionaryItem{}).Where("id = ?", *data[i].CountryInternal).
				Limit(1).Select("name").Find(&currency)
			if currency != "" {
				data[i].Currency = &currency
			}
		}

		if data[i].OurSignatoryInternal != nil {
			var ourSignatory string
			global.DB.Model(&model.DictionaryItem{}).Where("id = ?", *data[i].OurSignatoryInternal).
				Limit(1).Select("name").Find(&ourSignatory)
			if ourSignatory != "" {
				data[i].OurSignatory = &ourSignatory
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
