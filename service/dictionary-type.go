package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
)

type DictionaryTypeGet struct {
	Id int64
}

type DictionaryTypeCreate struct {
	UserId  int64
	Name    string `json:"name" binding:"required"` //名称
	Sort    int    `json:"sort,omitempty"`          //顺序值
	Status  *bool  `json:"status"`                  //是否启用
	Remarks string `json:"remarks,omitempty"`       //备注
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DictionaryTypeUpdate struct {
	UserId  int64
	Id      int64
	Name    *string `json:"name"`    //名称
	Sort    *int    `json:"sort"`    //顺序值
	Status  *bool   `json:"status"`  //是否启用
	Remarks *string `json:"remarks"` //备注
}

type DictionaryTypeDelete struct {
	Id int64
}

type DictionaryTypeGetList struct {
	list.Input
	NameInclude string `json:"name_include,omitempty"`
}

type DictionaryTypeOutput struct {
	Creator      *int64  `json:"creator"`
	LastModifier *int64  `json:"last_modifier"`
	Id           int64   `json:"id"`
	Name         string  `json:"name"`    //名称
	Sort         *int    `json:"sort"`    //顺序值
	Remarks      *string `json:"remarks"` //备注
}

func (d *DictionaryTypeGet) Get() (output *DictionaryTypeOutput, errCode int) {
	err := global.DB.Model(model.DictionaryType{}).
		Where("id = ?", d.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}
	return output, util.Success
}

func (d *DictionaryTypeCreate) Create() (errCode int) {
	var paramOut model.DictionaryType
	if d.UserId > 0 {
		paramOut.Creator = &d.UserId
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
		return util.ErrorFailToCreateRecord
	}
	return util.Success
}

func (d *DictionaryTypeUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	if d.UserId > 0 {
		paramOut["last_modifier"] = d.UserId
	}

	if d.Name != nil {
		if *d.Name != "" {
			paramOut["name"] = d.Name
		} else {
			return util.ErrorInvalidJSONParameters
		}
	}

	if d.Sort != nil {
		if *d.Sort > 0 {
			paramOut["sort"] = d.Sort
		} else if *d.Sort == 0 {
			paramOut["sort"] = nil
		} else {
			return util.ErrorInvalidJSONParameters
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

	err := global.DB.Model(&model.DictionaryType{}).
		Where("id = ?", d.Id).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}

func (d *DictionaryTypeDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.DictionaryType
	err := global.DB.Where("id = ?", d.Id).
		Find(&record).
		Delete(&record).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}
	return util.Success
}

func (d *DictionaryTypeGetList) GetList() (
	outputs []DictionaryTypeOutput, errCode int, paging *list.PagingOutput) {
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

	//outputs
	db.Model(&model.DictionaryType{}).Find(&outputs)

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
