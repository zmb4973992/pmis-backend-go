package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type RelatedPartyGet struct {
	ID int
}

type RelatedPartyCreate struct {
	Creator      int
	LastModifier int

	ChineseName             string `json:"chinese_name,omitempty"`
	EnglishName             string `json:"english_name,omitempty"`
	Address                 string `json:"address,omitempty"`
	UniformSocialCreditCode string `json:"uniform_social_credit_code,omitempty"` //统一社会信用代码
	Telephone               string `json:"telephone,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type RelatedPartyUpdate struct {
	LastModifier int
	ID           int

	ChineseName             *string `json:"chinese_name"`
	EnglishName             *string `json:"english_name"`
	Address                 *string `json:"address"`
	UniformSocialCreditCode *string `json:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string `json:"telephone"`
}

type RelatedPartyDelete struct {
	Deleter int
	ID      int
}

type RelatedPartyGetList struct {
	dto.ListInput

	ChineseNameInclude string `json:"chinese_name_include,omitempty"`
	EnglishNameInclude string `json:"english_name_include,omitempty"`
}

//以下为出参

type RelatedPartyOutput struct {
	Creator      *int `json:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" gorm:"id"`

	ChineseName             *string `json:"chinese_name" gorm:"chinese_name"`
	EnglishName             *string `json:"english_name" gorm:"english_name"`
	Address                 *string `json:"address" gorm:"address"`
	UniformSocialCreditCode *string `json:"uniform_social_credit_code" gorm:"uniform_social_credit_code"` //统一社会信用代码
	Telephone               *string `json:"telephone" gorm:"telephone"`
}

func (r *RelatedPartyGet) Get() response.Common {
	var result RelatedPartyOutput
	err := global.DB.Model(&model.RelatedParty{}).
		Where("id = ?", r.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (r *RelatedPartyCreate) Create() response.Common {
	var paramOut model.RelatedParty
	if r.Creator > 0 {
		paramOut.Creator = &r.Creator
	}

	if r.LastModifier > 0 {
		paramOut.LastModifier = &r.LastModifier
	}

	if r.ChineseName != "" {
		paramOut.ChineseName = &r.ChineseName
	}

	if r.EnglishName != "" {
		paramOut.EnglishName = &r.EnglishName
	}

	if r.Address != "" {
		paramOut.Address = &r.Address
	}

	if r.UniformSocialCreditCode != "" {
		paramOut.UniformSocialCreditCode = &r.UniformSocialCreditCode
	}

	if r.Telephone != "" {
		paramOut.Telephone = &r.Telephone
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "deleter", "CreateAt", "UpdatedAt", "DeletedAt")

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

func (r *RelatedPartyUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if r.LastModifier > 0 {
		paramOut["last_modifier"] = r.LastModifier
	}

	if r.ChineseName != nil {
		if *r.ChineseName != "" {
			paramOut["chinese_name"] = r.ChineseName
		} else {
			paramOut["chinese_name"] = nil
		}
	}

	if r.EnglishName != nil {
		if *r.EnglishName != "" {
			paramOut["english_name"] = r.EnglishName
		} else {
			paramOut["english_name"] = nil
		}
	}

	if r.Address != nil {
		if *r.Address != "" {
			paramOut["address"] = r.Address
		} else {
			paramOut["address"] = nil
		}
	}

	if r.UniformSocialCreditCode != nil {
		if *r.UniformSocialCreditCode != "" {
			paramOut["uniform_social_credit_code"] = r.UniformSocialCreditCode
		} else {
			paramOut["uniform_social_credit_code"] = nil
		}
	}

	if r.Telephone != nil {
		if *r.Telephone != "" {
			paramOut["telephone"] = r.Telephone
		} else {
			paramOut["telephone"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.RelatedParty{}).Where("id = ?", r.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (r *RelatedPartyDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.RelatedParty
	global.DB.Where("id = ?", r.ID).Find(&record)
	record.Deleter = &r.Deleter
	err := global.DB.Where("id = ?", r.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (r *RelatedPartyGetList) GetList() response.List {
	db := global.DB.Model(&model.RelatedParty{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if r.ChineseNameInclude != "" {
		db = db.Where("chinese_name like ?", "%"+r.ChineseNameInclude+"%")
	}

	if r.EnglishNameInclude != "" {
		db = db.Where("english_name like ?", "%"+r.EnglishNameInclude+"%")
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := r.SortingInput.OrderBy
	desc := r.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.RelatedParty{}, orderBy)
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
	if r.PagingInput.Page > 0 {
		page = r.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if r.PagingInput.PageSize > 0 &&
		r.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = r.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []RelatedPartyOutput
	db.Model(&model.RelatedParty{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
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
