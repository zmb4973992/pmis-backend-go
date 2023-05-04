package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type DictionaryDetailGet struct {
	SnowID int64
}

type DictionaryDetailCreate struct {
	Creator              int64
	LastModifier         int64
	DictionaryTypeSnowID int64  `json:"dictionary_type_snow_id" binding:"required,gt=0"` //字典类型id
	Name                 string `json:"name" binding:"required"`                         //名称
	Sequence             int    `json:"sequence,omitempty"`                              //顺序值
	Remarks              string `json:"remarks,omitempty"`                               //备注
}

type DictionaryDetailCreateInBatches struct {
	Data []DictionaryDetailCreate `json:"data"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DictionaryDetailUpdate struct {
	LastModifier int64
	SnowID       int64
	Name         *string `json:"name"`     //名称
	Sequence     *int    `json:"sequence"` //顺序值
	Remarks      *string `json:"remarks"`  //备注
}

type DictionaryDetailDelete struct {
	SnowID int64
}

type DictionaryDetailGetArray struct {
	ListInput
	DictionaryTypeSnowID int64  `json:"dictionary_type_snow_id,omitempty"`
	DictionaryTypeName   string `json:"dictionary_type_name,omitempty"`
}

type DictionaryDetailGetList struct {
	ListInput
	DictionaryTypeSnowID int64  `json:"dictionary_type_snow_id,omitempty"`
	DictionaryTypeName   string `json:"dictionary_type_name,omitempty"`
}

//以下为出参

type DictionaryDetailOutput struct {
	Creator              *int64  `json:"creator"`
	LastModifier         *int64  `json:"last_modifier"`
	SnowID               int64   `json:"snow_id"`
	DictionaryTypeSnowID int64   `json:"dictionary_type_snow_id"` //字典类型
	Name                 string  `json:"name"`                    //名称
	Sequence             *int    `json:"sequence"`                //顺序值
	Remarks              *string `json:"remarks"`                 //备注
}

func (d *DictionaryDetailGet) Get() response.Common {
	var result DictionaryDetailOutput
	err := global.DB.Model(model.DictionaryDetail{}).
		Where("id = ?", d.SnowID).First(&result).Error
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

	paramOut.DictionaryTypeSnowID = d.DictionaryTypeSnowID

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
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (d *DictionaryDetailCreateInBatches) CreateInBatches() response.Common {
	var paramOut []model.DictionaryDetail
	for i := range d.Data {
		var record model.DictionaryDetail

		if d.Data[i].Creator > 0 {
			record.Creator = &d.Data[i].Creator
		}

		if d.Data[i].LastModifier > 0 {
			record.LastModifier = &d.Data[i].LastModifier
		}

		record.DictionaryTypeSnowID = d.Data[i].DictionaryTypeSnowID

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
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (d *DictionaryDetailUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if d.LastModifier > 0 {
		paramOut["last_modifier"] = d.LastModifier
	}

	//if d.DictionaryTypeSnowID != nil {
	//	if *d.DictionaryTypeSnowID > 0 {
	//		paramOut["dictionary_type_id"] = d.DictionaryTypeSnowID
	//	} else if *d.DictionaryTypeSnowID == 0 {
	//		paramOut["dictionary_type_id"] = nil
	//	} else {
	//		return response.Failure(util.ErrorInvalidJSONParameters)
	//	}
	//}

	if d.Name != nil {
		if *d.Name != "" {
			paramOut["name"] = d.Name
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if d.Sequence != nil {
		if *d.Sequence > 0 {
			paramOut["sequence"] = d.Sequence
		} else if *d.Sequence == 0 {
			paramOut["sequence"] = nil
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
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
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.DictionaryDetail{}).Where("id = ?", d.SnowID).
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
	global.DB.Where("id = ?", d.SnowID).Find(&record)
	err := global.DB.Where("id = ?", d.SnowID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

//func (d *DictionaryDetailGetArray) GetArray() response.Common {
//	db := global.DB.Model(&model.DictionaryDetail{})
//	// 顺序：where -> count -> Order -> limit -> offset -> array
//
//	//where
//	if d.DictionaryTypeSnowID != 0 {
//		db = db.Where("dictionary_type_id = ?", d.DictionaryTypeSnowID)
//	}
//
//	//Order
//	orderBy := d.SortingInput.OrderBy
//	desc := d.SortingInput.Desc
//	//如果排序字段为空
//	if orderBy == "" {
//		//如果要求降序排列
//		if desc == true {
//			db = db.Order("id desc")
//		}
//	} else { //如果有排序字段
//		//先看排序字段是否存在于表中
//		exists := util.FieldIsInModel(&model.DictionaryDetail{}, orderBy)
//		if !exists {
//			return response.Failure(util.ErrorSortingFieldDoesNotExist)
//		}
//		//如果要求降序排列
//		if desc == true {
//			db = db.Order(orderBy + " desc")
//		} else { //如果没有要求排序方式
//			db = db.Order(orderBy)
//		}
//	}
//
//	//limit
//	page := 1
//	if d.PagingInput.Page > 0 {
//		page = d.PagingInput.Page
//	}
//	pageSize := global.Config.DefaultPageSize
//	if d.PagingInput.PageSize != nil && *d.PagingInput.PageSize >= 0 &&
//		*d.PagingInput.PageSize <= global.Config.MaxPageSize {
//		pageSize = *d.PagingInput.PageSize
//	}
//	db = db.Limit(pageSize)
//
//	//offset
//	offset := (page - 1) * pageSize
//	db = db.Offset(offset)
//
//	//array
//	var array []string
//	db.Model(&model.DictionaryDetail{}).Select("name").Find(&array)
//
//	if len(array) == 0 {
//		return response.Failure(util.ErrorRecordNotFound)
//	}
//
//	return response.Common{
//		Data:    array,
//		Code:    util.Success,
//		Message: util.GetMessage(util.Success),
//	}
//}

func (d *DictionaryDetailGetList) GetList() response.List {
	db := global.DB.Model(&model.DictionaryDetail{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if d.DictionaryTypeSnowID != 0 {
		db = db.Where("dictionary_type_id = ?", d.DictionaryTypeSnowID)
	}

	if d.DictionaryTypeName != "" {
		var dictionaryTypeID int
		global.DB.Model(&model.DictionaryType{}).Where("name = ?", d.DictionaryTypeName).
			Select("id").Limit(1).Find(&dictionaryTypeID)
		if dictionaryTypeID > 0 {
			db = db.Where("dictionary_type_id = ?", dictionaryTypeID)
		} else {
			return response.FailureForList(util.ErrorDictionaryTypeNameNotFound)
		}
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
		Paging: &PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
