package service

import (
	"gorm.io/gorm"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type dictionaryType struct{}

func (dictionaryType) Get(dictionaryTypeID int) response.Common {
	var result dto.DictionaryTypeOutput
	err := global.DB.Model(model.DictionaryType{}).
		Where("id = ?", dictionaryTypeID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}
	return response.SucceedWithData(result)
}

func (dictionaryType) Create(paramIn dto.DictionaryTypeCreate) response.Common {
	var paramOut model.DictionaryType
	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	paramOut.Name = paramIn.Name

	if paramIn.Sort != 0 {
		paramOut.Sort = &paramIn.Sort
	}

	if paramIn.Remarks != "" {
		paramOut.Remarks = &paramIn.Remarks
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (dictionaryType) CreateInBatches(paramIn []dto.DictionaryTypeCreate) response.Common {
	var paramOut []model.DictionaryType
	for i := range paramIn {
		var record model.DictionaryType

		if paramIn[i].Creator > 0 {
			record.Creator = &paramIn[i].Creator
		}

		if paramIn[i].LastModifier > 0 {
			record.LastModifier = &paramIn[i].LastModifier
		}

		record.Name = paramIn[i].Name

		if paramIn[i].Sort != 0 {
			record.Sort = &paramIn[i].Sort
		}

		if paramIn[i].Remarks != "" {
			record.Remarks = &paramIn[i].Remarks
		}

		paramOut = append(paramOut, record)
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (dictionaryType) Update(paramIn dto.DictionaryTypeUpdate) response.Common {
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

	if paramIn.Sort != nil {
		if *paramIn.Sort > 0 {
			paramOut["sort"] = paramIn.Sort
		} else if *paramIn.Sort == 0 {
			paramOut["sort"] = nil
		} else {
			return response.Fail(util.ErrorInvalidJSONParameters)
		}
	}

	if paramIn.Remarks != nil {
		if *paramIn.Remarks != "" {
			paramOut["remarks"] = paramIn.Remarks
		} else {
			paramOut["remarks"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.DictionaryType{}).Where("id = ?", paramIn.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (dictionaryType) Delete(paramIn dto.DictionaryTypeDelete) response.Common {
	//由于删除需要做两件事：软删除+记录删除人，所以需要用事务
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//这里记录删除人，在事务中必须放在前面
		//如果放后面，由于是软删除，系统会找不到这条记录，导致无法更新
		err := tx.Debug().Model(&model.DictionaryType{}).Where("id = ?", paramIn.ID).
			Update("deleter", paramIn.Deleter).Error
		if err != nil {
			return err
		}
		//这里删除记录
		err = tx.Delete(&model.DictionaryType{}, paramIn.ID).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (dictionaryType) GetArray(paramIn dto.DictionaryTypeList) response.Common {
	db := global.DB.Model(&model.DictionaryType{})
	// 顺序：where -> count -> Order -> limit -> offset -> array

	//where
	if paramIn.NameInclude != "" {
		db = db.Where("name like ?", "%"+paramIn.NameInclude+"%")
	}

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
		exists := util.FieldIsInModel(model.DictionaryType{}, orderBy)
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

func (dictionaryType) GetList(paramIn dto.DictionaryTypeList) response.List {
	db := global.DB.Model(&model.DictionaryType{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if paramIn.NameInclude != "" {
		db = db.Where("name like ?", "%"+paramIn.NameInclude+"%")
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
		exists := util.FieldIsInModel(model.DictionaryType{}, orderBy)
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
	var data []dto.DictionaryTypeOutput
	db.Model(&model.DictionaryType{}).Find(&data)

	if len(data) == 0 {
		return response.FailForList(util.ErrorRecordNotFound)
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetTotalNumberOfPages(numberOfRecords, pageSize)

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
