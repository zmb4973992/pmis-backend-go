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
	ID int64
}

type ErrorLogCreate struct {
	UserID int64

	Detail            string `json:"detail,omitempty"`
	MainCategory      string `json:"main_category,omitempty"`
	SecondaryCategory string `json:"secondary_category,omitempty"`
	IsResolved        bool   `json:"is_resolved,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ErrorLogUpdate struct {
	UserID int64
	ID     int64

	Detail            *string `json:"detail"`
	Date              *string `json:"date"`
	MainCategory      *string `json:"main_category"`
	SecondaryCategory *string `json:"secondary_category"`
	IsResolved        *bool   `json:"is_resolved"`
}

type ErrorLogDelete struct {
	ID int64
}

type ErrorLogGetList struct {
	list.Input

	DetailInclude     string `json:"detail_include,omitempty" `
	Date              string `json:"date,omitempty"`
	MainCategory      string `json:"main_category,omitempty"`
	SecondaryCategory string `json:"secondary_category,omitempty"`
	IsResolved        bool   `json:"is_resolved,omitempty"`
}

//以下为出参

type ErrorLogOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`

	Detail            *string `json:"detail"`
	Date              *string `json:"date"`
	MainCategory      *string `json:"main_category"`
	SecondaryCategory *string `json:"secondary_category"`
	IsResolved        *bool   `json:"is_resolved"`
}

func (e *ErrorLogGet) Get() response.Common {
	var result ErrorLogOutput
	err := global.DB.Model(model.ErrorLog{}).
		Where("id = ?", e.ID).First(&result).Error
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
	if e.UserID > 0 {
		paramOut.Creator = &e.UserID
	}

	if e.Detail != "" {
		paramOut.Detail = &e.Detail
	}

	datetime := time.Now()
	paramOut.Datetime = &datetime

	if e.MainCategory != "" {
		paramOut.MainCategory = &e.MainCategory
	}

	if e.SecondaryCategory != "" {
		paramOut.SecondaryCategory = &e.SecondaryCategory
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

	if e.UserID > 0 {
		paramOut["last_modifier"] = e.UserID
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

	if e.MainCategory != nil {
		if *e.MainCategory != "" {
			paramOut["main_category"] = e.MainCategory
		} else {
			paramOut["main_category"] = nil
		}
	}

	if e.SecondaryCategory != nil {
		if *e.SecondaryCategory != "" {
			paramOut["secondary_category"] = e.SecondaryCategory
		} else {
			paramOut["secondary_category"] = nil
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
	paramOutForCounting := util.MapCopy(paramOut, "UserID",
		"UserID", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.ErrorLog{}).Where("id = ?", e.ID).
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
	global.DB.Where("id = ?", e.ID).Find(&record)
	err := global.DB.Where("id = ?", e.ID).Delete(&record).Error

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

	if e.MainCategory != "" {
		db = db.Where("main_category = ?", e.MainCategory)
	}

	if e.SecondaryCategory != "" {
		db = db.Where("secondary_category = ?", e.SecondaryCategory)
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
			db = db.Order("id desc")
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
