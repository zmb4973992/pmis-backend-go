package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type DepartmentGet struct {
	DepartmentID int
}

func (d DepartmentGet) Do() response.Common {
	var result dto.DepartmentOutput
	err := global.DB.Model(model.Department{}).
		Where("id = ?", d.DepartmentID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}

	return response.SucceedWithData(result)
}

type IGet interface {
	Do() response.Common
}

type factoryForDepartmentGet struct{}

func (f factoryForDepartmentGet) create() IGet {
	return DepartmentGet{
		DepartmentID: 0,
	}
}

func NewDepartmentGet() IGet {
	return DepartmentGet{}
}

func test() {
	departmentGet := new(DepartmentGet)
	departmentGet.Do()
}

//
//func (departmentService) DepartmentGet(DepartmentID int) response.Common {
//	var result dto.DepartmentOutput
//
//	err := global.DB.Model(model.Department{}).
//		Where("id = ?", DepartmentID).First(&result).Error
//	if err != nil {
//		global.SugaredLogger.Errorln(err)
//		return response.Fail(util.ErrorRecordNotFound)
//	}
//
//	return response.SucceedWithData(result)
//}
//
//func (departmentService) Create(paramIn dto.DepartmentCreate) response.Common {
//	var paramOut model.Department
//
//	if paramIn.Creator > 0 {
//		paramOut.Creator = &paramIn.Creator
//	}
//
//	if paramIn.LastModifier > 0 {
//		paramOut.LastModifier = &paramIn.LastModifier
//	}
//
//	paramOut.Name = paramIn.Name
//
//	paramOut.LevelName = paramIn.LevelName
//
//	paramOut.SuperiorID = &paramIn.SuperiorID
//
//	err := global.DB.Create(&paramOut).Error
//	if err != nil {
//		global.SugaredLogger.Errorln(err)
//		return response.Fail(util.ErrorFailToCreateRecord)
//	}
//	return response.Succeed()
//}
//
//func (departmentService) Update(paramIn dto.DepartmentUpdate) response.Common {
//	paramOut := make(map[string]any)
//
//	if paramIn.LastModifier > 0 {
//		paramOut["last_modifier"] = paramIn.LastModifier
//	}
//
//	if paramIn.Name != nil {
//		if *paramIn.Name != "" {
//			paramOut["name"] = paramIn.Name
//		} else {
//			paramOut["name"] = nil
//		}
//	}
//
//	if paramIn.LevelName != nil {
//		if *paramIn.LevelName != "" {
//			paramOut["level_name"] = paramIn.LevelName
//		} else {
//			paramOut["level_name"] = nil
//		}
//	}
//
//	if paramIn.SuperiorID != nil {
//		if *paramIn.SuperiorID > 0 {
//			paramOut["superior_id"] = paramIn.SuperiorID
//		} else if *paramIn.SuperiorID == 0 {
//			paramOut["superior_id"] = nil
//		} else {
//			return response.Fail(util.ErrorInvalidJSONParameters)
//		}
//	}
//
//	//计算有修改值的字段数，分别进行不同处理
//	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")
//
//	if len(paramOutForCounting) == 0 {
//		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
//	}
//
//	err := global.DB.Model(&model.Department{}).
//		Where("id = ?", paramIn.ID).Updates(paramOut).Error
//	if err != nil {
//		global.SugaredLogger.Errorln(err)
//		return response.Fail(util.ErrorFailToUpdateRecord)
//	}
//
//	return response.Succeed()
//}
//
//func (departmentService) Delete(paramIn dto.DepartmentDelete) response.Common {
//	//由于删除需要做两件事：软删除+记录删除人，所以需要用事务
//	err := global.DB.Transaction(func(tx *gorm.DB) error {
//		//这里记录删除人，在事务中必须放在前面
//		//如果放后面，由于是软删除，系统会找不到这条记录，导致无法更新
//		err := tx.Debug().Model(&model.Department{}).Where("id = ?", paramIn.ID).
//			Update("deleter", paramIn.Deleter).Error
//		if err != nil {
//			return err
//		}
//		//这里删除记录
//		err = tx.Delete(&model.Department{}, paramIn.ID).Error
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//
//	if err != nil {
//		global.SugaredLogger.Errorln(err)
//		return response.Fail(util.ErrorFailToDeleteRecord)
//	}
//	return response.Succeed()
//}
//
//func (departmentService) GetArray(paramIn dto.DepartmentList) response.Common {
//	db := global.DB.Model(&model.Department{})
//	// 顺序：where -> count -> order -> limit -> offset -> data
//
//	//where
//	if paramIn.SuperiorID > 0 {
//		db = db.Where("superior_id = ?", paramIn.SuperiorID)
//	}
//
//	if paramIn.LevelName != "" {
//		db = db.Where("level_name = ?", paramIn.LevelName)
//	}
//
//	if paramIn.Name != "" {
//		db = db.Where("name = ?", paramIn.Name)
//	}
//
//	if paramIn.NameLike != "" {
//		db = db.Where("name like ?", "%"+paramIn.NameLike+"%")
//	}
//
//	if paramIn.IsShowedByRole {
//		biggestRoleName := util.GetBiggestRoleName(paramIn.UserID)
//		if biggestRoleName == "事业部级" {
//			businessDivisionIDs := util.GetBusinessDivisionIDs(paramIn.UserID)
//			db = db.Where("superior_id in ?", businessDivisionIDs)
//		} else if biggestRoleName == "部门级" || biggestRoleName == "项目级" {
//			departmentIDs := util.GetDepartmentIDs(paramIn.UserID)
//			db = db.Where("id in ?", departmentIDs)
//		}
//	}
//
//	// count
//	var count int64
//	db.Count(&count)
//
//	//order
//	orderBy := paramIn.SortingInput.OrderBy
//	desc := paramIn.SortingInput.Desc
//	//如果排序字段为空
//	if orderBy == "" {
//		//如果要求降序排列
//		if desc == true {
//			db = db.Order("id desc")
//		}
//	} else { //如果有排序字段
//		//先看排序字段是否存在于表中
//		exists := util.FieldIsInModel(model.Department{}, orderBy)
//		if !exists {
//			return response.Fail(util.ErrorSortingFieldDoesNotExist)
//		}
//		//如果要求降序排列
//		if desc == true {
//			db = db.Order(orderBy + " desc")
//		} else { //如果没有要求排序方式
//			db = db.Order(orderBy)
//		}
//	}
//
//	//limit
//	page := 1
//	if paramIn.PagingInput.Page > 0 {
//		page = paramIn.PagingInput.Page
//	}
//	pageSize := global.Config.DefaultPageSize
//	if paramIn.PagingInput.PageSize > 0 &&
//		paramIn.PagingInput.PageSize <= global.Config.MaxPageSize {
//		pageSize = paramIn.PagingInput.PageSize
//	}
//	db = db.Limit(pageSize)
//
//	//offset
//	offset := (page - 1) * pageSize
//	db = db.Offset(offset)
//
//	//array
//	var array []string
//	db.Model(&model.DictionaryType{}).Select("name").Find(&array)
//
//	if len(array) == 0 {
//		return response.Fail(util.ErrorRecordNotFound)
//	}
//
//	return response.Common{
//		Data:    array,
//		Code:    util.Success,
//		Message: util.GetMessage(util.Success),
//	}
//}
//
//func (departmentService) List(paramIn dto.DepartmentList) response.List {
//	db := global.DB.Model(&model.Department{})
//	// 顺序：where -> count -> order -> limit -> offset -> data
//
//	//where
//	if paramIn.SuperiorID > 0 {
//		db = db.Where("superior_id = ?", paramIn.SuperiorID)
//	}
//
//	if paramIn.LevelName != "" {
//		db = db.Where("level_name = ?", paramIn.LevelName)
//	}
//
//	if paramIn.Name != "" {
//		db = db.Where("name = ?", paramIn.Name)
//	}
//
//	if paramIn.NameLike != "" {
//		db = db.Where("name like ?", "%"+paramIn.NameLike+"%")
//	}
//
//	if paramIn.IsShowedByRole {
//		biggestRoleName := util.GetBiggestRoleName(paramIn.UserID)
//		if biggestRoleName == "事业部级" {
//			businessDivisionIDs := util.GetBusinessDivisionIDs(paramIn.UserID)
//			db = db.Where("superior_id in ?", businessDivisionIDs)
//		} else if biggestRoleName == "部门级" || biggestRoleName == "项目级" {
//			departmentIDs := util.GetDepartmentIDs(paramIn.UserID)
//			db = db.Where("id in ?", departmentIDs)
//		}
//	}
//
//	// count
//	var count int64
//	db.Count(&count)
//
//	//order
//	orderBy := paramIn.SortingInput.OrderBy
//	desc := paramIn.SortingInput.Desc
//	//如果排序字段为空
//	if orderBy == "" {
//		//如果要求降序排列
//		if desc == true {
//			db = db.Order("id desc")
//		}
//	} else { //如果有排序字段
//		//先看排序字段是否存在于表中
//		exists := util.FieldIsInModel(model.Department{}, orderBy)
//		if !exists {
//			return response.FailForList(util.ErrorSortingFieldDoesNotExist)
//		}
//		//如果要求降序排列
//		if desc == true {
//			db = db.Order(orderBy + " desc")
//		} else { //如果没有要求排序方式
//			db = db.Order(orderBy)
//		}
//	}
//
//	//limit
//	page := 1
//	if paramIn.PagingInput.Page > 0 {
//		page = paramIn.PagingInput.Page
//	}
//	pageSize := global.Config.DefaultPageSize
//	if paramIn.PagingInput.PageSize > 0 &&
//		paramIn.PagingInput.PageSize <= global.Config.MaxPageSize {
//		pageSize = paramIn.PagingInput.PageSize
//	}
//	db = db.Limit(pageSize)
//
//	//offset
//	offset := (page - 1) * pageSize
//	db = db.Offset(offset)
//
//	//data
//	var data []dto.DepartmentOutput
//	db.Model(&model.Department{}).Find(&data)
//
//	if len(data) == 0 {
//		return response.FailForList(util.ErrorRecordNotFound)
//	}
//
//	numberOfRecords := int(count)
//	numberOfPages := util.GetTotalNumberOfPages(numberOfRecords, pageSize)
//
//	return response.List{
//		Data: data,
//		Paging: &dto.PagingOutput{
//			Page:            page,
//			PageSize:        pageSize,
//			NumberOfPages:   numberOfPages,
//			NumberOfRecords: numberOfRecords,
//		},
//		Code:    util.Success,
//		Message: util.GetMessage(util.Success),
//	}
//}
