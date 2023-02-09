package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type DepartmentGet struct {
	ID int
}

type DepartmentCreate struct {
	Creator      int
	LastModifier int
	Name         string `json:"name" binding:"required"`             //名称
	LevelName    string `json:"level_name" binding:"required"`       //级别，如公司、事业部、部门等
	SuperiorID   int    `json:"superior_id" binding:"required,gt=0"` //上级机构ID
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DepartmentUpdate struct {
	LastModifier int
	ID           int
	Name         *string `json:"name"`        //名称
	LevelName    *string `json:"level_name"`  //级别，如公司、事业部、部门等
	SuperiorID   *int    `json:"superior_id"` //上级机构ID
}

type DepartmentDelete struct {
	Deleter int
	ID      int
}

type DepartmentGetArray struct {
	dto.ListInput
	dto.AuthInput
	SuperiorID  int    `json:"superior_id,omitempty"`
	LevelName   string `json:"level_name,omitempty"`
	Name        string `json:"name,omitempty"`
	NameInclude string `json:"name_include,omitempty"`
}

type DepartmentGetList struct {
	dto.ListInput
	dto.AuthInput
	SuperiorID  int    `json:"superior_id,omitempty"`
	LevelName   string `json:"level_name,omitempty"`
	Name        string `json:"name,omitempty"`
	NameInclude string `json:"name_include,omitempty"`
}

type DepartmentOutput struct {
	Creator      *int    `json:"creator" gorm:"creator"`
	LastModifier *int    `json:"last_modifier" gorm:"last_modifier"`
	ID           int     `json:"id" gorm:"id"`
	Name         string  `json:"name" gorm:"name"`               //名称
	LevelName    *string `json:"level_name" gorm:"level_name"`   //级别，如公司、事业部、部门等
	SuperiorID   *int    `json:"superior_id" gorm:"superior_id"` //上级机构id
}

func (d *DepartmentGet) Get() response.Common {
	var result DepartmentOutput

	err := global.DB.Model(model.Department{}).
		Where("id = ?", d.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}

	return response.SucceedWithData(result)
}

func (d *DepartmentCreate) Create() response.Common {
	var paramOut model.Department

	if d.Creator > 0 {
		paramOut.Creator = &d.Creator
	}

	if d.LastModifier > 0 {
		paramOut.LastModifier = &d.LastModifier
	}

	paramOut.Name = d.Name

	paramOut.LevelName = d.LevelName

	paramOut.SuperiorID = &d.SuperiorID

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (d *DepartmentUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if d.LastModifier > 0 {
		paramOut["last_modifier"] = d.LastModifier
	}

	if d.Name != nil {
		if *d.Name != "" {
			paramOut["name"] = d.Name
		} else {
			paramOut["name"] = nil
		}
	}

	if d.LevelName != nil {
		if *d.LevelName != "" {
			paramOut["level_name"] = d.LevelName
		} else {
			paramOut["level_name"] = nil
		}
	}

	if d.SuperiorID != nil {
		if *d.SuperiorID > 0 {
			paramOut["superior_id"] = d.SuperiorID
		} else if *d.SuperiorID == 0 {
			paramOut["superior_id"] = nil
		} else {
			return response.Fail(util.ErrorInvalidJSONParameters)
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Department{}).
		Where("id = ?", d.ID).Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (d *DepartmentDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.Department
	global.DB.Where("id = ?", d.ID).Find(&record)
	record.Deleter = &d.Deleter
	err := global.DB.Where("id = ?", d.ID).Delete(&record).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (d *DepartmentGetArray) GetArray() response.Common {
	db := global.DB.Model(&model.Department{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if d.SuperiorID > 0 {
		db = db.Where("superior_id = ?", d.SuperiorID)
	}

	if d.LevelName != "" {
		db = db.Where("level_name = ?", d.LevelName)
	}

	if d.Name != "" {
		db = db.Where("name = ?", d.Name)
	}

	if d.NameInclude != "" {
		db = db.Where("name like ?", "%"+d.NameInclude+"%")
	}

	if d.IsShowedByRole {
		biggestRoleName := util.GetBiggestRoleName(d.UserID)
		if biggestRoleName == "事业部级" {
			businessDivisionIDs := util.GetBusinessDivisionIDs(d.UserID)
			db = db.Where("superior_id in ?", businessDivisionIDs)
		} else if biggestRoleName == "部门级" || biggestRoleName == "项目级" {
			departmentIDs := util.GetDepartmentIDs(d.UserID)
			db = db.Where("id in ?", departmentIDs)
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
		exists := util.FieldIsInModel(&model.Department{}, orderBy)
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
	db.Model(&model.DictionaryType{}).Select("name").Find(&array)

	if len(array) == 0 {
		return response.Fail(util.ErrorRecordNotFound)
	}

	return response.Common{
		Data:    array,
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}

func (d *DepartmentGetList) GetList() response.List {
	db := global.DB.Model(&model.Department{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if d.SuperiorID > 0 {
		db = db.Where("superior_id = ?", d.SuperiorID)
	}

	if d.LevelName != "" {
		db = db.Where("level_name = ?", d.LevelName)
	}

	if d.Name != "" {
		db = db.Where("name = ?", d.Name)
	}

	if d.NameInclude != "" {
		db = db.Where("name like ?", "%"+d.NameInclude+"%")
	}

	if d.IsShowedByRole {
		biggestRoleName := util.GetBiggestRoleName(d.UserID)
		if biggestRoleName == "事业部级" {
			businessDivisionIDs := util.GetBusinessDivisionIDs(d.UserID)
			db = db.Where("superior_id in ?", businessDivisionIDs)
		} else if biggestRoleName == "部门级" || biggestRoleName == "项目级" {
			departmentIDs := util.GetDepartmentIDs(d.UserID)
			db = db.Where("id in ?", departmentIDs)
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
		exists := util.FieldIsInModel(&model.Department{}, orderBy)
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
	var data []DepartmentOutput
	db.Model(&model.Department{}).Find(&data)

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
