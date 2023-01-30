package service

import (
	"gorm.io/gorm"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type dictionaryItemService struct{}

func (dictionaryItemService) Get(dictionaryItemID int) response.Common {
	var result dto.DictionaryItemOutput
	err := global.DB.Model(model.DictionaryItem{}).
		Where("id = ?", dictionaryItemID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (dictionaryItemService) Create(paramIn dto.DictionaryItemCreate) response.Common {
	var paramOut model.DictionaryItem
	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	paramOut.DictionaryTypeID = paramIn.DictionaryTypeID

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
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (dictionaryItemService) CreateInBatches(paramIn []dto.DictionaryItemCreate) response.Common {
	var paramOut []model.DictionaryItem
	for i := range paramIn {
		var record model.DictionaryItem

		if paramIn[i].Creator > 0 {
			record.Creator = &paramIn[i].Creator
		}

		if paramIn[i].LastModifier > 0 {
			record.LastModifier = &paramIn[i].LastModifier
		}

		record.DictionaryTypeID = paramIn[i].DictionaryTypeID

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
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (dictionaryItemService) Update(paramIn dto.DictionaryItemUpdate) response.Common {
	paramOut := make(map[string]any)

	if paramIn.LastModifier > 0 {
		paramOut["last_modifier"] = paramIn.LastModifier
	}

	if paramIn.DictionaryTypeID != nil {
		if *paramIn.DictionaryTypeID > 0 {
			paramOut["dictionary_type_id"] = paramIn.DictionaryTypeID
		} else if *paramIn.DictionaryTypeID == 0 {
			paramOut["dictionary_type_id"] = nil
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
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
			return response.Failure(util.ErrorInvalidJSONParameters)
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
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.DictionaryItem{}).Where("id = ?", paramIn.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (dictionaryItemService) Delete(paramIn dto.DictionaryItemDelete) response.Common {
	//由于删除需要做两件事：软删除+记录删除人，所以需要用事务
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//这里记录删除人，在事务中必须放在前面
		//如果放后面，由于是软删除，系统会找不到这条记录，导致无法更新
		err := tx.Debug().Model(&model.DictionaryItem{}).Where("id = ?", paramIn.ID).
			Update("deleter", paramIn.Deleter).Error
		if err != nil {
			return err
		}
		//这里删除记录
		err = tx.Delete(&model.DictionaryItem{}, paramIn.ID).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (dictionaryItemService) GetArray(paramIn dto.DictionaryItemList) response.Common {
	db := global.DB.Model(&model.DictionaryItem{})
	// 顺序：where -> count -> order -> limit -> offset -> array

	//where
	if paramIn.DictionaryTypeID != 0 {
		db = db.Where("dictionary_type_id = ?", paramIn.DictionaryTypeID)
	}

	//order
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
		exists := util.FieldIsInModel(model.DictionaryItem{}, orderBy)
		if !exists {
			return response.Failure(util.ErrorSortingFieldDoesNotExist)
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
	db.Model(&model.DictionaryItem{}).Select("name").Find(&array)

	if len(array) == 0 {
		return response.Failure(util.ErrorRecordNotFound)
	}

	return response.Common{
		Data:    array,
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}

func (dictionaryItemService) GetList(paramIn dto.DictionaryItemList) response.List {
	db := global.DB.Model(&model.DictionaryItem{})
	// 顺序：where -> count -> order -> limit -> offset -> data

	//where
	if paramIn.DictionaryTypeID != 0 {
		db = db.Where("dictionary_type_id = ?", paramIn.DictionaryTypeID)
	}

	// count
	var count int64
	db.Count(&count)

	//order
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
		exists := util.FieldIsInModel(model.DictionaryItem{}, orderBy)
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
	var data []dto.DictionaryItemOutput
	db.Model(&model.DictionaryItem{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
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
