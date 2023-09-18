package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type OrganizationGet struct {
	Id int64
}

type OrganizationCreate struct {
	UserId     int64
	SuperiorId int64  `json:"superior_id" binding:"required"` //上级机构id
	Name       string `json:"name" binding:"required"`        //名称
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type OrganizationUpdate struct {
	LastModifier int64
	Id           int64
	Name         *string `json:"name"`        //名称
	SuperiorId   *int64  `json:"superior_id"` //上级机构id
}

type OrganizationDelete struct {
	Id int64
}

type OrganizationGetList struct {
	list.Input
	UserId      int64  `json:"-"`
	Name        string `json:"name,omitempty"`
	IsValid     *bool  `json:"is_valid"`
	NameInclude string `json:"name_include,omitempty"`
}

type OrganizationOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`
	Name         string `json:"name"`        //名称
	SuperiorId   *int   `json:"superior_id"` //上级机构id
	IsValid      *bool  `json:"is_valid"`    //是否有效
}

func (d *OrganizationGet) Get() (output *OrganizationOutput, errCode int) {
	err := global.DB.Model(model.Organization{}).
		Where("id = ?", d.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	return output, util.Success
}

func (d *OrganizationCreate) Create() (errCode int) {
	var paramOut model.Organization

	if d.UserId > 0 {
		paramOut.Creator = &d.UserId
	}

	paramOut.Name = d.Name

	paramOut.SuperiorId = &d.SuperiorId

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}
	return util.Success
}

func (d *OrganizationUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	if d.LastModifier > 0 {
		paramOut["last_modifier"] = d.LastModifier
	}

	if d.Name != nil {
		if *d.Name != "" {
			paramOut["name"] = d.Name
		} else {
			return util.ErrorInvalidJSONParameters
		}
	}

	if d.SuperiorId != nil {
		if *d.SuperiorId > 0 {
			paramOut["superior_id"] = d.SuperiorId
		} else if *d.SuperiorId == -1 {
			paramOut["superior_id"] = nil
		} else {
			return util.ErrorInvalidJSONParameters
		}
	}

	err := global.DB.Model(&model.Organization{}).
		Where("id = ?", d.Id).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}

func (d *OrganizationDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.Organization
	err := global.DB.Where("id = ?", d.Id).
		Find(&record).
		Delete(&record).Error
	if err != nil {
		return util.ErrorFailToDeleteRecord
	}
	return util.Success
}

func (o *OrganizationGetList) GetList() (
	outputs []OrganizationOutput, errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.Organization{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where
	organizationIds := util.GetOrganizationIdsForDataAuthority(o.UserId)
	db = db.Where("id in ?", organizationIds)

	if o.Name != "" {
		db = db.Where("name = ?", o.Name)
	}

	if o.IsValid != nil {
		if *o.IsValid == true {
			db = db.Where("is_valid = ?", true)
		} else if *o.IsValid == false {
			db = db.Where("is_valid = ?", false)
		}
	}

	if o.NameInclude != "" {
		db = db.Where("name like ?", "%"+o.NameInclude+"%")
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := o.SortingInput.OrderBy
	desc := o.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Organization{}, orderBy)
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
	if o.PagingInput.Page > 0 {
		page = o.PagingInput.Page
	}
	pageSize := global.Config.Paging.DefaultPageSize
	if o.PagingInput.PageSize != nil && *o.PagingInput.PageSize >= 0 &&
		*o.PagingInput.PageSize <= global.Config.Paging.MaxPageSize {
		pageSize = *o.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.Organization{}).Find(&outputs)

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
