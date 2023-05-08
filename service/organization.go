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

type OrganizationGet struct {
	SnowID int64
}

type OrganizationCreate struct {
	Creator        int64
	LastModifier   int64
	SuperiorSnowID int64  `json:"superior_snow_id" binding:"required"` //上级机构ID
	Name           string `json:"name" binding:"required"`             //名称
	//Sort           int    `json:"sort" binding:"required"`       //级别，如公司、事业部、部门等
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type OrganizationUpdate struct {
	LastModifier   int64
	SnowID         int64
	Name           *string `json:"name"`             //名称
	SuperiorSnowID *int64  `json:"superior_snow_id"` //上级机构ID
}

type OrganizationDelete struct {
	SnowID int64
}

type OrganizationGetList struct {
	list.Input
	//list.DataScopeInput
	SuperiorSnowID int64  `json:"superior_snow_id,omitempty"`
	Name           string `json:"name,omitempty"`
	NameInclude    string `json:"name_include,omitempty"`
}

type OrganizationOutput struct {
	Creator        *int64 `json:"creator"`
	LastModifier   *int64 `json:"last_modifier"`
	SnowID         int64  `json:"snow_id"`
	Name           string `json:"name"`             //名称
	SuperiorSnowID *int   `json:"superior_snow_id"` //上级机构id
}

func (d *OrganizationGet) Get() response.Common {
	var result OrganizationOutput

	err := global.DB.Model(model.Organization{}).
		Where("snow_id = ?", d.SnowID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	return response.SuccessWithData(result)
}

func (d *OrganizationCreate) Create() response.Common {
	var paramOut model.Organization

	if d.Creator > 0 {
		paramOut.Creator = &d.Creator
	}

	if d.LastModifier > 0 {
		paramOut.LastModifier = &d.LastModifier
	}

	paramOut.SnowID = idgen.NextId()

	paramOut.Name = d.Name

	paramOut.SuperiorSnowID = &d.SuperiorSnowID

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

	if d.SuperiorSnowID != nil {
		if *d.SuperiorSnowID > 0 {
			paramOut["superior_snow_id"] = d.SuperiorSnowID
		} else if *d.SuperiorSnowID == -1 {
			paramOut["superior_snow_id"] = nil
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

	err := global.DB.Model(&model.Organization{}).
		Where("snow_id = ?", d.SnowID).Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (d *OrganizationDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.Organization
	global.DB.Where("snow_id = ?", d.SnowID).Find(&record)
	err := global.DB.Where("snow_id = ?", d.SnowID).Delete(&record).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (d *OrganizationGetList) GetList() response.List {
	db := global.DB.Model(&model.Organization{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if d.SuperiorSnowID > 0 {
		db = db.Where("superior_snow_id = ?", d.SuperiorSnowID)
	}

	if d.Name != "" {
		db = db.Where("name = ?", d.Name)
	}

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
			db = db.Order("snow_id desc")
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

//func DeptList2DeptTree(deptList []model.SysDept, pCode string) []model.SysDept {
//	var deptTree []model.SysDept
//	for _, v := range deptList {
//		if v.ParentCode == pCode {
//			v.Children = DeptList2DeptTree(deptList, v.DeptCode)
//			deptTree = append(deptTree, v)
//		}
//	}
//	return deptTree
//}
