package service

import (
	"github.com/mitchellh/mapstructure"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

// errorLogService 没有数据、只有方法，所有的数据都放在DTO里
//这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
//所有的增删改查都交给DAO层处理，否则service层会非常庞大
type errorLogService struct{}

func (errorLogService) Get(errorLogID int) response.Common {
	var result dto.ErrorLogGetDTO
	err := global.DB.Model(model.ErrorLog{}).
		Where("id = ?", errorLogID).First(&result).Error
	if err != nil {
		return response.Failure(util.ErrorRecordNotFound)
	}
	if result.Date != nil {
		date := *result.Date
		*result.Date = date[:10]
	}
	return response.SuccessWithData(result)
}

func (errorLogService) Create(paramIn *dto.ErrorLogCreateOrUpdateDTO) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut model.ErrorLog
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}

	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.Detail == "" {
		paramOut.Detail = nil
	} else {
		paramOut.Detail = paramIn.Detail
	}

	if *paramIn.Date == "" {
		paramOut.Date = nil
	} else {
		date, err := time.Parse("2006-01-02", *paramIn.Date)
		if err != nil {
			return response.Failure(util.ErrorInvalidJSONParameters)
		} else {
			paramOut.Date = &date
		}
	}

	if *paramIn.MajorCategory == "" {
		paramOut.MajorCategory = nil
	} else {
		paramOut.MajorCategory = paramIn.MajorCategory
	}

	if *paramIn.MinorCategory == "" {
		paramOut.MinorCategory = nil
	} else {
		paramOut.MinorCategory = paramIn.MinorCategory
	}

	if *paramIn.IsResolved == false {
		temp := false
		paramOut.IsResolved = &temp
	} else {
		temp := true
		paramOut.IsResolved = &temp
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

// Update 更新为什么要用dto？首先因为很多数据需要绑定，也就是一定要传参；
// 其次是需要清洗
func (errorLogService) Update(paramIn *dto.ErrorLogCreateOrUpdateDTO) response.Common {
	var paramOut model.ErrorLog
	paramOut.ID = paramIn.ID
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.Detail == "" {
		paramOut.Detail = nil
	} else {
		paramOut.Detail = paramIn.Detail
	}

	if *paramIn.Date == "" {
		paramOut.Date = nil
	} else {
		date, err := time.Parse("2006-01-02", *paramIn.Date)
		if err != nil {
			return response.Failure(util.ErrorInvalidJSONParameters)
		} else {
			paramOut.Date = &date
		}
	}

	if *paramIn.MajorCategory == "" {
		paramOut.MajorCategory = nil
	} else {
		paramOut.MajorCategory = paramIn.MajorCategory
	}

	if *paramIn.MinorCategory == "" {
		paramOut.MinorCategory = nil
	} else {
		paramOut.MinorCategory = paramIn.MinorCategory
	}

	if *paramIn.IsResolved == false {
		temp := false
		paramOut.IsResolved = &temp
	} else {
		temp := true
		paramOut.IsResolved = &temp
	}

	//清洗完毕，开始update
	err := global.DB.Where("id = ?", paramOut.ID).Omit("created_at", "creator").Save(&paramOut).Error
	//拿到dao层的返回结果，进行处理
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (errorLogService) Delete(errorLogID int) response.Common {
	err := global.DB.Delete(&model.ErrorLog{}, errorLogID).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (errorLogService) List(paramIn dto.DisassemblyListDTO) response.List {
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

	if paramIn.ProjectID != nil {
		sqlCondition.Equal("project_id", *paramIn.ProjectID)
	}

	if paramIn.SuperiorID != nil {
		sqlCondition.Equal("superior_id", *paramIn.SuperiorID)
	}

	if paramIn.Level != nil {
		sqlCondition.Equal("level", *paramIn.Level)
	}

	if paramIn.LevelGte != nil {
		sqlCondition.Gte("level", *paramIn.LevelGte)
	}

	if paramIn.LevelLte != nil {
		sqlCondition.Lte("level", *paramIn.LevelLte)
	}

	//这部分是用于order的参数
	orderBy := paramIn.OrderBy
	if orderBy != "" {
		ok := sqlCondition.ValidateColumn(orderBy, model.ErrorLog{})
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

	tempList := sqlCondition.Find(global.DB, model.ErrorLog{})
	totalRecords := sqlCondition.Count(global.DB, model.ErrorLog{})
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var list []dto.DisassemblyGetDTO
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
