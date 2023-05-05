package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type ErrorLogGet struct {
	SnowID int64
}

type ErrorLogCreate struct {
	Creator      int64
	LastModifier int64

	Detail        string `json:"detail,omitempty" `
	Date          string `json:"date,omitempty"`
	MajorCategory string `json:"major_category,omitempty"`
	MinorCategory string `json:"minor_category,omitempty"`
	IsResolved    bool   `json:"is_resolved,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ErrorLogUpdate struct {
	LastModifier int64
	SnowID       int64

	Detail        *string `json:"detail"`
	Date          *string `json:"date"`
	MajorCategory *string `json:"major_category"`
	MinorCategory *string `json:"minor_category"`
	IsResolved    *bool   `json:"is_resolved"`
}

type ErrorLogDelete struct {
	SnowID int64
}

type ErrorLogGetList struct {
	list.Input

	DetailInclude string `json:"detail_include,omitempty" `
	Date          string `json:"date,omitempty"`
	MajorCategory string `json:"major_category,omitempty"`
	MinorCategory string `json:"minor_category,omitempty"`
	IsResolved    bool   `json:"is_resolved,omitempty"`
}

//以下为出参

type ErrorLogOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	SnowID       int64  `json:"snow_id"`

	Detail        *string `json:"detail"`
	Date          *string `json:"date"`
	MajorCategory *string `json:"major_category"`
	MinorCategory *string `json:"minor_category"`
	IsResolved    *bool   `json:"is_resolved"`
}

func (e *ErrorLogGet) Get() response.Common {
	var result ErrorLogOutput
	err := global.DB.Model(model.ErrorLog{}).
		Where("snow_id = ?", e.SnowID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	if result.Date != nil {
		date := *result.Date
		*result.Date = date[:10]
	}
	return response.SuccessWithData(result)
}

func (e *ErrorLogCreate) Create() response.Common {
	var paramOut model.ErrorLog
	if e.Creator > 0 {
		paramOut.Creator = &e.Creator
	}

	if e.LastModifier > 0 {
		paramOut.LastModifier = &e.LastModifier
	}

	if e.Detail != "" {
		paramOut.Detail = &e.Detail
	}

	if e.Date != "" {
		date, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Failure(util.ErrorInvalidJSONParameters)
		} else {
			paramOut.Date = &date
		}
	}

	if e.MajorCategory != "" {
		paramOut.MajorCategory = &e.MajorCategory
	}

	if e.MinorCategory != "" {
		paramOut.MinorCategory = &e.MinorCategory
	}

	if e.IsResolved != false {
		paramOut.IsResolved = &e.IsResolved
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (e *ErrorLogUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if e.LastModifier > 0 {
		paramOut["last_modifier"] = e.LastModifier
	}

	if e.Detail != nil {
		if *e.Detail != "" {
			paramOut["detail"] = e.Detail
		} else {
			paramOut["detail"] = nil
		}
	}

	if e.Date != nil {
		if *e.Date != "" {
			date, err := time.Parse("2006-01-02", *e.Date)
			if err != nil {
				return response.Failure(util.ErrorInvalidJSONParameters)
			}
			paramOut["date"] = date
		} else {
			paramOut["date"] = nil
		}
	}

	if e.MajorCategory != nil {
		if *e.MajorCategory != "" {
			paramOut["major_category"] = e.MajorCategory
		} else {
			paramOut["major_category"] = nil
		}
	}

	if e.MinorCategory != nil {
		if *e.MinorCategory != "" {
			paramOut["minor_category"] = e.MinorCategory
		} else {
			paramOut["minor_category"] = nil
		}
	}

	if e.IsResolved != nil {
		if *e.IsResolved != false {
			paramOut["is_resolved"] = e.IsResolved
		} else {
			paramOut["is_resolved"] = nil
		}
	}
	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.ErrorLog{}).Where("snow_id = ?", e.SnowID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (e *ErrorLogDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.ErrorLog
	global.DB.Where("snow_id = ?", e.SnowID).Find(&record)
	err := global.DB.Where("snow_id = ?", e.SnowID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

// GetList date待修改
func (e *ErrorLogGetList) GetList() response.List {
	db := global.DB.Model(&model.ErrorLog{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if e.DetailInclude != "" {
		db = db.Where("detail like ?", "%"+e.DetailInclude+"%")
	}

	//待完成
	if e.Date != "" {

	}

	if e.MajorCategory != "" {
		db = db.Where("major_category = ?", e.MajorCategory)
	}

	if e.MinorCategory != "" {
		db = db.Where("minor_category = ?", e.MinorCategory)
	}

	if e.IsResolved != false {
		db = db.Where("is_resolved = ?", e.IsResolved)
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := e.SortingInput.OrderBy
	desc := e.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("snow_id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.ErrorLog{}, orderBy)
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
	if e.PagingInput.Page > 0 {
		page = e.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if e.PagingInput.PageSize != nil && *e.PagingInput.PageSize >= 0 &&
		*e.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *e.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []ErrorLogOutput
	db.Model(&model.ErrorLog{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return response.List{
		Data: data,
		Paging: &list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
