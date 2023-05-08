package service

import (
	"github.com/yitter/idgenerator-go/idgen"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type RoleGet struct {
	SnowID int64
}

type RoleCreate struct {
	Creator      int64
	LastModifier int64
	//连接关联表的id

	//连接dictionary_item表的id

	//日期

	//数字(允许为0、nil)
	SuperiorSnowID *int64 `json:"superior_snow_id"`

	Name string `json:"name" binding:"required"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type RoleUpdate struct {
	LastModifier int64
	SnowID       int64
	//连接关联表的id

	//连接dictionary_item表的id

	//日期

	//允许为0的数字
	SuperiorSnowID *int64 `json:"superior_snow_id"`

	//允许为null的字符串
	Name *string `json:"name"`
}

type RoleDelete struct {
	SnowID int64
}

type RoleGetList struct {
	list.Input
}

//以下为出参

type RoleOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	SnowID       int64  `json:"snow_id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示

	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示

	//关联表的详情，不需要gorm查询，需要在json中显示

	//dictionary_item表的详情，不需要gorm查询，需要在json中显示

	//其他属性
	Name           *string `json:"name"`
	SuperiorSnowID *int64  `json:"superior_snow_id"`
}

func (r *RoleGet) Get() response.Common {
	var result RoleOutput
	err := global.DB.Model(model.Role{}).
		Where("snow_id = ?", r.SnowID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	return response.SuccessWithData(result)
}

func (r *RoleCreate) Create() response.Common {
	var paramOut model.Role

	if r.Creator > 0 {
		paramOut.Creator = &r.Creator
	}
	if r.LastModifier > 0 {
		paramOut.LastModifier = &r.LastModifier
	}

	paramOut.SnowID = idgen.NextId()

	//允许为0的数字
	{
		if r.SuperiorSnowID != nil {
			paramOut.SuperiorSnowID = r.SuperiorSnowID
		}
	}

	//允许为null的字符串
	{
		if r.Name != "" {
			paramOut.Name = r.Name
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "CreateAt", "UpdatedAt", "SnowID")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (r *RoleUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if r.LastModifier > 0 {
		paramOut["last_modifier"] = r.LastModifier
	}

	//允许为0的数字
	{
		if r.SuperiorSnowID != nil {
			if *r.SuperiorSnowID != -1 {
				paramOut["superior_snow_id"] = r.SuperiorSnowID
			} else {
				paramOut["superior_snow_id"] = nil
			}
		}
	}

	//允许为null的字符串
	{
		if r.Name != nil {
			if *r.Name != "" {
				paramOut["name"] = r.Name
			} else {
				return response.Failure(util.ErrorInvalidJSONParameters)
			}
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Role{}).Where("snow_id = ?", r.SnowID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (r *RoleDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.Role
	global.DB.Where("snow_id = ?", r.SnowID).Find(&record)
	err := global.DB.Where("snow_id = ?", r.SnowID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (r *RoleGetList) GetList() response.List {
	db := global.DB.Model(&model.Role{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where

	//count
	var count int64
	db.Count(&count)

	//order
	orderBy := r.SortingInput.OrderBy
	desc := r.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("snow_id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Role{}, orderBy)
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
	if r.PagingInput.Page > 0 {
		page = r.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if r.PagingInput.PageSize != nil && *r.PagingInput.PageSize >= 0 &&
		*r.PagingInput.PageSize <= global.Config.MaxPageSize {

		pageSize = *r.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []RoleOutput
	db.Model(&model.Role{}).Find(&data)

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
