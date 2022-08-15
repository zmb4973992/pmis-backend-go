package service

import (
	"learn-go/dao"
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
	"learn-go/serializer/response"
	"learn-go/util"
)

// ProjectDisassemblyService 没有数据、只有方法，所有的数据都放在DTO里
//这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
//所有的增删改查都交给DAO层处理，否则service层会非常庞大
type projectDisassemblyService struct{}

func (projectDisassemblyService) Get(projectDisassemblyID int) response.Common {
	result := dao.ProjectDisassemblyDAO.Get(projectDisassemblyID)
	if result == nil {
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (projectDisassemblyService) Create(paramIn *dto.ProjectDisassemblyCreateAndUpdateDTO) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut model.ProjectDisassembly
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
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

	err := dao.ProjectDisassemblyDAO.Create(&paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

// Update 更新为什么要用dto？首先因为很多数据需要绑定，也就是一定要传参；
// 其次是需要清洗
func (projectDisassemblyService) Update(paramIn *dto.ProjectDisassemblyCreateAndUpdateDTO) response.Common {
	var paramOut model.ProjectDisassembly
	paramOut.ID = paramIn.ID
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
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
	err := dao.ProjectDisassemblyDAO.Update(&paramOut)
	//拿到dao层的返回结果，进行处理
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (projectDisassemblyService) Delete(projectDisassemblyID int) response.Common {
	err := dao.ProjectDisassemblyDAO.Delete(projectDisassemblyID)
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (projectDisassemblyService) List(paramIn dto.ProjectDisassemblyListDTO) response.List {
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
	if paramIn.IDGte != nil {
		sqlCondition.Gte("id", *paramIn.IDGte)
	}
	if paramIn.IDLte != nil {
		sqlCondition.Lte("id", *paramIn.IDLte)
	}
	if paramIn.Name != nil && *paramIn.Name != "" {
		sqlCondition = sqlCondition.Equal("name", *paramIn.Name)
	}
	if paramIn.NameInclude != nil && *paramIn.NameInclude != "" {
		sqlCondition = sqlCondition.Include("name", *paramIn.NameInclude)
	}

	//这部分是用于order的参数
	orderBy := paramIn.OrderBy
	if orderBy != "" {
		ok := sqlCondition.ValidateColumn(orderBy, model.ProjectDisassembly{})
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

	list := sqlCondition.Find(model.ProjectDisassembly{})
	totalRecords := sqlCondition.Count(model.ProjectDisassembly{})
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(list) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
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
