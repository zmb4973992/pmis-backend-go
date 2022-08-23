package service

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"learn-go/dao"
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
	"learn-go/serializer/response"
	"learn-go/util"
)

// errorLogService 没有数据、只有方法，所有的数据都放在DTO里
//这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
//所有的增删改查都交给DAO层处理，否则service层会非常庞大
type errorLogService struct{}

func (errorLogService) Get(errorLogID int) response.Common {
	//var param dto.ErrorLogGetDTO
	//把基础的拆解信息查出来
	var errorLog model.ErrorLog
	err := global.DB.Where("id = ?", errorLogID).First(&errorLog).Error
	if err != nil {
		return response.Failure(util.ErrorRecordNotFound)
	}

	//把所有查出的结果赋值给输出变量
	detail := errorLog.Detail
	fmt.Println(detail)

	//if disassembly.Level != nil {
	//	param.Level = disassembly.Level
	//}
	//if disassembly.Weight != nil {
	//	param.Weight = disassembly.Weight
	//}
	//if disassembly.SuperiorID != nil {
	//	param.SuperiorID = disassembly.SuperiorID
	//}
	//return &param
	//
	//errorLog := dao.DisassemblyDAO.Get(errorLogID)
	//if errorLog == nil {
	//	return response.Failure(util.ErrorRecordNotFound)
	//}
	return response.SuccessWithData(errorLog)
}

func (errorLogService) Create(paramIn *dto.DisassemblyCreateOrUpdateDTO) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut model.Disassembly
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}

	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.Name == "" { //这里不需要对paramIn.Name进行非空判定，因为前面的dto已经设定了必须绑定
		paramOut.Name = nil
	} else {
		paramOut.Name = paramIn.Name
	}
	if *paramIn.Level == -1 {
		paramOut.Level = nil
	} else {
		paramOut.Level = paramIn.Level
	}
	if *paramIn.ProjectID == -1 {
		paramOut.ProjectID = nil
	} else {
		paramOut.ProjectID = paramIn.ProjectID
	}
	if *paramIn.Weight == -1 {
		paramOut.Weight = nil
	} else {
		paramOut.Weight = paramIn.Weight
	}
	if *paramIn.SuperiorID == -1 {
		paramOut.SuperiorID = nil
	} else {
		paramOut.SuperiorID = paramIn.SuperiorID
	}

	err := dao.DisassemblyDAO.Create(&paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

// Update 更新为什么要用dto？首先因为很多数据需要绑定，也就是一定要传参；
// 其次是需要清洗
func (errorLogService) Update(paramIn *dto.DisassemblyCreateOrUpdateDTO) response.Common {
	var paramOut model.Disassembly
	paramOut.ID = paramIn.ID
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.Name == "" { //这里不需要对paramIn.Name进行非空判定，因为前面的dto已经设定了必须绑定
		paramOut.Name = nil
	} else {
		paramOut.Name = paramIn.Name
	}
	if *paramIn.Level == -1 {
		paramOut.Level = nil
	} else {
		paramOut.Level = paramIn.Level
	}
	if *paramIn.ProjectID == -1 {
		paramOut.ProjectID = nil
	} else {
		paramOut.ProjectID = paramIn.ProjectID
	}
	if *paramIn.Weight == -1 {
		paramOut.Weight = nil
	} else {
		paramOut.Weight = paramIn.Weight
	}
	if *paramIn.SuperiorID == -1 {
		paramOut.SuperiorID = nil
	} else {
		paramOut.SuperiorID = paramIn.SuperiorID
	}

	//清洗完毕，开始update
	err := dao.DisassemblyDAO.Update(&paramOut)
	//拿到dao层的返回结果，进行处理
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (errorLogService) Delete(errorLogID int) response.Common {
	err := dao.DisassemblyDAO.Delete(errorLogID)
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
		ok := sqlCondition.ValidateColumn(orderBy, model.Disassembly{})
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

	tempList := sqlCondition.Find(model.Disassembly{})
	totalRecords := sqlCondition.Count(model.Disassembly{})
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
