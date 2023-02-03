package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type relatedParty struct{}

func (*relatedParty) Get(relatedPartyID int) response.Common {
	var result dto.RelatedPartyOutput
	err := global.DB.Model(&model.RelatedParty{}).
		Where("id = ?", relatedPartyID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}
	return response.SucceedWithData(result)
}

func (*relatedParty) Create(paramIn dto.RelatedPartyCreate) response.Common {
	var paramOut model.RelatedParty
	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	if paramIn.ChineseName != "" {
		paramOut.ChineseName = &paramIn.ChineseName
	}

	if paramIn.EnglishName != "" {
		paramOut.EnglishName = &paramIn.EnglishName
	}

	if paramIn.Address != "" {
		paramOut.Address = &paramIn.Address
	}

	if paramIn.UniformSocialCreditCode != "" {
		paramOut.UniformSocialCreditCode = &paramIn.UniformSocialCreditCode
	}

	if paramIn.Telephone != "" {
		paramOut.Telephone = &paramIn.Telephone
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		response.Fail(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "deleter", "CreateAt", "UpdatedAt", "DeletedAt")

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

func (*relatedParty) Update(paramIn dto.RelatedPartyUpdate) response.Common {
	paramOut := make(map[string]any)

	if paramIn.LastModifier > 0 {
		paramOut["last_modifier"] = paramIn.LastModifier
	}

	if paramIn.ChineseName != nil {
		if *paramIn.ChineseName != "" {
			paramOut["chinese_name"] = paramIn.ChineseName
		} else {
			paramOut["chinese_name"] = nil
		}
	}

	if paramIn.EnglishName != nil {
		if *paramIn.EnglishName != "" {
			paramOut["english_name"] = paramIn.EnglishName
		} else {
			paramOut["english_name"] = nil
		}
	}

	if paramIn.Address != nil {
		if *paramIn.Address != "" {
			paramOut["address"] = paramIn.Address
		} else {
			paramOut["address"] = nil
		}
	}

	if paramIn.UniformSocialCreditCode != nil {
		if *paramIn.UniformSocialCreditCode != "" {
			paramOut["uniform_social_credit_code"] = paramIn.UniformSocialCreditCode
		} else {
			paramOut["uniform_social_credit_code"] = nil
		}
	}

	if paramIn.Telephone != nil {
		if *paramIn.Telephone != "" {
			paramOut["telephone"] = paramIn.Telephone
		} else {
			paramOut["telephone"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.RelatedParty{}).Where("id = ?", paramIn.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (*relatedParty) Delete(paramIn dto.RelatedPartyDelete) response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.RelatedParty
	global.DB.Where("id = ?", paramIn.ID).Find(&record)
	record.Deleter = &paramIn.Deleter
	err := global.DB.Where("id = ?", paramIn.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (*relatedParty) GetList(paramIn dto.RelatedPartyList) response.List {
	db := global.DB.Model(&model.RelatedParty{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if paramIn.ChineseNameInclude != "" {
		db = db.Where("chinese_name like ?", "%"+paramIn.ChineseNameInclude+"%")
	}

	if paramIn.EnglishNameInclude != "" {
		db = db.Where("english_name like ?", "%"+paramIn.EnglishNameInclude+"%")
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := paramIn.SortingInput.OrderBy
	desc := paramIn.SortingInput.Desc
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
	if paramIn.PagingInput.Page > 0 {
		page = paramIn.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if paramIn.PagingInput.PageSize > 0 &&
		paramIn.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = paramIn.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []dto.RelatedPartyOutput
	db.Model(&model.RelatedParty{}).Find(&data)

	if len(data) == 0 {
		return response.FailForList(util.ErrorRecordNotFound)
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
