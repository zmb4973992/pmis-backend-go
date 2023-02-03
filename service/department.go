package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type department struct{}

func (*department) Get(departmentID int) response.Common {
	var result dto.DepartmentOutput

	err := global.DB.Model(model.Department{}).
		Where("id = ?", departmentID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}

	return response.SucceedWithData(result)
}

func (*department) Create(paramIn dto.DepartmentCreate) response.Common {
	var paramOut model.Department

	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	paramOut.Name = paramIn.Name

	paramOut.LevelName = paramIn.LevelName

	paramOut.SuperiorID = &paramIn.SuperiorID

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (*department) Update(paramIn dto.DepartmentUpdate) response.Common {
	paramOut := make(map[string]any)

	if paramIn.LastModifier > 0 {
		paramOut["last_modifier"] = paramIn.LastModifier
	}

	if paramIn.Name != nil {
		if *paramIn.Name != "" {
			paramOut["name"] = paramIn.Name
		} else {
			paramOut["name"] = nil
		}
	}

	if paramIn.LevelName != nil {
		if *paramIn.LevelName != "" {
			paramOut["level_name"] = paramIn.LevelName
		} else {
			paramOut["level_name"] = nil
		}
	}

	if paramIn.SuperiorID != nil {
		if *paramIn.SuperiorID > 0 {
			paramOut["superior_id"] = paramIn.SuperiorID
		} else if *paramIn.SuperiorID == 0 {
			paramOut["superior_id"] = nil
		} else {
			return response.Fail(util.ErrorInvalidJSONParameters)
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Department{}).
		Where("id = ?", paramIn.ID).Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (*department) Delete(paramIn dto.DepartmentDelete) response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.Department
	global.DB.Where("id = ?", paramIn.ID).Find(&record)
	record.Deleter = &paramIn.Deleter
	err := global.DB.Where("id = ?", paramIn.ID).Delete(&record).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (*department) GetArray(paramIn dto.DepartmentList) response.Common {
	db := global.DB.Model(&model.Department{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if paramIn.SuperiorID > 0 {
		db = db.Where("superior_id = ?", paramIn.SuperiorID)
	}

	if paramIn.LevelName != "" {
		db = db.Where("level_name = ?", paramIn.LevelName)
	}

	if paramIn.Name != "" {
		db = db.Where("name = ?", paramIn.Name)
	}

	if paramIn.NameLike != "" {
		db = db.Where("name like ?", "%"+paramIn.NameLike+"%")
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
		exists := util.FieldIsInModel(&model.Department{}, orderBy)
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
	db.Model(&model.DictionaryType{}).Select("name").Find(&array)

	if len(array) == 0 {
		return response.Fail(util.ErrorRecordNotFound)
	}

	return response.Common{
		Data:    array,
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}

func (*department) List(paramIn dto.DepartmentList) response.List {
	db := global.DB.Model(&model.Department{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if paramIn.SuperiorID > 0 {
		db = db.Where("superior_id = ?", paramIn.SuperiorID)
	}

	if paramIn.LevelName != "" {
		db = db.Where("level_name = ?", paramIn.LevelName)
	}

	if paramIn.Name != "" {
		db = db.Where("name = ?", paramIn.Name)
	}

	if paramIn.NameLike != "" {
		db = db.Where("name like ?", "%"+paramIn.NameLike+"%")
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
		exists := util.FieldIsInModel(&model.Department{}, orderBy)
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
	var data []dto.DepartmentOutput
	db.Model(&model.Department{}).Find(&data)

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
