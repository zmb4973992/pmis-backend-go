package service

import (
	"errors"
	"gorm.io/gorm"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
	"strconv"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type RoleGet struct {
	Id int64
}

type RoleCreate struct {
	UserId int64
	//连接关联表的id

	//连接dictionary_item表的id

	//日期

	//数字(允许为0、nil)
	SuperiorId      *int64 `json:"superior_id"`
	Name            string `json:"name" binding:"required"`
	DataAuthorityId int64  `json:"data_authority_id" binding:"required"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type RoleUpdate struct {
	UserId int64
	Id     int64
	//连接关联表的id

	//连接dictionary_item表的id

	//日期

	//允许为0的数字
	//SuperiorId *int64 `json:"superior_id"`

	//允许为null的字符串
	Name            *string `json:"name"`
	DataAuthorityId *int64  `json:"data_authority_id"`
}

type RoleDelete struct {
	Id int64
}

type RoleGetList struct {
	list.Input
}

type RoleUpdateUsers struct {
	UserId int64

	RoleId  int64    `json:"-"`
	UserIds *[]int64 `json:"user_ids"`
}

type RoleUpdateMenus struct {
	UserId int64

	RoleId  int64    `json:"-"`
	MenuIds *[]int64 `json:"menu_ids"`
}

//以下为出参

type RoleOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示

	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示

	//关联表的详情，不需要gorm查询，需要在json中显示

	//dictionary_item表的详情，不需要gorm查询，需要在json中显示

	//其他属性
	Name       *string `json:"name"`
	SuperiorId *int64  `json:"superior_id"`
}

func (r *RoleGet) Get() (output *RoleOutput, errCode int) {
	err := global.DB.Model(model.Role{}).
		Where("id = ?", r.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	return output, util.Success
}

func (r *RoleCreate) Create() (errCode int) {
	var paramOut model.Role

	if r.UserId > 0 {
		paramOut.Creator = &r.UserId
	}

	//允许为0的数字
	{
		if r.SuperiorId != nil {
			paramOut.SuperiorId = r.SuperiorId
		}
	}

	//允许为null的字符串
	{
		if r.Name != "" {
			paramOut.Name = r.Name
		}
	}

	paramOut.DataAuthorityId = r.DataAuthorityId

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}
	return util.Success
}

func (r *RoleUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	if r.UserId > 0 {
		paramOut["last_modifier"] = r.UserId
	}

	//允许为null的字符串
	{
		if r.Name != nil {
			if *r.Name != "" {
				paramOut["name"] = r.Name
			} else {
				return util.ErrorInvalidJSONParameters
			}
		}
	}

	if r.DataAuthorityId != nil {
		if *r.DataAuthorityId == -1 {
			paramOut["data_authority_id"] = nil
		} else {
			paramOut["data_authority_id"] = r.DataAuthorityId
		}
	}

	err := global.DB.Model(&model.Role{}).
		Where("id = ?", r.Id).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}

func (r *RoleDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.Role
	err := global.DB.Where("id = ?", r.Id).
		Find(&record).
		Delete(&record).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}
	return util.Success
}

func (r *RoleGetList) GetList() (outputs []RoleOutput,
	errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.Role{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

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
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Role{}, orderBy)
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
	if r.PagingInput.Page > 0 {
		page = r.PagingInput.Page
	}
	pageSize := global.Config.Paging.DefaultPageSize
	if r.PagingInput.PageSize != nil && *r.PagingInput.PageSize >= 0 &&
		*r.PagingInput.PageSize <= global.Config.Paging.MaxPageSize {

		pageSize = *r.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.Role{}).Find(&outputs)

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

func (r *RoleUpdateUsers) Update() (errCode int) {
	if r.UserIds == nil {
		return util.ErrorInvalidJSONParameters
	}

	if len(*r.UserIds) == 0 {
		err := global.DB.Where("role_id = ?", r.RoleId).
			Delete(&model.UserAndRole{}).Error
		if err != nil {
			return util.ErrorFailToDeleteRecord
		}
		return util.Success
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//先删掉原始记录
		err := tx.Where("role_id = ?", r.RoleId).
			Delete(&model.UserAndRole{}).Error
		if err != nil {
			return util.GenerateCustomError(util.ErrorFailToDeleteRecord)
		}

		//再增加新的记录
		var paramOut []model.UserAndRole
		for _, userId := range *r.UserIds {
			var record model.UserAndRole

			if r.UserId > 0 {
				record.LastModifier = &r.UserId
			}

			record.RoleId = r.RoleId
			record.UserId = userId
			paramOut = append(paramOut, record)
		}

		err = global.DB.Create(&paramOut).Error
		if err != nil {
			return util.GenerateCustomError(util.ErrorFailToCreateRecord)
		}

		//更新casbin的rbac分组规则
		var param1 rbacUpdateGroupingPolicyByGroup
		param1.Group = strconv.FormatInt(r.RoleId, 10)
		for _, userId := range *r.UserIds {
			param1.Members = append(param1.Members, strconv.FormatInt(userId, 10))
		}
		err = param1.Update()
		if err != nil {
			return util.GenerateCustomError(util.ErrorFailToUpdateRBACGroupingPolicies)
		}

		return nil
	})

	switch {
	case err == nil:
		return util.Success
	case errors.Is(err, ErrorFailToCreateRecord):
		return util.ErrorFailToCreateRecord
	case errors.Is(err, ErrorFailToDeleteRecord):
		return util.ErrorFailToDeleteRecord
	case errors.Is(err, ErrorFieldsToBeCreatedNotFound):
		return util.ErrorFieldsToBeCreatedNotFound
	case errors.Is(err, ErrorFailToUpdateRBACGroupingPolicies):
		return util.ErrorFailToUpdateRBACGroupingPolicies
	default:
		return util.ErrorFailToUpdateRecord
	}
}

func (r *RoleUpdateMenus) Update() (errCode int) {
	if r.MenuIds == nil {
		return util.ErrorInvalidJSONParameters
	}

	if len(*r.MenuIds) == 0 {
		err := global.DB.Where("role_id = ?", r.RoleId).
			Delete(&model.RoleAndMenu{}).Error
		if err != nil {
			return util.ErrorFailToDeleteRecord
		}
		return util.Success
	}

	//先删掉原始记录
	err := global.DB.Where("role_id = ?", r.RoleId).
		Delete(&model.RoleAndMenu{}).Error
	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	//再增加新的记录
	var paramOut []model.RoleAndMenu
	for _, menuId := range *r.MenuIds {
		var record model.RoleAndMenu

		if r.UserId > 0 {
			record.LastModifier = &r.UserId
		}

		record.RoleId = r.RoleId
		record.MenuId = menuId
		paramOut = append(paramOut, record)
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}

	//更新casbin的rbac的策略
	var param1 rbacUpdatePolicyByRoleId
	param1.RoleId = r.RoleId
	err = param1.Update()
	if err != nil {
		return util.ErrorFailToUpdateRBACPoliciesByRoleId
	}

	return util.Success
}
