package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
	"strconv"
	"strings"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type ContractGet struct {
	ContractID int64
	UserID     int64
}

type ContractCreate struct {
	UserID int64
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
	FileIDs     string `json:"file_ids,omitempty"`
	Operator    string `json:"operator,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ContractUpdate struct {
	UserID     int64
	ContractID int64
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
	FileIDs     []int64 `json:"file_ids"`
	Operator    *string `json:"operator"`

	IgnoreDataAuthority bool `json:"-"`
}

type ContractDelete struct {
	UserID     int64
	ContractID int64
}

type ContractGetList struct {
	list.Input
	UserID         int64  `json:"-"`
	ProjectID      int64  `json:"project_id,omitempty"`
	RelatedPartyID int64  `json:"related_party_id,omitempty"`
	FundDirection  int64  `json:"fund_direction,omitempty"`
	NameInclude    string `json:"name_include,omitempty"`

	//是否忽略数据权限的限制，用于请求数据范围外的全部数据
	IgnoreDataAuthority bool `json:"ignore_data_authority"`
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

	Name          *string      `json:"name"`
	Code          *string      `json:"code"`
	Content       *string      `json:"content"`
	Deliverable   *string      `json:"deliverable"`
	PenaltyRule   *string      `json:"penalty_rule"`
	Operator      *string      `json:"operator"`
	FileIDs       *string      `json:"-"`
	FilesExternal []FileOutput `json:"files" gorm:"-"`

	//用来告诉前端，该记录是否为数据范围内，用来判定是否能访问详情、是否需要做跳转等
	Authorized bool `json:"authorized" gorm:"-"`
}

type contractCheckAuth struct {
	UserID     int64
	ContractID int64
}

func (c *ContractGet) Get() (output *ContractOutput, errCode int) {
	err := global.DB.Model(model.Contract{}).
		Where("id = ?", c.ContractID).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	var auth contractCheckAuth
	auth.ContractID = c.ContractID
	auth.UserID = c.UserID
	authorized := auth.checkAuth()

	if !authorized {
		return nil, util.ErrorUnauthorized
	}

	//查询关联表的详情
	{
		//查项目信息
		if output.ProjectID != nil {
			var record *ProjectOutput
			res := global.DB.Model(&model.Project{}).
				Where("id = ?", output.ProjectID).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.ProjectExternal = record
			}
		}
		//查部门信息
		if output.OrganizationID != nil {
			var record *OrganizationOutput
			res := global.DB.Model(&model.Organization{}).
				Where("id = ?", output.OrganizationID).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.OrganizationExternal = record
			}
		}
		//查相关方信息
		if output.RelatedPartyID != nil {
			var record *RelatedPartyOutput
			res := global.DB.Model(&model.RelatedParty{}).
				Where("id = ?", output.RelatedPartyID).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.RelatedPartyExternal = record
			}
		}

		//查文件信息
		if output.FileIDs != nil {
			tempFileIDs := strings.Split(*output.FileIDs, ",")
			var fileIDs []int64
			for i := range tempFileIDs {
				fileID, err1 := strconv.ParseInt(tempFileIDs[i], 10, 64)
				if err1 != nil {
					continue
				}
				fileIDs = append(fileIDs, fileID)
			}

			var records []FileOutput
			global.DB.Model(&model.File{}).
				Where("id in ?", fileIDs).
				Find(&records)

			ip := global.Config.DownloadConfig.LocalIP
			port := global.Config.AppConfig.HttpPort
			accessPath := global.Config.DownloadConfig.RelativePath
			for i := range records {
				records[i].Url = "http://" + ip + ":" + port + accessPath + strconv.FormatInt(records[i].ID, 10)
			}
			output.FilesExternal = records
		}
	}

	//查询dictionary_item表的详情
	{
		if output.FundDirection != nil {
			var record *DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", output.FundDirection).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.FundDirectionExternal = record
			}
		}
		if output.Currency != nil {
			var record *DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", output.Currency).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.CurrencyExternal = record
			}
		}
		if output.OurSignatory != nil {
			var record *DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", output.OurSignatory).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.OurSignatoryExternal = record
			}
		}
		if output.Type != nil {
			var record *DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", output.Type).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.TypeExternal = record
			}
		}
	}

	//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
	//需要取年月日(即前9位)
	{
		if output.SigningDate != nil {
			temp := *output.SigningDate
			temp1 := temp[:10]
			output.SigningDate = &temp1
		}
		if output.EffectiveDate != nil {
			temp := *output.EffectiveDate
			temp1 := temp[:10]
			output.EffectiveDate = &temp1
		}
		if output.CommissioningDate != nil {
			temp := *output.CommissioningDate
			temp1 := temp[:10]
			output.CommissioningDate = &temp1
		}
		if output.CompletionDate != nil {
			temp := *output.CompletionDate
			temp1 := temp[:10]
			output.CompletionDate = &temp1
		}
	}

	return output, util.Success
}

func (c *ContractCreate) Create() (errCode int) {
	var paramOut model.Contract

	if c.UserID > 0 {
		paramOut.Creator = &c.UserID
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
				return util.ErrorInvalidDateFormat
			}
			paramOut.SigningDate = &signingDate
		}

		if c.EffectiveDate != "" {
			effectiveDate, err := time.Parse("2006-01-02", c.EffectiveDate)
			if err != nil {
				return util.ErrorInvalidDateFormat
			}
			paramOut.EffectiveDate = &effectiveDate
		}

		if c.CommissioningDate != "" {
			commissioningDate, err := time.Parse("2006-01-02", c.CommissioningDate)
			if err != nil {
				return util.ErrorInvalidDateFormat
			}
			paramOut.CommissioningDate = &commissioningDate
		}

		if c.CompletionDate != "" {
			completionDate, err := time.Parse("2006-01-02", c.CompletionDate)
			if err != nil {
				return util.ErrorInvalidDateFormat
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

		if c.FileIDs != "" {
			paramOut.FileIDs = &c.FileIDs
		}

		if c.Operator != "" {
			paramOut.Operator = &c.Operator
		}
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}
	return util.Success
}

func (c *ContractUpdate) Update() (errCode int) {
	var result ContractOutput
	err := global.DB.Model(model.Contract{}).
		Where("id = ?", c.ContractID).
		First(&result).Error
	if err != nil {
		return util.ErrorRecordNotFound
	}

	if c.IgnoreDataAuthority == false {
		var authorization contractCheckAuth
		authorization.ContractID = c.ContractID
		authorization.UserID = c.UserID
		authorized := authorization.checkAuth()
		if !authorized {
			return util.ErrorUnauthorized
		}
	}

	paramOut := make(map[string]any)

	if c.UserID > 0 {
		paramOut["last_modifier"] = c.UserID
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
				var err1 error
				paramOut["signing_date"], err1 = time.Parse("2006-01-02", *c.SigningDate)
				if err1 != nil {
					return util.ErrorInvalidJSONParameters
				}
			} else {
				paramOut["signing_date"] = nil
			}
		}
		if c.EffectiveDate != nil {
			if *c.EffectiveDate != "" {
				var err1 error
				paramOut["effective_date"], err1 = time.Parse("2006-01-02", *c.EffectiveDate)
				if err1 != nil {
					return util.ErrorInvalidJSONParameters
				}
			} else {
				paramOut["effective_date"] = nil
			}
		}
		if c.CommissioningDate != nil {
			if *c.CommissioningDate != "" {
				var err1 error
				paramOut["commissioning_date"], err1 = time.Parse("2006-01-02", *c.CommissioningDate)
				if err1 != nil {
					return util.ErrorInvalidJSONParameters
				}
			} else {
				paramOut["commissioning_date"] = nil
			}
		}
		if c.CompletionDate != nil {
			if *c.CompletionDate != "" {
				var err1 error
				paramOut["completion_date"], err1 = time.Parse("2006-01-02", *c.CompletionDate)
				if err1 != nil {
					return util.ErrorInvalidJSONParameters
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

		if c.FileIDs != nil {
			if len(c.FileIDs) > 0 {
				var fileIDs []string
				for _, v := range c.FileIDs {
					fileIDs = append(fileIDs, strconv.FormatInt(v, 10))
				}
				paramOut["file_ids"] = strings.Join(fileIDs, ",")
			} else {
				paramOut["file_ids"] = nil
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

	err = global.DB.Model(&model.Contract{}).
		Where("id = ?", c.ContractID).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}

func (c *ContractDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.Contract

	err := global.DB.Where("id = ?", c.ContractID).
		Find(&record).
		Delete(&record).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	return util.Success
}

func (c *ContractGetList) GetList() (outputs []ContractOutput,
	errCode int, paging *list.PagingOutput) {
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
	if c.IgnoreDataAuthority == false {
		organizationIDs := util.GetOrganizationIDsForDataAuthority(c.UserID)
		//先找出项目的数据范围
		var projectIDs []int64
		global.DB.Model(&model.Project{}).Where("organization_id in ?", organizationIDs).
			Select("id").Find(&projectIDs)
		//然后再加上组织的数据范围
		db = db.Where("organization_id in ? or project_id in ?",
			organizationIDs, projectIDs)
	}

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

	//outputs
	db.Model(&model.Contract{}).Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	for i := range outputs {
		//查询关联表的详情
		{
			//查项目信息
			if outputs[i].ProjectID != nil {
				var record ProjectOutput
				res := global.DB.Model(&model.Project{}).
					Where("id = ?", *outputs[i].ProjectID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].ProjectExternal = &record
				}
			}
			//查部门信息
			if outputs[i].OrganizationID != nil {
				var record OrganizationOutput
				res := global.DB.Model(&model.Organization{}).
					Where("id = ?", *outputs[i].OrganizationID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].OrganizationExternal = &record
				}
			}
			//查相关方信息
			if outputs[i].RelatedPartyID != nil {
				var record RelatedPartyOutput
				res := global.DB.Model(&model.RelatedParty{}).
					Where("id = ?", *outputs[i].RelatedPartyID).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].RelatedPartyExternal = &record
				}
			}
		}

		//查dictionary_item表的详情
		{
			if outputs[i].FundDirection != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *outputs[i].FundDirection).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].FundDirectionExternal = &record
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
			if outputs[i].Currency != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *outputs[i].Currency).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].CurrencyExternal = &record
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
		}

		//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
		//需要取年月日(即前9位)
		{
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
			if outputs[i].CompletionDate != nil {
				temp := *outputs[i].CompletionDate
				*outputs[i].CompletionDate = temp[:10]
			}
		}

		if c.IgnoreDataAuthority == true {
			var authorize contractCheckAuth
			authorize.ContractID = outputs[i].ID
			authorize.UserID = c.UserID
			outputs[i].Authorized = authorize.checkAuth()
		} else {
			outputs[i].Authorized = true
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

// 该方法一定要在确定记录存在后再调用
func (c *contractCheckAuth) checkAuth() (authorized bool) {
	//用来确定数据范围内的组织id
	organizationIDs := util.GetOrganizationIDsForDataAuthority(c.UserID)

	if len(organizationIDs) == 0 {
		return false
	}

	//看看在数据范围内是否有该记录
	var count int64
	global.DB.Model(model.Contract{}).
		Joins("join (select id as project_id from project where organization_id in ?) as temp1 on contract.project_id = temp1.project_id", organizationIDs).
		Where("id = ?", c.ContractID).
		Count(&count)

	if count > 0 {
		return true
	}

	return false
}
