package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

type errorLog struct{}

func (*errorLog) Get(errorLogID int) response.Common {
	var result dto.ErrorLogOutput
	err := global.DB.Model(model.ErrorLog{}).
		Where("id = ?", errorLogID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}
	if result.Date != nil {
		date := *result.Date
		*result.Date = date[:10]
	}
	return response.SucceedWithData(result)
}

func (*errorLog) Create(paramIn dto.ErrorLogCreate) response.Common {
	var paramOut model.ErrorLog
	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	if paramIn.Detail != "" {
		paramOut.Detail = &paramIn.Detail
	}

	if paramIn.Date != "" {
		date, err := time.Parse("2006-01-02", paramIn.Date)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Fail(util.ErrorInvalidJSONParameters)
		} else {
			paramOut.Date = &date
		}
	}

	if paramIn.MajorCategory != "" {
		paramOut.MajorCategory = &paramIn.MajorCategory
	}

	if paramIn.MinorCategory != "" {
		paramOut.MinorCategory = &paramIn.MinorCategory
	}

	if paramIn.IsResolved != false {
		paramOut.IsResolved = &paramIn.IsResolved
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (*errorLog) Update(paramIn dto.ErrorLogUpdate) response.Common {
	paramOut := make(map[string]any)

	if paramIn.LastModifier > 0 {
		paramOut["last_modifier"] = paramIn.LastModifier
	}

	if paramIn.Detail != nil {
		if *paramIn.Detail != "" {
			paramOut["detail"] = paramIn.Detail
		} else {
			paramOut["detail"] = nil
		}
	}

	if paramIn.Date != nil {
		if *paramIn.Date != "" {
			date, err := time.Parse("2006-01-02", *paramIn.Date)
			if err != nil {
				return response.Fail(util.ErrorInvalidJSONParameters)
			}
			paramOut["date"] = date
		} else {
			paramOut["date"] = nil
		}
	}

	if paramIn.MajorCategory != nil {
		if *paramIn.MajorCategory != "" {
			paramOut["major_category"] = paramIn.MajorCategory
		} else {
			paramOut["major_category"] = nil
		}
	}

	if paramIn.MinorCategory != nil {
		if *paramIn.MinorCategory != "" {
			paramOut["minor-category"] = paramIn.MinorCategory
		} else {
			paramOut["minor-category"] = nil
		}
	}

	if paramIn.IsResolved != nil {
		if *paramIn.IsResolved != false {
			paramOut["is_resolved"] = paramIn.IsResolved
		} else {
			paramOut["is_resolved"] = nil
		}
	}
	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.ErrorLog{}).Where("id = ?", paramIn.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (*errorLog) Delete(paramIn dto.ErrorLogDelete) response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.ErrorLog
	global.DB.Where("id = ?", paramIn.ID).Find(&record)
	record.Deleter = &paramIn.Deleter
	err := global.DB.Where("id = ?", paramIn.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

// GetList date待修改
func (*errorLog) GetList(paramIn dto.ErrorLogList) response.List {
	db := global.DB.Model(&model.ErrorLog{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if paramIn.DetailInclude != "" {
		db = db.Where("detail like ?", "%"+paramIn.DetailInclude+"%")
	}

	//待完成
	if paramIn.Date != "" {

	}

	if paramIn.MajorCategory != "" {
		db = db.Where("major_category = ?", paramIn.MajorCategory)
	}

	if paramIn.MinorCategory != "" {
		db = db.Where("minor_category = ?", paramIn.MinorCategory)
	}

	if paramIn.IsResolved != false {
		db = db.Where("is_resolved = ?", paramIn.IsResolved)
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
		exists := util.FieldIsInModel(&model.ErrorLog{}, orderBy)
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
	var data []dto.ErrorLogOutput
	db.Model(&model.ErrorLog{}).Find(&data)

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
