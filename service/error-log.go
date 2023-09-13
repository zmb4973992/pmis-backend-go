package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type ErrorLogGet struct {
	Id int64
}

type ErrorLogCreate struct {
	UserId int64

	Detail            string `json:"detail,omitempty"`
	MainCategory      string `json:"main_category,omitempty"`
	SecondaryCategory string `json:"secondary_category,omitempty"`
	IsResolved        bool   `json:"is_resolved,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type ErrorLogUpdate struct {
	UserId int64
	Id     int64

	Detail            *string `json:"detail"`
	MainCategory      *string `json:"main_category"`
	SecondaryCategory *string `json:"secondary_category"`
	IsResolved        *bool   `json:"is_resolved"`
}

type ErrorLogDelete struct {
	Id int64
}

type ErrorLogGetList struct {
	list.Input

	DetailInclude     string `json:"detail_include,omitempty" `
	MainCategory      string `json:"main_category,omitempty"`
	SecondaryCategory string `json:"secondary_category,omitempty"`
	IsResolved        bool   `json:"is_resolved,omitempty"`
}

//以下为出参

type ErrorLogOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`

	Detail            *string `json:"detail"`
	Date              *string `json:"date"`
	MainCategory      *string `json:"main_category"`
	SecondaryCategory *string `json:"secondary_category"`
	IsResolved        *bool   `json:"is_resolved"`
}

func (e *ErrorLogGet) Get() (output *ErrorLogOutput, errCode int) {
	err := global.DB.Model(model.ErrorLog{}).
		Where("id = ?", e.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	return output, util.Success
}

func (e *ErrorLogCreate) Create() (errCode int) {
	var paramOut model.ErrorLog
	if e.UserId > 0 {
		paramOut.Creator = &e.UserId
	}

	if e.Detail != "" {
		paramOut.Detail = &e.Detail
	}

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
		return util.ErrorFailToCreateRecord
	}

	return util.Success
}

func (e *ErrorLogUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	if e.UserId > 0 {
		paramOut["last_modifier"] = e.UserId
	}

	if e.Detail != nil {
		if *e.Detail != "" {
			paramOut["detail"] = e.Detail
		} else {
			paramOut["detail"] = nil
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

	err := global.DB.Model(&model.ErrorLog{}).
		Where("id = ?", e.Id).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}

func (e *ErrorLogDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.ErrorLog
	err := global.DB.Where("id = ?", e.Id).
		Find(&record).
		Delete(&record).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}
	return util.Success
}

func (e *ErrorLogGetList) GetList() (
	outputs []model.ErrorLog, errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.ErrorLog{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if e.DetailInclude != "" {
		db = db.Where("detail like ?", "%"+e.DetailInclude+"%")
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
			return nil, util.ErrorSortingFieldDoesNotExist, nil
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
	pageSize := global.Config.Paging.DefaultPageSize
	if e.PagingInput.PageSize != nil && *e.PagingInput.PageSize >= 0 &&
		*e.PagingInput.PageSize <= global.Config.Paging.MaxPageSize {
		pageSize = *e.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.ErrorLog{}).Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return outputs,
		util.Success,
		&list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		}
}
