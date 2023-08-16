package service

import (
	"gorm.io/gorm"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"strconv"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type RoleGet struct {
	ID int64
}

type RoleCreate struct {
	UserID int64
	//连接关联表的id

	//连接dictionary_item表的id

	//日期

	//数字(允许为0、nil)
	SuperiorID      *int64 `json:"superior_id"`
	Name            string `json:"name" binding:"required"`
	DataAuthorityID int64  `json:"data_authority_id" binding:"required"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type RoleUpdate struct {
	UserID int64
	ID     int64
	//连接关联表的id

	//连接dictionary_item表的id

	//日期

	//允许为0的数字
	//SuperiorID *int64 `json:"superior_id"`

	//允许为null的字符串
	Name            *string `json:"name"`
	DataAuthorityID *int64  `json:"data_authority_id"`
}

type RoleDelete struct {
	ID int64
}

type RoleGetList struct {
	list.Input
}

type RoleUpdateUsers struct {
	UserID int64

	RoleID  int64    `json:"-"`
	UserIDs *[]int64 `json:"user_ids"`
}

type RoleUpdateMenus struct {
	UserID int64

	RoleID  int64    `json:"-"`
	MenuIDs *[]int64 `json:"menu_ids"`
}

//以下为出参

type RoleOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示

	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示

	//关联表的详情，不需要gorm查询，需要在json中显示

	//dictionary_item表的详情，不需要gorm查询，需要在json中显示

	//其他属性
	Name       *string `json:"name"`
	SuperiorID *int64  `json:"superior_id"`
}

func (r *RoleGet) Get() response.Common {
	var result RoleOutput
	err := global.DB.Model(model.Role{}).
		Where("id = ?", r.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	return response.SuccessWithData(result)
}

func (r *RoleCreate) Create() response.Common {
	var paramOut model.Role

	if r.UserID > 0 {
		paramOut.Creator = &r.UserID
	}

	//允许为0的数字
	{
		if r.SuperiorID != nil {
			paramOut.SuperiorID = r.SuperiorID
		}
	}

	//允许为null的字符串
	{
		if r.Name != "" {
			paramOut.Name = r.Name
		}
	}

	paramOut.DataAuthorityID = r.DataAuthorityID

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"UserID", "UserID", "CreateAt", "UpdatedAt")

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

	if r.UserID > 0 {
		paramOut["last_modifier"] = r.UserID
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

	if r.DataAuthorityID != nil {
		if *r.DataAuthorityID == -1 {
			paramOut["data_authority_id"] = nil
		} else {
			paramOut["data_authority_id"] = r.DataAuthorityID
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "UserID",
		"UserID", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Role{}).Where("id = ?", r.ID).
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
	global.DB.Where("id = ?", r.ID).Find(&record)
	err := global.DB.Where("id = ?", r.ID).Delete(&record).Error

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
			db = db.Order("id desc")
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

func (r *RoleUpdateUsers) Update() response.Common {
	if r.UserIDs == nil {
		return response.Failure(util.ErrorInvalidJSONParameters)
	}

	if len(*r.UserIDs) == 0 {
		err := global.DB.Where("role_id = ?", r.RoleID).Delete(&model.UserAndRole{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Failure(util.ErrorFailToDeleteRecord)
		}
		return response.Success()
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//先删掉原始记录
		err := tx.Where("role_id = ?", r.RoleID).Delete(&model.UserAndRole{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToDeleteRecord
		}

		//再增加新的记录
		var paramOut []model.UserAndRole
		for _, userID := range *r.UserIDs {
			var record model.UserAndRole

			if r.UserID > 0 {
				record.LastModifier = &r.UserID
			}

			record.RoleID = r.RoleID
			record.UserID = userID
			paramOut = append(paramOut, record)
		}

		for i := range paramOut {
			//计算有修改值的字段数，分别进行不同处理
			tempParamOut, err := util.StructToMap(paramOut[i])
			if err != nil {
				return ErrorFailToUpdateRecord
			}
			paramOutForCounting := util.MapCopy(tempParamOut,
				"UserID", "UserID", "CreateAt", "UpdatedAt")

			if len(paramOutForCounting) == 0 {
				return ErrorFieldsToBeCreatedNotFound
			}
		}

		err = global.DB.Create(&paramOut).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToCreateRecord
		}

		//更新casbin的rbac分组规则
		var param1 rbacUpdateGroupingPolicyByGroup
		param1.Group = strconv.FormatInt(r.RoleID, 10)
		for _, userID := range *r.UserIDs {
			param1.Members = append(param1.Members, strconv.FormatInt(userID, 10))
		}
		err = param1.Update()
		if err != nil {
			return ErrorFailToUpdateRBACGroupingPolicies
		}

		return nil
	})

	switch err {
	case nil:
		return response.Success()
	case ErrorFailToCreateRecord:
		return response.Failure(util.ErrorFailToCreateRecord)
	case ErrorFailToDeleteRecord:
		return response.Failure(util.ErrorFailToDeleteRecord)
	case ErrorFieldsToBeCreatedNotFound:
		return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
	case ErrorFailToUpdateRBACGroupingPolicies:
		return response.Failure(util.ErrorFailToUpdateRBACGroupingPolicies)
	default:
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
}

func (r *RoleUpdateMenus) Update() response.Common {
	if r.MenuIDs == nil {
		return response.Failure(util.ErrorInvalidJSONParameters)
	}

	if len(*r.MenuIDs) == 0 {
		err := global.DB.Where("role_id = ?", r.RoleID).Delete(&model.RoleAndMenu{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Failure(util.ErrorFailToDeleteRecord)
		}
		return response.Success()
	}

	//先删掉原始记录
	err := global.DB.Where("role_id = ?", r.RoleID).Delete(&model.RoleAndMenu{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	//再增加新的记录
	var paramOut []model.RoleAndMenu
	for _, menuID := range *r.MenuIDs {
		var record model.RoleAndMenu

		if r.UserID > 0 {
			record.LastModifier = &r.UserID
		}

		record.RoleID = r.RoleID
		record.MenuID = menuID
		paramOut = append(paramOut, record)
	}

	for i := range paramOut {
		//计算有修改值的字段数，分别进行不同处理
		tempParamOut, err1 := util.StructToMap(paramOut[i])
		if err1 != nil {
			return response.Failure(util.ErrorFailToUpdateRecord)
		}
		paramOutForCounting := util.MapCopy(tempParamOut,
			"UserID", "UserID", "CreateAt", "UpdatedAt")

		if len(paramOutForCounting) == 0 {
			return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
		}
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	//更新casbin的rbac的策略
	var param1 rbacUpdatePolicyByRoleID
	param1.RoleID = r.RoleID
	err = param1.Update()
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRBACPoliciesByRoleID)
	}

	return response.Success()
}
