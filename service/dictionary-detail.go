package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type DictionaryDetailGet struct {
	ID int64
}

type DictionaryDetailCreate struct {
	Creator          int64
	LastModifier     int64
	DictionaryTypeID int64  `json:"dictionary_type_id" binding:"required"` //字典类型id
	Name             string `json:"name" binding:"required"`               //名称
	Sort             int    `json:"sort,omitempty"`                        //顺序值
	Status           *bool  `json:"status"`                                //是否启用
	Remarks          string `json:"remarks,omitempty"`                     //备注
}

type DictionaryDetailCreateInBatches struct {
	Data []DictionaryDetailCreate `json:"data"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DictionaryDetailUpdate struct {
	LastModifier int64
	ID           int64
	Name         *string `json:"name"`    //名称
	Sort         *int    `json:"sort"`    //顺序值
	Status       *bool   `json:"status"`  //是否启用
	Remarks      *string `json:"remarks"` //备注
}

type DictionaryDetailDelete struct {
	ID int64
}

type DictionaryDetailGetArray struct {
	list.Input
	DictionaryTypeID   int64  `json:"dictionary_type_id,omitempty"`
	DictionaryTypeName string `json:"dictionary_type_name,omitempty"`
}

type DictionaryDetailGetList struct {
	list.Input
	DictionaryTypeID int64 `json:"dictionary_type_id,omitempty"`
}

//以下为出参

type DictionaryDetailOutput struct {
	Creator          *int64  `json:"creator"`
	LastModifier     *int64  `json:"last_modifier"`
	ID               int64   `json:"id"`
	DictionaryTypeID int64   `json:"dictionary_type_id"` //字典类型
	Name             string  `json:"name"`               //名称
	Sort             *int    `json:"sort"`               //顺序值
	Status           *bool   `json:"status"`             //是否启用
	Remarks          *string `json:"remarks"`            //备注
}

func (d *DictionaryDetailGet) Get() response.Common {
	var result DictionaryDetailOutput
	err := global.DB.Model(model.DictionaryDetail{}).
		Where("id = ?", d.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (d *DictionaryDetailCreate) Create() response.Common {
	var paramOut model.DictionaryDetail
	if d.Creator > 0 {
		paramOut.Creator = &d.Creator
	}

	if d.LastModifier > 0 {
		paramOut.LastModifier = &d.LastModifier
	}

	paramOut.DictionaryTypeID = d.DictionaryTypeID

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

//func (d *DictionaryDetailCreateInBatches) CreateInBatches() response.Common {
//	var paramOut []model.DictionaryDetail
//	for i := range d.Data {
//		var record model.DictionaryDetail
//
//		if d.Data[i].Creator > 0 {
//			record.Creator = &d.Data[i].Creator
//		}
//
//		if d.Data[i].LastModifier > 0 {
//			record.LastModifier = &d.Data[i].LastModifier
//		}
//
//		record.DictionaryTypeID = d.Data[i].DictionaryTypeID
//
//		record.Name = d.Data[i].Name
//
//		if d.Data[i].Sort != 0 {
//			record.Sort = &d.Data[i].Sort
//		}
//
//		if d.Data[i].Remarks != "" {
//			record.Remarks = &d.Data[i].Remarks
//		}
//
//		paramOut = append(paramOut, record)
//	}
//
//	err := global.DB.Create(&paramOut).Error
//	if err != nil {
//		global.SugaredLogger.Errorln(err)
//		return response.Failure(util.ErrorFailToCreateRecord)
//	}
//	return response.Success()
//}

func (d *DictionaryDetailUpdate) Update() response.Common {
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
		"LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.DictionaryDetail{}).Where("id = ?", d.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (d *DictionaryDetailDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.DictionaryDetail
	global.DB.Where("id = ?", d.ID).Find(&record)
	err := global.DB.Where("id = ?", d.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (d *DictionaryDetailGetList) GetList() response.List {
	db := global.DB.Model(&model.DictionaryDetail{})
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
		exists := util.FieldIsInModel(&model.DictionaryDetail{}, orderBy)
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
	var data []DictionaryDetailOutput
	db.Model(&model.DictionaryDetail{}).Find(&data)

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
