package service

import (
	"github.com/mitchellh/mapstructure"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

// projectService 没有数据、只有方法，所有的数据都放在DTO里
//这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
//所有的增删改查都交给DAO层处理，否则service层会非常庞大
type projectService struct{}

func (projectService) Get(projectID int) response.Common {
	var result dto.ProjectGetDTO
	//查主表
	err := global.DB.Model(model.Project{}).
		Where("id = ?", projectID).First(&result).Error
	if err != nil {
		return response.Failure(util.ErrorRecordNotFound)
	}
	//如果有部门id，就查部门信息
	if result.DepartmentID != nil {
		err = global.DB.Model(model.Department{}).
			Where("id=?", result.DepartmentID).First(&result.Department).Error
		if err != nil {
			result.Department = nil
		}
	}
	return response.SuccessWithData(result)
}

func (projectService) Create(paramIn *dto.ProjectCreateOrUpdateDTO) response.Common {
	//对dto进行清洗，生成需要的model
	var paramOut model.Project
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}
	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}
	if *paramIn.ProjectCode == "" {
		paramOut.ProjectCode = nil
	} else {
		paramOut.ProjectCode = paramIn.ProjectCode
	}
	if *paramIn.ProjectFullName == "" {
		paramOut.ProjectFullName = nil
	} else {
		paramOut.ProjectFullName = paramIn.ProjectFullName
	}
	if *paramIn.ProjectShortName == "" {
		paramOut.ProjectShortName = nil
	} else {
		paramOut.ProjectShortName = paramIn.ProjectShortName
	}
	if *paramIn.Country == "" {
		paramOut.Country = nil
	} else {
		paramOut.Country = paramIn.Country
	}
	if *paramIn.Province == "" {
		paramOut.Province = nil
	} else {
		paramOut.Province = paramIn.Province
	}
	if *paramIn.ProjectType == "" {
		paramOut.ProjectType = nil
	} else {
		paramOut.ProjectType = paramIn.ProjectType
	}
	if *paramIn.Amount == -1 {
		paramOut.Amount = nil
	} else {
		paramOut.Amount = paramIn.Amount
	}
	if *paramIn.Currency == "" {
		paramOut.Currency = nil
	} else {
		paramOut.Currency = paramIn.Currency
	}
	if *paramIn.ExchangeRate == -1 {
		paramOut.ExchangeRate = nil
	} else {
		paramOut.ExchangeRate = paramIn.ExchangeRate
	}
	if *paramIn.DepartmentID == -1 {
		paramOut.DepartmentID = nil
	} else {
		paramOut.DepartmentID = paramIn.DepartmentID
	}
	if *paramIn.RelatedPartyID == -1 {
		paramOut.RelatedPartyID = nil
	} else {
		paramOut.RelatedPartyID = paramIn.RelatedPartyID
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (projectService) CreateInBatches(paramIn []dto.ProjectCreateOrUpdateDTO) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut []model.Project
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	for i := range paramIn {
		var record model.Project
		//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
		if paramIn[i].Creator != nil {
			record.Creator = paramIn[i].Creator
		}
		if paramIn[i].LastModifier != nil {
			record.LastModifier = paramIn[i].LastModifier
		}
		if *paramIn[i].ProjectCode == "" {
			record.ProjectCode = nil
		} else {
			record.ProjectCode = paramIn[i].ProjectCode
		}
		if *paramIn[i].ProjectFullName == "" {
			record.ProjectFullName = nil
		} else {
			record.ProjectFullName = paramIn[i].ProjectFullName
		}
		if *paramIn[i].ProjectShortName == "" {
			record.ProjectShortName = nil
		} else {
			record.ProjectShortName = paramIn[i].ProjectShortName
		}
		if *paramIn[i].Country == "" {
			record.Country = nil
		} else {
			record.Country = paramIn[i].Country
		}
		if *paramIn[i].Province == "" {
			record.Province = nil
		} else {
			record.Province = paramIn[i].Province
		}
		if *paramIn[i].ProjectType == "" {
			record.ProjectType = nil
		} else {
			record.ProjectType = paramIn[i].ProjectType
		}
		if *paramIn[i].Amount == -1 {
			record.Amount = nil
		} else {
			record.Amount = paramIn[i].Amount
		}
		if *paramIn[i].Currency == "" {
			record.Currency = nil
		} else {
			record.Currency = paramIn[i].Currency
		}
		if *paramIn[i].ExchangeRate == -1 {
			record.ExchangeRate = nil
		} else {
			record.ExchangeRate = paramIn[i].ExchangeRate
		}
		if *paramIn[i].DepartmentID == -1 {
			record.DepartmentID = nil
		} else {
			record.DepartmentID = paramIn[i].DepartmentID
		}
		if *paramIn[i].RelatedPartyID == -1 {
			record.RelatedPartyID = nil
		} else {
			record.RelatedPartyID = paramIn[i].RelatedPartyID
		}
		paramOut = append(paramOut, record)
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

// Update 更新为什么要用dto？首先因为很多数据需要绑定，也就是一定要传参；
// 其次是需要清洗
func (projectService) Update(paramIn *dto.ProjectCreateOrUpdateDTO) response.Common {
	var paramOut model.Project
	paramOut.ID = paramIn.ID
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.ProjectCode == "" {
		paramOut.ProjectCode = nil
	} else {
		paramOut.ProjectCode = paramIn.ProjectCode
	}
	if *paramIn.ProjectFullName == "" {
		paramOut.ProjectFullName = nil
	} else {
		paramOut.ProjectFullName = paramIn.ProjectFullName
	}
	if *paramIn.ProjectShortName == "" {
		paramOut.ProjectShortName = nil
	} else {
		paramOut.ProjectShortName = paramIn.ProjectShortName
	}
	if *paramIn.Country == "" {
		paramOut.Country = nil
	} else {
		paramOut.Country = paramIn.Country
	}
	if *paramIn.Province == "" {
		paramOut.Province = nil
	} else {
		paramOut.Province = paramIn.Province
	}
	if *paramIn.ProjectType == "" {
		paramOut.ProjectType = nil
	} else {
		paramOut.ProjectType = paramIn.ProjectType
	}
	if *paramIn.Amount == -1 {
		paramOut.Amount = nil
	} else {
		paramOut.Amount = paramIn.Amount
	}
	if *paramIn.Currency == "" {
		paramOut.Currency = nil
	} else {
		paramOut.Currency = paramIn.Currency
	}
	if *paramIn.ExchangeRate == -1 {
		paramOut.ExchangeRate = nil
	} else {
		paramOut.ExchangeRate = paramIn.ExchangeRate
	}
	if *paramIn.DepartmentID == -1 {
		paramOut.DepartmentID = nil
	} else {
		paramOut.DepartmentID = paramIn.DepartmentID
	}
	if *paramIn.RelatedPartyID == -1 {
		paramOut.RelatedPartyID = nil
	} else {
		paramOut.RelatedPartyID = paramIn.RelatedPartyID
	}

	//清洗完毕，开始update
	err := global.DB.Where("id = ?", paramOut.ID).Omit("created_at", "creator").Save(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (projectService) Delete(projectID int) response.Common {
	err := global.DB.Delete(&model.Project{}, projectID).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (projectService) List(paramIn dto.ProjectListDTO) response.List {
	db := global.DB
	//生成sql查询条件
	sqlCondition := util.NewSqlCondition()
	//对paramIn进行清洗
	//这部分是用于where的参数
	//如果用户角色是管理员或公司级，就不作处理
	if util.IsInSlice("管理员", paramIn.RoleNames) ||
		util.IsInSlice("公司级", paramIn.RoleNames) {
	} else if util.IsInSlice("事业部级", paramIn.RoleNames) {
		var departmentIDs []int
		if len(paramIn.BusinessDivisionIDs) > 0 {
			global.DB.Model(&model.Department{}).
				Where("superior_id in ?", paramIn.BusinessDivisionIDs).
				Select("id").Find(&departmentIDs)
		}
		if len(departmentIDs) > 0 {
			sqlCondition.In("department_id", departmentIDs)
		} else {
			sqlCondition.Where("department_id", -1)
		}

	} else if util.SliceIncludes(paramIn.RoleNames, "部门级") {
		if len(paramIn.DepartmentIDs) > 0 {
			sqlCondition.In("department_id", paramIn.DepartmentIDs)
		} else {
			sqlCondition.Where("department_id", -1)
		}

	} else { //为以后的”项目级“预留的功能
		if len(paramIn.DepartmentIDs) > 0 {
			sqlCondition.In("department_id", paramIn.DepartmentIDs)
		} else {
			sqlCondition.Where("department_id", -1)
		}
	}

	if paramIn.Page > 0 {
		sqlCondition.Paging.Page = paramIn.Page
	}
	//如果参数里的pageSize是整数且大于0、小于等于上限：
	maxPagingSize := global.Config.PagingConfig.MaxPageSize
	if paramIn.PageSize > 0 && paramIn.PageSize <= maxPagingSize {
		sqlCondition.Paging.PageSize = paramIn.PageSize
	}

	if len(paramIn.DepartmentIDIn) > 0 {
		sqlCondition.In("department_id", paramIn.DepartmentIDIn)
	}

	if paramIn.DepartmentNameLike != nil && *paramIn.DepartmentNameLike != "" {
		var departmentIDs []int
		err := global.DB.Model(model.Department{}).
			Where("name LIKE ?", "%"+*paramIn.DepartmentNameLike+"%").
			Select("id").Find(&departmentIDs).Error
		if err != nil {
			return response.FailureForList(util.ErrorInvalidJSONParameters)
		}
		sqlCondition.In("department_id", departmentIDs)
	}

	if paramIn.ProjectNameLike != nil && *paramIn.ProjectNameLike != "" {
		db = db.Where("project_full_name like ?", "%"+*paramIn.ProjectNameLike+"%").
			Or("project_short_name like ?", "%"+*paramIn.ProjectNameLike+"%")
	}

	//这部分是用于order的参数
	orderBy := paramIn.OrderBy
	if orderBy != "" {
		ok := sqlCondition.ValidateColumn(orderBy, model.Project{})
		if ok {
			sqlCondition.Sorting.OrderBy = orderBy
		}
	}
	desc := paramIn.Desc
	if desc == true {
		sqlCondition.Sorting.Desc = true
	} else {
		sqlCondition.Sorting.Desc = false
	}
	//需要先count再find，因为find会改变db的指针
	totalRecords := sqlCondition.Count(db, model.Project{})
	tempList := sqlCondition.Find(db, model.Project{})
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	//tempList是map，需要转成structure才能使用
	var list []dto.ProjectGetDTO
	_ = mapstructure.Decode(&tempList, &list)

	for i := range list {
		//如果有部门id，就查部门信息
		if list[i].DepartmentID != nil {
			err := global.DB.Model(model.Department{}).
				Where("id=?", list[i].DepartmentID).First(&list[i].Department).Error
			if err != nil {
				list[i].Department = nil
			}
		}
	}

	return response.List{
		Data: list,
		Paging: &dto.PagingDTO{
			Page:         sqlCondition.Paging.Page,
			PageSize:     sqlCondition.Paging.PageSize,
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
