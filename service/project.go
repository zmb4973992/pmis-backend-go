package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type project struct{}

func (*project) Get(projectID int) response.Common {
	var result dto.ProjectOutput
	err := global.DB.Model(model.Project{}).
		Where("id = ?", projectID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
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

func (*project) Create(paramIn dto.ProjectCreate) response.Common {
	var paramOut model.Project

	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	if paramIn.ProjectCode != "" {
		paramOut.ProjectCode = &paramIn.ProjectCode
	}

	if paramIn.ProjectFullName != "" {
		paramOut.ProjectFullName = &paramIn.ProjectFullName
	}

	if paramIn.ProjectShortName != "" {
		paramOut.ProjectShortName = &paramIn.ProjectShortName
	}

	if paramIn.Country != "" {
		paramOut.Country = &paramIn.Country
	}

	if paramIn.Province != "" {
		paramOut.Province = &paramIn.Province
	}

	if paramIn.ProjectType != "" {
		paramOut.ProjectType = &paramIn.ProjectType
	}

	if paramIn.Amount != nil {
		paramOut.Amount = paramIn.Amount
	}

	if paramIn.Currency != "" {
		paramOut.Currency = &paramIn.Currency
	}

	if paramIn.ExchangeRate != nil {
		paramOut.ExchangeRate = paramIn.ExchangeRate
	}

	if paramIn.DepartmentID != 0 {
		paramOut.DepartmentID = &paramIn.DepartmentID
	}

	if paramIn.RelatedPartyID != 0 {
		paramOut.RelatedPartyID = &paramIn.RelatedPartyID
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (*project) Update(paramIn dto.ProjectUpdate) response.Common {
	paramOut := make(map[string]any)

	if paramIn.LastModifier > 0 {
		paramOut["last_modifier"] = paramIn.LastModifier
	}

	if paramIn.ProjectCode != nil {
		if *paramIn.ProjectCode != "" {
			paramOut["project_code"] = paramIn.ProjectCode
		} else {
			paramOut["project_code"] = nil
		}
	}

	if paramIn.ProjectFullName != nil {
		if *paramIn.ProjectFullName != "" {
			paramOut["project_full_name"] = paramIn.ProjectFullName
		} else {
			paramOut["project_full_name"] = nil
		}
	}

	if paramIn.ProjectShortName != nil {
		if *paramIn.ProjectShortName != "" {
			paramOut["project_short_name"] = paramIn.ProjectShortName
		} else {
			paramOut["project_short_name"] = nil
		}
	}

	if paramIn.Country != nil {
		if *paramIn.Country != "" {
			paramOut["country"] = paramIn.Country
		} else {
			paramOut["country"] = nil
		}
	}

	if paramIn.Province != nil {
		if *paramIn.Province != "" {
			paramOut["province"] = paramIn.Province
		} else {
			paramOut["province"] = nil
		}
	}

	if paramIn.ProjectType != nil {
		if *paramIn.ProjectType != "" {
			paramOut["project_type"] = paramIn.ProjectType
		} else {
			paramOut["project_type"] = nil
		}
	}

	if paramIn.Amount != nil {
		if *paramIn.Amount != -1 {
			paramOut["amount"] = paramIn.Amount
		} else {
			paramOut["amount"] = nil
		}
	}

	if paramIn.Currency != nil {
		if *paramIn.Currency != "" {
			paramOut["currency"] = paramIn.Currency
		} else {
			paramOut["currency"] = nil
		}
	}

	if paramIn.ExchangeRate != nil {
		if *paramIn.ExchangeRate != -1 {
			paramOut["exchange_rate"] = paramIn.ExchangeRate
		} else {
			paramOut["exchange_rate"] = nil
		}
	}

	if paramIn.DepartmentID != nil {
		if *paramIn.DepartmentID != 0 {
			paramOut["department_id"] = paramIn.DepartmentID
		} else {
			paramOut["department_id"] = nil
		}
	}

	if paramIn.RelatedPartyID != nil {
		if *paramIn.RelatedPartyID != 0 {
			paramOut["related_party_id"] = paramIn.RelatedPartyID
		} else {
			paramOut["related_party_id"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Project{}).Where("id = ?", paramIn.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (*project) Delete(paramIn dto.ProjectDelete) response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.Project
	global.DB.Where("id = ?", paramIn.ID).Find(&record)
	record.Deleter = &paramIn.Deleter
	err := global.DB.Where("id = ?", paramIn.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (*project) GetArray(paramIn dto.ProjectList) response.Common {
	db := global.DB.Model(&model.Project{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if paramIn.ProjectNameLike != "" {
		db = db.Where("project_full_name like ?", "%"+paramIn.ProjectNameLike+"%").
			Or("project_short_name like ?", "%"+paramIn.ProjectNameLike+"%")
	}

	if paramIn.DepartmentNameLike != "" {
		var departmentIDs []int
		global.DB.Model(&model.Department{}).Where("name like ?", "%"+paramIn.DepartmentNameLike+"%").
			Select("id").Find(&departmentIDs)
		if len(departmentIDs) > 0 {
			db = db.Where("department_id in ?", departmentIDs)
		}
	}

	if len(paramIn.DepartmentIDIn) > 0 {
		db = db.Where("department_id in ?", paramIn.DepartmentIDIn)
	}

	if paramIn.IsShowedByRole {
		biggestRoleName := util.GetBiggestRoleName(paramIn.UserID)
		if biggestRoleName == "事业部级" {
			businessDivisionIDs := util.GetBusinessDivisionIDs(paramIn.UserID)
			db = db.Where("superior_id in ?", businessDivisionIDs)
		} else if biggestRoleName == "部门级" || biggestRoleName == "项目级" {
			departmentIDs := util.GetDepartmentIDs(paramIn.UserID)
			db = db.Where("id in ?", departmentIDs)
		}
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

func (*project) GetList(paramIn dto.ProjectList) response.List {
	db := global.DB.Model(&model.Project{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if paramIn.ProjectNameLike != "" {
		db = db.Where("project_full_name like ?", "%"+paramIn.ProjectNameLike+"%").
			Or("project_short_name like ?", "%"+paramIn.ProjectNameLike+"%")
	}

	if paramIn.DepartmentNameLike != "" {
		var departmentIDs []int
		global.DB.Model(&model.Department{}).Where("name like ?", "%"+paramIn.DepartmentNameLike+"%").
			Select("id").Find(&departmentIDs)
		if len(departmentIDs) > 0 {
			db = db.Where("department_id in ?", departmentIDs)
		}
	}

	if len(paramIn.DepartmentIDIn) > 0 {
		db = db.Where("department_id in ?", paramIn.DepartmentIDIn)
	}

	if paramIn.IsShowedByRole {
		biggestRoleName := util.GetBiggestRoleName(paramIn.UserID)
		if biggestRoleName == "事业部级" {
			businessDivisionIDs := util.GetBusinessDivisionIDs(paramIn.UserID)
			db = db.Where("superior_id in ?", businessDivisionIDs)
		} else if biggestRoleName == "部门级" || biggestRoleName == "项目级" {
			departmentIDs := util.GetDepartmentIDs(paramIn.UserID)
			db = db.Where("id in ?", departmentIDs)
		}
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
	var data []dto.ProjectOutput
	db.Model(&model.Project{}).Find(&data)

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
