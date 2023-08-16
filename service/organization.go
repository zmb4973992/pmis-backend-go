package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type OrganizationGet struct {
	ID int64
}

type OrganizationCreate struct {
	UserID     int64
	SuperiorID int64  `json:"superior_id" binding:"required"` //上级机构ID
	Name       string `json:"name" binding:"required"`        //名称
	//Sort           int    `json:"sort" binding:"required"`       //级别，如公司、事业部、部门等
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type OrganizationUpdate struct {
	LastModifier int64
	ID           int64
	Name         *string `json:"name"`        //名称
	SuperiorID   *int64  `json:"superior_id"` //上级机构ID
}

type OrganizationDelete struct {
	ID int64
}

type OrganizationGetList struct {
	list.Input
	UserID      int64  `json:"-"`
	Name        string `json:"name,omitempty"`
	IsValid     *bool  `json:"is_valid"`
	NameInclude string `json:"name_include,omitempty"`
}

type OrganizationOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`
	Name         string `json:"name"`        //名称
	SuperiorID   *int   `json:"superior_id"` //上级机构id
	IsValid      *bool  `json:"is_valid"`    //是否有效
}

func (d *OrganizationGet) Get() response.Common {
	var result OrganizationOutput

	err := global.DB.Model(model.Organization{}).
		Where("id = ?", d.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	return response.SuccessWithData(result)
}

func (d *OrganizationCreate) Create() response.Common {
	var paramOut model.Organization

	if d.UserID > 0 {
		paramOut.Creator = &d.UserID
	}

	paramOut.Name = d.Name

	paramOut.SuperiorID = &d.SuperiorID

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (d *OrganizationUpdate) Update() response.Common {
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

	if d.SuperiorID != nil {
		if *d.SuperiorID > 0 {
			paramOut["superior_id"] = d.SuperiorID
		} else if *d.SuperiorID == -1 {
			paramOut["superior_id"] = nil
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "UserID",
		"UserID", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Organization{}).
		Where("id = ?", d.ID).Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (d *OrganizationDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.Organization
	global.DB.Where("id = ?", d.ID).Find(&record)
	err := global.DB.Where("id = ?", d.ID).Delete(&record).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (o *OrganizationGetList) GetList() response.List {
	db := global.DB.Model(&model.Organization{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	organizationIDs := util.GetOrganizationIDsForDataAuthority(o.UserID)
	db = db.Where("id in ?", organizationIDs)

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
	if o.PagingInput.Page > 0 {
		page = o.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if o.PagingInput.PageSize != nil && *o.PagingInput.PageSize >= 0 &&
		*o.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *o.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []OrganizationOutput
	db.Model(&model.Organization{}).Find(&data)

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
