package service

import (
	"github.com/mitchellh/mapstructure"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

// DepartmentService 没有数据、只有方法，所有的数据都放在DTO里
// 这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
// 所有的增删改查都交给DAO层处理，否则service层会非常庞大
type departmentService struct{}

func (departmentService) Get(departmentID int) response.Common {
	var result dto.DepartmentGetDTO

	err := global.DB.Model(model.Department{}).Where("id = ?", departmentID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	return response.SuccessWithData(result)
}

func (departmentService) Create(paramIn *dto.DepartmentCreateOrUpdateDTO) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut model.Department

	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}

	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	//如果dto的字段是指针，那么需要看字段binding是否为required，然后判定是否为空或-1再进行处理；
	//如果dto的字段不是指针，需要看字段binding是否为required，然后可以赋值
	paramOut.Name = paramIn.Name
	paramOut.Level = paramIn.Level

	if *paramIn.SuperiorID != -1 {
		paramOut.SuperiorID = paramIn.SuperiorID
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

// Update 更新为什么要用dto？首先因为很多数据需要绑定，也就是一定要传参；
// 其次是需要清洗
func (departmentService) Update(paramIn *dto.DepartmentCreateOrUpdateDTO) response.Common {
	var paramOut model.Department
	paramOut.ID = paramIn.ID

	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	paramOut.Name = paramIn.Name
	paramOut.Level = paramIn.Level

	if *paramIn.SuperiorID != -1 {
		paramOut.SuperiorID = paramIn.SuperiorID
	}

	//清洗完毕，开始update
	err := global.DB.Where("id = ?", paramOut.ID).Omit("created_at", "creator").Save(&paramOut).Error
	//拿到dao层的返回结果，进行处理
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (departmentService) Delete(departmentID int) response.Common {
	err := global.DB.Delete(&model.Department{}, departmentID).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (departmentService) List(paramIn dto.DepartmentListDTO) response.List {
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
		ok := sqlCondition.ValidateColumn(orderBy, model.Department{})
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
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var list []dto.DepartmentGetDTO
	_ = mapstructure.Decode(&tempList, &list)

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
