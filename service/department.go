package service

import (
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type departmentService struct{}

func (departmentService) Get(departmentID int) response.Common {
	var result dto.DepartmentOutput

	err := global.DB.Model(model.Department{}).
		Where("id = ?", departmentID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	return response.SuccessWithData(result)
}

func (departmentService) Create(paramIn dto.DepartmentCreate) response.Common {
	var paramOut model.Department

	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	paramOut.Name = paramIn.Name

	paramOut.Level = paramIn.Level

	paramOut.SuperiorID = &paramIn.SuperiorID

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (departmentService) Update(paramIn dto.DepartmentUpdate) response.Common {
	paramOut := make(map[string]any)

	if paramIn.LastModifier > 0 {
		paramOut["last_modifier"] = paramIn.LastModifier
	}

	if paramIn.Name != nil {
		if *paramIn.Name != "" {
			paramOut["name"] = paramIn.Name
		} else {
			paramOut["name"] = nil
		}
	}

	if paramIn.Level != nil {
		if *paramIn.Level != "" {
			paramOut["level"] = paramIn.Level
		} else {
			paramOut["level"] = nil
		}
	}

	if paramIn.SuperiorID != nil {
		if *paramIn.SuperiorID > 0 {
			paramOut["superior_id"] = paramIn.SuperiorID
		} else if *paramIn.SuperiorID == 0 {
			paramOut["superior_id"] = nil
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Department{}).
		Where("id = ?", paramIn.ID).Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (departmentService) Delete(paramIn dto.DepartmentDelete) response.Common {
	//由于删除需要做两件事：软删除+记录删除人，所以需要用事务
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//这里记录删除人，在事务中必须放在前面
		//如果放后面，由于是软删除，系统会找不到这条记录，导致无法更新
		err := tx.Debug().Model(&model.Department{}).Where("id = ?", paramIn.ID).
			Update("deleter", paramIn.Deleter).Error
		if err != nil {
			return err
		}
		//这里删除记录
		err = tx.Delete(&model.Department{}, paramIn.ID).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

// 老式写法，待修改
func (departmentService) List(paramIn dto.DepartmentListOld) response.List {
	//生成sql查询条件
	sqlCondition := util.NewSqlCondition()

	//对paramIn进行清洗
	//这部分是用于where的参数
	if paramIn.Page > 0 {
		sqlCondition.Paging.Page = paramIn.Page
	}
	//如果参数里的pageSize是整数且大于0、小于等于上限：
	maxPagingSize := global.Config.PagingConfig.MaxPageSize
	if paramIn.PageSize > 0 && paramIn.PageSize <= maxPagingSize {
		sqlCondition.Paging.PageSize = paramIn.PageSize
	}
	if id := paramIn.ID; id > 0 {
		sqlCondition.Equal("id", id)
	}
	if paramIn.SuperiorID != nil {
		sqlCondition.Equal("superior_id", *paramIn.SuperiorID)
	}
	if paramIn.Level != nil && *paramIn.Level != "" {
		sqlCondition.Equal("level", *paramIn.Level)
	}
	if paramIn.Name != nil && *paramIn.Name != "" {
		sqlCondition = sqlCondition.Equal("name", *paramIn.Name)
	}
	if paramIn.NameLike != nil && *paramIn.NameLike != "" {
		sqlCondition = sqlCondition.Like("name", *paramIn.NameLike)
	}
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
				sqlCondition.In("id", departmentIDs)
			} else {
				sqlCondition.Where("id", -1)
			}

		} else if util.SliceIncludes(paramIn.RoleNames, "部门级") {
			if len(paramIn.DepartmentIDs) > 0 {
				sqlCondition.In("id", paramIn.DepartmentIDs)
			} else {
				sqlCondition.Where("id", -1)
			}

		} else { //为以后的”项目级“预留的功能
			if len(paramIn.DepartmentIDs) > 0 {
				sqlCondition.In("id", paramIn.DepartmentIDs)
			} else {
				sqlCondition.Where("id", -1)
			}
		}
	}

	//这部分是用于order的参数
	orderBy := paramIn.OrderBy
	if orderBy != "" {
		ok := sqlCondition.FieldIsInModel(model.Department{}, orderBy)
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

	totalRecords := sqlCondition.Count(global.DB, model.Department{})
	tempList := sqlCondition.Find(global.DB, model.Department{})
	totalPages := util.GetTotalNumberOfPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var list []dto.DepartmentOutputOld
	_ = mapstructure.Decode(&tempList, &list)

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
