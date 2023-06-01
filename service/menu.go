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

type MenuGet struct {
	ID int64
}

type MenuCreate struct {
	Creator      int64
	LastModifier int64
	//连接关联表的id

	//连接dictionary_item表的id

	//日期

	//数字(允许为0、nil)
	SuperiorID    int64  `json:"superior_id,omitempty"`
	Path          string `json:"path" binding:"required"`
	Group         string `json:"group"  binding:"required"`
	Name          string `json:"name"  binding:"required"`
	HiddenInSider bool   `json:"hidden_in_sider" `
	Component     string `json:"component" binding:"required"`
	Sort          int    `json:"sort" binding:"required"`
	KeepAlive     bool   `json:"keep_alive" `
	Title         string `json:"title" binding:"required"`
	Icon          string `json:"icon,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type MenuUpdate struct {
	LastModifier int64
	ID           int64
	//连接关联表的id

	//连接dictionary_item表的id

	//日期

	//允许为0的数字
	//SuperiorID *int64 `json:"superior_id"`

	//允许为null的字符串

	SuperiorID    *int64  `json:"superior_id"`
	Path          *string `json:"path"`
	Group         *string `json:"group"`
	Name          *string `json:"name"`
	HiddenInSider *bool   `json:"hidden_in_sider"`
	Component     *string `json:"component"`
	Sort          *int    `json:"sort"`
	KeepAlive     *bool   `json:"keep_alive"`
	Title         *string `json:"title"`
	Icon          *string `json:"icon"`
}

type MenuDelete struct {
	ID int64
}

type MenuGetList struct {
	list.Input
	list.DataScopeInput
	Group string `json:"group,omitempty"`
}

type MenuGetTree struct {
	list.Input
	list.DataScopeInput
	Group string `json:"group,omitempty"`
}

type MenuUpdateApis struct {
	Creator      int64
	LastModifier int64

	MenuID int64    `json:"-"`
	ApiIDs *[]int64 `json:"api_ids"`
}

//以下为出参

type Meta struct {
	Hidden    *bool   `json:"hidden"`
	KeepAlive *bool   `json:"keep_alive"`
	Title     *string `json:"title"`
	Icon      *string `json:"icon"`
}

type MenuOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示

	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示

	//关联表的详情，不需要gorm查询，需要在json中显示

	//dictionary_item表的详情，不需要gorm查询，需要在json中显示

	//其他属性
	SuperiorID *int64  `json:"superior_id"`
	Path       *string `json:"path"`
	Group      *string `json:"group"`
	Name       *string `json:"name"`
	Component  *string `json:"component"`
	Sort       *int    `json:"sort"`
	Meta       `json:"meta"`
	Children   []MenuOutput `json:"children" gorm:"-"`
}

func (m *MenuGet) Get() response.Common {
	var result MenuOutput
	err := global.DB.Model(model.Menu{}).
		Where("id = ?", m.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	return response.SuccessWithData(result)
}

func (m *MenuCreate) Create() response.Common {
	var paramOut model.Menu

	if m.Creator > 0 {
		paramOut.Creator = &m.Creator
	}
	if m.LastModifier > 0 {
		paramOut.LastModifier = &m.LastModifier
	}

	if m.SuperiorID > 0 {
		paramOut.SuperiorID = &m.SuperiorID
	}

	if m.Path != "" {
		paramOut.Path = &m.Path
	}

	if m.Group != "" {
		paramOut.Group = m.Group
	}

	if m.Name != "" {
		paramOut.Name = m.Name
	}

	paramOut.Hidden = m.HiddenInSider

	if m.Component != "" {
		paramOut.Component = &m.Component
	}

	if m.Sort > 0 {
		paramOut.Sort = &m.Sort
	}

	paramOut.KeepAlive = &m.KeepAlive

	if m.Title != "" {
		paramOut.Title = &m.Title
	}

	if m.Icon != "" {
		paramOut.Icon = &m.Icon
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "CreateAt", "UpdatedAt")

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

func (m *MenuUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if m.LastModifier > 0 {
		paramOut["last_modifier"] = m.LastModifier
	}

	//允许为0的数字
	//{
	//	if m.SuperiorID != nil {
	//		if *m.SuperiorID != -1 {
	//			paramOut["superior_id"] = m.SuperiorID
	//		} else {
	//			paramOut["superior_id"] = nil
	//		}
	//	}
	//}

	if m.SuperiorID != nil {
		if *m.SuperiorID == -1 {
			paramOut["superior_id"] = nil
		} else {
			paramOut["superior_id"] = *m.SuperiorID
		}
	}

	if m.Path != nil {
		if *m.Path != "" {
			paramOut["route_path"] = m.Path
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if m.Group != nil {
		if *m.Group != "" {
			paramOut["group"] = m.Group
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	//允许为null的字符串
	if m.Name != nil {
		if *m.Name != "" {
			paramOut["name"] = m.Name
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if m.Component != nil {
		if *m.Component != "" {
			paramOut["component_path"] = m.Component
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if m.Sort != nil {
		if *m.Sort == -1 {
			paramOut["sort"] = nil
		} else {
			paramOut["sort"] = m.Sort
		}
	}

	if m.KeepAlive != nil {
		paramOut["keep_alive"] = m.KeepAlive
	}

	if m.Title != nil {
		if *m.Title != "" {
			paramOut["title"] = m.Title
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if m.Icon != nil {
		if *m.Icon != "" {
			paramOut["icon"] = m.Icon
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Menu{}).Where("id = ?", m.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (m *MenuDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.Menu
	global.DB.Where("id = ?", m.ID).Find(&record)
	err := global.DB.Where("id = ?", m.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (m *MenuUpdateApis) Update() response.Common {
	if m.ApiIDs == nil {
		return response.Failure(util.ErrorInvalidJSONParameters)
	}

	if len(*m.ApiIDs) == 0 {
		err := global.DB.Where("menu_id = ?", m.MenuID).Delete(&model.Menu{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Failure(util.ErrorFailToDeleteRecord)
		}
		return response.Success()
	}

	//先删掉原始记录
	err := global.DB.Where("menu_id = ?", m.MenuID).Delete(&model.MenuAndApi{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	//再增加新的记录
	var paramOut []model.MenuAndApi
	for _, apiID := range *m.ApiIDs {
		var record model.MenuAndApi
		if m.Creator > 0 {
			record.Creator = &m.Creator
		}
		if m.LastModifier > 0 {
			record.LastModifier = &m.LastModifier
		}

		record.MenuID = m.MenuID
		record.ApiID = apiID
		paramOut = append(paramOut, record)
	}

	for i := range paramOut {
		//计算有修改值的字段数，分别进行不同处理
		tempParamOut, err1 := util.StructToMap(paramOut[i])
		if err1 != nil {
			return response.Failure(util.ErrorFailToUpdateRecord)
		}
		paramOutForCounting := util.MapCopy(tempParamOut,
			"Creator", "LastModifier", "CreateAt", "UpdatedAt")

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
	var param1 rbacUpdatePolicyByMenuID
	param1.MenuID = m.MenuID
	err = param1.Update()
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRBACPoliciesByMenuID)
	}

	return response.Success()
}

func (m *MenuGetList) GetList() response.List {
	if m.UserID == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	db := global.DB.Model(&model.Menu{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if m.Group != "" {
		db = db.Where("group = ?", m.Group)
	}

	var roleIDs []int64
	global.DB.Model(&model.UserAndRole{}).Where("user_id = ?", m.UserID).
		Select("role_id").Find(&roleIDs)
	var menuIDs []int64
	global.DB.Model(&model.RoleAndMenu{}).Where("role_id in ?", roleIDs).
		Select("menu_id").Find(&menuIDs)
	db = db.Where("id in ?", menuIDs)

	//count
	var count int64
	db.Count(&count)

	//order
	orderBy := m.SortingInput.OrderBy
	desc := m.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Menu{}, orderBy)
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
	if m.PagingInput.Page > 0 {
		page = m.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if m.PagingInput.PageSize != nil && *m.PagingInput.PageSize >= 0 &&
		*m.PagingInput.PageSize <= global.Config.MaxPageSize {

		pageSize = *m.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []MenuOutput
	db.Model(&model.Menu{}).Find(&data)

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

func (m *MenuGetTree) GetTree() response.List {
	if m.UserID == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	db := global.DB.Model(&model.Menu{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	db = db.Where("superior_id is null")

	if m.Group != "" {
		db = db.Where("group = ?", m.Group)
	}

	var roleIDs []int64
	global.DB.Model(&model.UserAndRole{}).Where("user_id = ?", m.UserID).
		Select("role_id").Find(&roleIDs)
	var menuIDs []int64
	global.DB.Model(&model.RoleAndMenu{}).Where("role_id in ?", roleIDs).
		Select("menu_id").Find(&menuIDs)
	db = db.Where("id in ?", menuIDs)

	//count
	var count int64
	db.Count(&count)

	//order
	orderBy := m.SortingInput.OrderBy
	desc := m.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Menu{}, orderBy)
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
	if m.PagingInput.Page > 0 {
		page = m.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if m.PagingInput.PageSize != nil && *m.PagingInput.PageSize >= 0 &&
		*m.PagingInput.PageSize <= global.Config.MaxPageSize {

		pageSize = *m.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []MenuOutput
	db.Model(&model.Menu{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	for i := range data {
		data[i].Children = getMenuTree(data[i].ID)
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

func getMenuTree(superiorID int64) []MenuOutput {
	var result []MenuOutput
	res := global.DB.Model(model.Menu{}).
		Where("superior_id = ?", superiorID).Find(&result)
	if res.RowsAffected == 0 {
		return nil
	}

	for i := range result {
		result[i].Children = getMenuTree(result[i].ID)
	}
	return result
}
