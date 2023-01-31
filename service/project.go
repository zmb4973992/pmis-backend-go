package service

import (
	"github.com/mitchellh/mapstructure"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type projectService struct{}

func (projectService) Get(projectID int) response.Common {
	var result dto.ProjectOutput
	err := global.DB.Model(model.Project{}).
		Where("id = ?", projectID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}

	//如果有部门id，就查部门信息
	if result.DepartmentID != nil {
		err = global.DB.Model(model.Department{}).
			Where("id=?", result.DepartmentID).First(&result.Department).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			result.Department = nil
		}
	}
	return response.SucceedWithData(result)
}

func (projectService) Create(paramIn dto.ProjectCreate) response.Common {
	var paramOut model.Project

	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	if paramIn.ProjectCode != "" {
		paramOut.ProjectCode = &paramIn.ProjectCode
	}

	if paramIn.ProjectFullName != "" {
		paramOut.ProjectFullName = &paramIn.ProjectFullName
	}

	if paramIn.ProjectShortName != "" {
		paramOut.ProjectShortName = &paramIn.ProjectShortName
	}

	if paramIn.Country != "" {
		paramOut.Country = &paramIn.Country
	}

	if paramIn.Province != "" {
		paramOut.Province = &paramIn.Province
	}

	if paramIn.ProjectType != "" {
		paramOut.ProjectType = &paramIn.ProjectType
	}

	if paramIn.Amount != nil {
		paramOut.Amount = paramIn.Amount
	}

	if paramIn.Currency != "" {
		paramOut.Currency = &paramIn.Currency
	}

	if paramIn.ExchangeRate != nil {
		paramOut.ExchangeRate = paramIn.ExchangeRate
	}

	if paramIn.DepartmentID != 0 {
		paramOut.DepartmentID = &paramIn.DepartmentID
	}

	if paramIn.RelatedPartyID != 0 {
		paramOut.RelatedPartyID = &paramIn.RelatedPartyID
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (projectService) Update(paramIn *dto.ProjectUpdate) response.Common {
	paramOut := make(map[string]any)

	if paramIn.LastModifier > 0 {
		paramOut["last_modifier"] = paramIn.LastModifier
	}

	if paramIn.ProjectCode != nil {
		if *paramIn.ProjectCode != "" {
			paramOut["project_code"] = paramIn.ProjectCode
		} else {
			paramOut["project_code"] = nil
		}
	}

	if paramIn.ProjectFullName != nil {
		if *paramIn.ProjectFullName != "" {
			paramOut["project_full_name"] = paramIn.ProjectFullName
		} else {
			paramOut["project_full_name"] = nil
		}
	}

	if paramIn.ProjectShortName != nil {
		if *paramIn.ProjectShortName != "" {
			paramOut["project_short_name"] = paramIn.ProjectShortName
		} else {
			paramOut["project_short_name"] = nil
		}
	}

	if paramIn.Country != nil {
		if *paramIn.Country != "" {
			paramOut["country"] = paramIn.Country
		} else {
			paramOut["country"] = nil
		}
	}

	if paramIn.Province != nil {
		if *paramIn.Province != "" {
			paramOut["province"] = paramIn.Province
		} else {
			paramOut["province"] = nil
		}
	}

	if paramIn.ProjectType != nil {
		if *paramIn.ProjectType != "" {
			paramOut["project_type"] = paramIn.ProjectType
		} else {
			paramOut["project_type"] = nil
		}
	}

	if paramIn.Amount != nil {
		if *paramIn.Amount != -1 {
			paramOut["amount"] = paramIn.Amount
		} else {
			paramOut["amount"] = nil
		}
	}

	if paramIn.Currency != nil {
		if *paramIn.Currency != "" {
			paramOut["currency"] = paramIn.Currency
		} else {
			paramOut["currency"] = nil
		}
	}

	if paramIn.ExchangeRate != nil {
		if *paramIn.ExchangeRate != -1 {
			paramOut["exchange_rate"] = paramIn.ExchangeRate
		} else {
			paramOut["exchange_rate"] = nil
		}
	}

	if paramIn.DepartmentID != nil {
		if *paramIn.DepartmentID != 0 {
			paramOut["department_id"] = paramIn.DepartmentID
		} else {
			paramOut["department_id"] = nil
		}
	}

	if paramIn.RelatedPartyID != nil {
		if *paramIn.RelatedPartyID != 0 {
			paramOut["related_party_id"] = paramIn.RelatedPartyID
		} else {
			paramOut["related_party_id"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Project{}).Where("id = ?", paramIn.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (projectService) Delete(projectID int) response.Common {
	err := global.DB.Delete(&model.Project{}, projectID).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (projectService) List(paramIn dto.ProjectListOld) response.List {
	db := global.DB
	//生成sql查询条件
	sqlCondition := util.NewSqlCondition()
	//对paramIn进行清洗
	//这部分是用于where的参数
	if paramIn.VerifyRole != nil && *paramIn.VerifyRole == true {
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

		} else if util.SliceIncludesOld(paramIn.RoleNames, "部门级") {
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
			global.SugaredLogger.Errorln(err)
			return response.FailForList(util.ErrorInvalidJSONParameters)
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
		ok := sqlCondition.FieldIsInModel(model.Project{}, orderBy)
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
	totalPages := util.GetTotalNumberOfPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailForList(util.ErrorRecordNotFound)
	}

	//tempList是map，需要转成structure才能使用
	var list []dto.ProjectOutputOld
	_ = mapstructure.Decode(&tempList, &list)

	for i := range list {
		//如果有部门id，就查部门信息
		if list[i].DepartmentID != nil {
			err := global.DB.Model(model.Department{}).
				Where("id=?", list[i].DepartmentID).First(&list[i].Department).Error
			if err != nil {
				global.SugaredLogger.Errorln(err)
				list[i].Department = nil
			}
		}
	}

	return response.List{
		Data: list,
		Paging: &dto.PagingOutput{
			Page:            sqlCondition.Paging.Page,
			PageSize:        sqlCondition.Paging.PageSize,
			NumberOfPages:   totalPages,
			NumberOfRecords: totalRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
