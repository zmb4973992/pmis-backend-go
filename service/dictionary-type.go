package service

import (
	"gorm.io/gorm"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type dictionaryTypeService struct{}

func (dictionaryTypeService) Create(paramIn dto.DictionaryTypeCreate) response.Common {
	var paramOut model.DictionaryType
	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	paramOut.Name = paramIn.Name

	if paramIn.Sort != nil && *paramIn.Sort != -1 {
		paramOut.Sort = paramIn.Sort
	}

	if paramIn.Remarks != nil && *paramIn.Remarks != "" {
		paramOut.Remarks = paramIn.Remarks
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (dictionaryTypeService) CreateInBatches(paramIn []dto.DictionaryTypeCreate) response.Common {
	var paramOut []model.DictionaryType
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	for i := range paramIn {
		var record model.DictionaryType

		if paramIn[i].Creator > 0 {
			record.Creator = &paramIn[i].Creator
		}

		if paramIn[i].LastModifier > 0 {
			record.LastModifier = &paramIn[i].LastModifier
		}

		record.Name = paramIn[i].Name

		if paramIn[i].Sort != nil && *paramIn[i].Sort != -1 {
			record.Sort = paramIn[i].Sort
		}

		if paramIn[i].Remarks != nil && *paramIn[i].Remarks != "" {
			record.Remarks = paramIn[i].Remarks
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

func (dictionaryTypeService) Update(paramIn dto.DictionaryTypeUpdate) response.Common {
	paramOut := make(map[string]any)

	if paramIn.LastModifier > 0 {
		paramOut["last_modifier"] = paramIn.LastModifier
	}

	paramOut["name"] = paramIn.Name

	if paramIn.Sort != nil {
		if *paramIn.Sort != -1 {
			paramOut["sort"] = paramIn.Sort
		} else {
			paramOut["sort"] = nil
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

	err := global.DB.Model(&model.DictionaryType{}).Where("id = ?", paramIn.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (dictionaryTypeService) Delete(paramIn dto.DictionaryTypeDelete) response.Common {
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
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (dictionaryTypeService) List(paramIn dto.DictionaryTypeList) response.List {
	db := global.DB
	// 顺序：where -> count -> order -> limit -> offset -> data
	// count
	var count int64
	db.Model(&model.DictionaryType{}).Count(&count)

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
		exists := util.FieldIsInModel(model.DictionaryType{}, orderBy)
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
	var data []string
	db.Model(&model.DictionaryType{}).Select("name").Find(&data)

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
