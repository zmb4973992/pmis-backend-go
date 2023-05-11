package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type DictionaryTypeGet struct {
	ID int64
}

type DictionaryTypeCreate struct {
	Creator      int64
	LastModifier int64
	Name         string `json:"name" binding:"required"` //名称
	Sort         int    `json:"sort,omitempty"`          //顺序值
	Status       *bool  `json:"status"`                  //是否启用
	Remarks      string `json:"remarks,omitempty"`       //备注
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DictionaryTypeUpdate struct {
	LastModifier int64
	ID           int64
	Name         *string `json:"name"`    //名称
	Sort         *int    `json:"sort"`    //顺序值
	Status       *bool   `json:"status"`  //是否启用
	Remarks      *string `json:"remarks"` //备注
}

type DictionaryTypeDelete struct {
	ID int64
}

type DictionaryTypeGetArray struct {
	list.Input
	NameInclude string `json:"name_include,omitempty"`
}

type DictionaryTypeGetList struct {
	list.Input
	NameInclude string `json:"name_include,omitempty"`
}

type DictionaryTypeOutput struct {
	Creator      *int64  `json:"creator"`
	LastModifier *int64  `json:"last_modifier"`
	ID           int64   `json:"id"`
	Name         string  `json:"name"`    //名称
	Sort         *int    `json:"sort"`    //顺序值
	Remarks      *string `json:"remarks"` //备注
}

func (d *DictionaryTypeGet) Get() response.Common {
	var result DictionaryTypeOutput
	err := global.DB.Model(model.DictionaryType{}).
		Where("id = ?", d.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (d *DictionaryTypeCreate) Create() response.Common {
	var paramOut model.DictionaryType
	if d.Creator > 0 {
		paramOut.Creator = &d.Creator
	}

	if d.LastModifier > 0 {
		paramOut.LastModifier = &d.LastModifier
	}

	paramOut.Name = d.Name

	if d.Sort != 0 {
		paramOut.Sort = &d.Sort
	}

	if d.Remarks != "" {
		paramOut.Remarks = &d.Remarks
	}

	if d.Status != nil {
		paramOut.Status = d.Status
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (d *DictionaryTypeUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if d.LastModifier > 0 {
		paramOut["last_modifier"] = d.LastModifier
	}

	if d.Name != nil {
		if *d.Name != "" {
			paramOut["name"] = d.Name
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if d.Sort != nil {
		if *d.Sort > 0 {
			paramOut["sort"] = d.Sort
		} else if *d.Sort == 0 {
			paramOut["sort"] = nil
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if d.Status != nil {
		paramOut["status"] = d.Status
	}

	if d.Remarks != nil {
		if *d.Remarks != "" {
			paramOut["remarks"] = d.Remarks
		} else {
			paramOut["remarks"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt", "ID")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.DictionaryType{}).Where("id = ?", d.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (d *DictionaryTypeDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.DictionaryType
	global.DB.Where("id = ?", d.ID).Find(&record)
	err := global.DB.Where("id = ?", d.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (d *DictionaryTypeGetList) GetList() response.List {
	db := global.DB.Model(&model.DictionaryType{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if d.NameInclude != "" {
		db = db.Where("name like ?", "%"+d.NameInclude+"%")
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := d.SortingInput.OrderBy
	desc := d.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.DictionaryType{}, orderBy)
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
	if d.PagingInput.Page > 0 {
		page = d.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if d.PagingInput.PageSize != nil && *d.PagingInput.PageSize >= 0 &&
		*d.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *d.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []DictionaryTypeOutput
	db.Model(&model.DictionaryType{}).Find(&data)

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
