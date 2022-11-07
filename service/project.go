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
	err := global.DB.Model(model.Project{}).
		Where("id = ?", projectID).First(&result).Error
	if err != nil {
		return response.Failure(util.ErrorRecordNotFound)
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

	tempList := sqlCondition.Find(model.Project{})
	totalRecords := sqlCondition.Count(model.Project{})
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var list []dto.ProjectGetDTO
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
