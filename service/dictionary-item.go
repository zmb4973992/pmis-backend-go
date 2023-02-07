package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type DictionaryItemGet struct {
	ID int
}

type DictionaryItemCreate struct {
	Creator          int
	LastModifier     int
	DictionaryTypeID int    `json:"dictionary_type_id" binding:"required,gt=0"` //字典类型id
	Name             string `json:"name" binding:"required"`                    //名称
	Sequence         int    `json:"sequence,omitempty"`                         //顺序值
	Remarks          string `json:"remarks,omitempty"`                          //备注
}

type DictionaryItemCreateInBatches struct {
	Data []DictionaryItemCreate `json:"data"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DictionaryItemUpdate struct {
	LastModifier     int
	ID               int
	DictionaryTypeID *int    `json:"dictionary_type_id"` //字典类型id
	Name             *string `json:"name"`               //名称
	Sequence         *int    `json:"sequence"`           //顺序值
	Remarks          *string `json:"remarks"`            //备注
}

type DictionaryItemDelete struct {
	Deleter int
	ID      int
}

type DictionaryItemGetArray struct {
	dto.ListInput
	DictionaryTypeID int `json:"dictionary_type_id,omitempty"`
}

type DictionaryItemGetList struct {
	dto.ListInput
	DictionaryTypeID int `json:"dictionary_type_id,omitempty"`
}

//以下为出参

type DictionaryItemOutput struct {
	Creator          *int    `json:"creator" gorm:"creator"`
	LastModifier     *int    `json:"last_modifier" gorm:"last_modifier"`
	ID               int     `json:"id" gorm:"id"`
	DictionaryTypeID int     `json:"dictionary_type_id" gorm:"dictionary_type_id"` //字典类型id
	Name             string  `json:"name" gorm:"name"`                             //名称
	Sequence         *int    `json:"sequence" gorm:"sequence"`                     //顺序值
	Remarks          *string `json:"remarks" gorm:"remarks"`                       //备注
}

func (d *DictionaryItemGet) Get() response.Common {
	var result DictionaryItemOutput
	err := global.DB.Model(model.DictionaryItem{}).
		Where("id = ?", d.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}
	return response.SucceedWithData(result)
}

func (d *DictionaryItemCreate) Create() response.Common {
	var paramOut model.DictionaryItem
	if d.Creator > 0 {
		paramOut.Creator = &d.Creator
	}

	if d.LastModifier > 0 {
		paramOut.LastModifier = &d.LastModifier
	}

	paramOut.DictionaryTypeID = d.DictionaryTypeID

	paramOut.Name = d.Name

	if d.Sequence != 0 {
		paramOut.Sequence = &d.Sequence
	}

	if d.Remarks != "" {
		paramOut.Remarks = &d.Remarks
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (d *DictionaryItemCreateInBatches) CreateInBatches() response.Common {
	var paramOut []model.DictionaryItem
	for i := range d.Data {
		var record model.DictionaryItem

		if d.Data[i].Creator > 0 {
			record.Creator = &d.Data[i].Creator
		}

		if d.Data[i].LastModifier > 0 {
			record.LastModifier = &d.Data[i].LastModifier
		}

		record.DictionaryTypeID = d.Data[i].DictionaryTypeID

		record.Name = d.Data[i].Name

		if d.Data[i].Sequence != 0 {
			record.Sequence = &d.Data[i].Sequence
		}

		if d.Data[i].Remarks != "" {
			record.Remarks = &d.Data[i].Remarks
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

func (d *DictionaryItemUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if d.LastModifier > 0 {
		paramOut["last_modifier"] = d.LastModifier
	}

	if d.DictionaryTypeID != nil {
		if *d.DictionaryTypeID > 0 {
			paramOut["dictionary_type_id"] = d.DictionaryTypeID
		} else if *d.DictionaryTypeID == 0 {
			paramOut["dictionary_type_id"] = nil
		} else {
			return response.Fail(util.ErrorInvalidJSONParameters)
		}
	}

	if d.Name != nil {
		if *d.Name != "" {
			paramOut["name"] = d.Name
		} else {
			paramOut["name"] = nil
		}
	}

	if d.Sequence != nil {
		if *d.Sequence > 0 {
			paramOut["sequence"] = d.Sequence
		} else if *d.Sequence == 0 {
			paramOut["sequence"] = nil
		} else {
			return response.Fail(util.ErrorInvalidJSONParameters)
		}
	}

	if d.Remarks != nil {
		if *d.Remarks != "" {
			paramOut["remarks"] = d.Remarks
		} else {
			paramOut["remarks"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.DictionaryItem{}).Where("id = ?", d.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (d *DictionaryItemDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.DictionaryItem
	global.DB.Where("id = ?", d.ID).Find(&record)
	record.Deleter = &d.Deleter
	err := global.DB.Where("id = ?", d.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (d *DictionaryItemGetArray) GetArray() response.Common {
	db := global.DB.Model(&model.DictionaryItem{})
	// 顺序：where -> count -> Order -> limit -> offset -> array

	//where
	if d.DictionaryTypeID != 0 {
		db = db.Where("dictionary_type_id = ?", d.DictionaryTypeID)
	}

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
		exists := util.FieldIsInModel(&model.DictionaryItem{}, orderBy)
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
	if d.PagingInput.Page > 0 {
		page = d.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if d.PagingInput.PageSize > 0 &&
		d.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = d.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//array
	var array []string
	db.Model(&model.DictionaryItem{}).Select("name").Find(&array)

	if len(array) == 0 {
		return response.Fail(util.ErrorRecordNotFound)
	}

	return response.Common{
		Data:    array,
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}

func (d *DictionaryItemGetList) GetList() response.List {
	db := global.DB.Model(&model.DictionaryItem{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if d.DictionaryTypeID != 0 {
		db = db.Where("dictionary_type_id = ?", d.DictionaryTypeID)
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
		exists := util.FieldIsInModel(&model.DictionaryItem{}, orderBy)
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
	if d.PagingInput.Page > 0 {
		page = d.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if d.PagingInput.PageSize > 0 &&
		d.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = d.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []DictionaryItemOutput
	db.Model(&model.DictionaryItem{}).Find(&data)

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
