package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type user struct{}

func (*user) Get(userID int) response.Common {
	var result dto.UserOutput
	err := global.DB.Model(model.User{}).
		Where("id = ?", userID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}
	return response.SucceedWithData(result)
}

func (*user) Create(paramIn dto.UserCreate) response.Common {
	var paramOut model.User
	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	paramOut.Username = paramIn.Username
	encryptedPassword, err := util.Encrypt(paramIn.Password)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToEncrypt)
	}

	paramOut.Password = encryptedPassword

	if paramIn.IsValid != nil {
		paramOut.IsValid = paramIn.IsValid
	}

	if paramIn.FullName != "" {
		paramOut.FullName = &paramIn.FullName
	}

	if paramIn.EmailAddress != "" {
		paramOut.EmailAddress = &paramIn.EmailAddress
	}

	if paramIn.MobilePhoneNumber != "" {
		paramOut.MobilePhoneNumber = &paramIn.MobilePhoneNumber
	}
	if paramIn.EmployeeNumber != "" {
		paramOut.EmployeeNumber = &paramIn.EmployeeNumber
	}

	err = global.DB.Create(&paramOut).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (*user) Update(paramIn dto.UserUpdate) response.Common {
	paramOut := make(map[string]any)

	if paramIn.LastModifier > 0 {
		paramOut["last_modifier"] = paramIn.LastModifier
	}

	if paramIn.FullName != nil {
		if *paramIn.FullName != "" {
			paramOut["full_name"] = paramIn.FullName
		} else {
			paramOut["full_name"] = nil
		}
	}

	if paramIn.EmailAddress != nil {
		if *paramIn.EmailAddress != "" {
			paramOut["email_address"] = paramIn.EmailAddress
		} else {
			paramOut["email_address"] = nil
		}
	}

	if paramIn.IsValid != nil {
		paramOut["is_valid"] = paramIn.IsValid
	}

	if paramIn.MobilePhoneNumber != nil {
		if *paramIn.MobilePhoneNumber != "" {
			paramOut["mobile_phone_number"] = paramIn.MobilePhoneNumber
		} else {
			paramOut["mobile_phone_number"] = nil
		}
	}

	if paramIn.EmployeeNumber != nil {
		if *paramIn.EmployeeNumber != "" {
			paramOut["employee_number"] = paramIn.EmployeeNumber
		} else {
			paramOut["employee_number"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.User{}).Where("id = ?", paramIn.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (*user) Delete(paramIn dto.UserDelete) response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.User
	global.DB.Where("id = ?", paramIn.ID).Find(&record)
	record.Deleter = &paramIn.Deleter
	err := global.DB.Where("id = ?", paramIn.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (*user) List(paramIn dto.UserList) response.List {
	db := global.DB.Model(&model.User{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if paramIn.IsValid != nil {
		db = db.Where("is_valid = ?", *paramIn.IsValid)
	}

	if paramIn.UsernameInclude != "" {
		db = db.Where("username like ?", "%"+paramIn.UsernameInclude+"%")
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
		exists := util.FieldIsInModel(&model.User{}, orderBy)
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
	var data []dto.UserOutput
	db.Model(&model.User{}).Find(&data)

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
