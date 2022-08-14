package service

import (
	"github.com/mitchellh/mapstructure"
	"learn-go/dao"
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
	"learn-go/serializer/response"
	"learn-go/util"
	"time"
)

// operationRecordService 没有数据、只有方法，所有的数据都放在DTO里
//这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
//所有的增删改查都交给DAO层处理，否则service层会非常庞大
type operationRecordService struct{}

func (operationRecordService) Get(operationRecordID int) response.Common {
	result := dao.OperationRecordDAO.Get(operationRecordID)
	if result == nil {
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (operationRecordService) Create(paramIn *dto.OperationRecordCreateAndUpdateDTO) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut model.OperationRecord

	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if *paramIn.ProjectID == -1 { //这里不需要对paramIn.Name进行非空判定，因为前面的dto已经设定了必须绑定
		paramOut.ProjectID = nil
	} else {
		paramOut.ProjectID = paramIn.ProjectID
	}

	if *paramIn.OperatorID == -1 {
		paramOut.OperatorID = nil
	} else {
		paramOut.OperatorID = paramIn.OperatorID
	}

	if *paramIn.ProjectID == -1 {
		paramOut.ProjectID = nil
	} else {
		paramOut.ProjectID = paramIn.ProjectID
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

	if *paramIn.Action == "" {
		paramOut.Action = nil
	} else {
		paramOut.Action = paramIn.Action
	}

	if *paramIn.Detail == "" {
		paramOut.Detail = nil
	} else {
		paramOut.Detail = paramIn.Detail
	}

	err := dao.OperationRecordDAO.Create(&paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	return response.Success()
}

// Update 更新为什么要用dto？首先因为很多数据需要绑定，也就是一定要传参；
// 其次是需要清洗
func (operationRecordService) Update(paramIn *dto.OperationRecordCreateAndUpdateDTO) response.Common {
	var paramOut model.OperationRecord
	paramOut.ID = paramIn.ID
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	//这里不需要进行非空判定，因为前面的dto已经设定了必须绑定
	if *paramIn.ProjectID == -1 {
		paramOut.ProjectID = nil
	} else {
		paramOut.ProjectID = paramIn.ProjectID
	}

	if *paramIn.OperatorID == -1 {
		paramOut.OperatorID = nil
	} else {
		paramOut.OperatorID = paramIn.OperatorID
	}

	if *paramIn.ProjectID == -1 {
		paramOut.ProjectID = nil
	} else {
		paramOut.ProjectID = paramIn.ProjectID
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

	if *paramIn.Action == "" {
		paramOut.Action = nil
	} else {
		paramOut.Action = paramIn.Action
	}

	if *paramIn.Detail == "" {
		paramOut.Detail = nil
	} else {
		paramOut.Detail = paramIn.Detail
	}

	//清洗完毕，开始update
	err := dao.OperationRecordDAO.Update(&paramOut)
	//拿到dao层的返回结果，进行处理
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (operationRecordService) Delete(operationRecordID int) response.Common {
	err := dao.OperationRecordDAO.Delete(operationRecordID)
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (operationRecordService) List(paramIn dto.OperationRecordListDTO) response.List {
	//生成sql查询条件
	sqlCondition := util.NewSqlCondition()
	//对paramIn进行清洗
	//select columns
	//如果参数正确，那么指定字段的数据正常返回，其他字段返回空；
	//如果参数错误，就返回全部字段的数据
	if len(paramIn.SelectedColumns) > 0 {
		ok := sqlCondition.ValidateColumns(paramIn.SelectedColumns, model.OperationRecord{})
		if ok {
			sqlCondition.SelectedColumns = paramIn.SelectedColumns
		}
	}

	//分页
	if paramIn.Page > 0 {
		sqlCondition.Paging.Page = paramIn.Page
	}
	//如果参数里的pageSize是整数且大于0、小于等于上限：
	maxPagingSize := global.Config.PagingConfig.MaxPageSize
	if paramIn.PageSize > 0 && paramIn.PageSize <= maxPagingSize {
		sqlCondition.Paging.PageSize = paramIn.PageSize
	}

	//这部分是用于where的参数
	if id := paramIn.ID; id > 0 {
		sqlCondition.Equal("id", id)
	}
	if paramIn.ProjectID != nil && *paramIn.ProjectID != -1 {
		sqlCondition.Equal("project_id", *paramIn.ProjectID)
	}
	if paramIn.OperatorID != nil && *paramIn.OperatorID != -1 {
		sqlCondition.Equal("operator_id", *paramIn.OperatorID)
	}
	if paramIn.DateGte != nil && *paramIn.DateGte != "" {
		sqlCondition.Gte("date", *paramIn.DateGte)
	}
	if paramIn.DateLte != nil && *paramIn.DateLte != "" {
		sqlCondition.Lte("date", *paramIn.DateLte)
	}
	if paramIn.Action != nil && *paramIn.Action != "" {
		sqlCondition.Equal("action", *paramIn.Action)
	}

	//这部分是用于order的参数
	orderBy := paramIn.OrderBy
	if orderBy != "" {
		ok := sqlCondition.ValidateColumn(orderBy, model.OperationRecord{})
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

	tempList := sqlCondition.Find(model.OperationRecord{})
	totalRecords := sqlCondition.Count(model.OperationRecord{})
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		response.FailureForList(util.ErrorRecordNotFound)
	}

	//这里的tempList是基于model的，不能直接传给前端，要处理成dto才行
	//如果map的字段类型和struct的字段类型不匹配，数据不会同步过来
	var list []dto.OperationRecordGetDTO
	_ = mapstructure.Decode(&tempList, &list)

	//处理字段类型不匹配、或者有特殊格式要求的字段
	for k := range tempList {
		a := tempList[k]["date"].(*time.Time).Format("2006-01-02")
		list[k].Date = &a
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
