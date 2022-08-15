package service

import (
	"learn-go/dao"
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
	"learn-go/serializer/response"
	"learn-go/util"
)

// DepartmentService 没有数据、只有方法，所有的数据都放在DTO里
//这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
//所有的增删改查都交给DAO层处理，否则service层会非常庞大
type departmentService struct{}

func (departmentService) Get(departmentID int) response.Common {
	result := dao.DepartmentDAO.Get(departmentID)
	if result == nil {
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (departmentService) Create(paramIn *dto.DepartmentCreateAndUpdateDTO) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut model.Department
	paramOut.Name = paramIn.Name
	paramOut.Level = paramIn.Level
	//model.Department的SuperiorID为指针，需要处理
	if *paramIn.SuperiorID == -1 {
		paramOut.SuperiorID = nil
	} else {
		paramOut.SuperiorID = paramIn.SuperiorID
	}
	err := dao.DepartmentDAO.Create(&paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

// Update 更新为什么要用dto？首先因为很多数据需要绑定，也就是一定要传参；
// 其次是需要清洗
func (departmentService) Update(paramIn *dto.DepartmentCreateAndUpdateDTO) response.Common {
	var paramOut model.Department
	paramOut.ID = paramIn.ID
	paramOut.Name = paramIn.Name
	paramOut.Level = paramIn.Level
	//model.Department的SuperiorID为指针，需要处理
	if *paramIn.SuperiorID == -1 {
		paramOut.SuperiorID = nil
	} else {
		paramOut.SuperiorID = paramIn.SuperiorID
	}

	//清洗完毕，开始update
	err := dao.DepartmentDAO.Update(&paramOut)
	//拿到dao层的返回结果，进行处理
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (departmentService) Delete(departmentID int) response.Common {
	err := dao.DepartmentDAO.Delete(departmentID)
	if err != nil {
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
	if paramIn.IDGte != nil {
		sqlCondition.Gte("id", *paramIn.IDGte)
	}
	if paramIn.IDLte != nil {
		sqlCondition.Lte("id", *paramIn.IDLte)
	}
	if paramIn.Name != nil && *paramIn.Name != "" {
		sqlCondition = sqlCondition.Equal("name", paramIn.Name)
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

	list := sqlCondition.Find(model.Department{})
	totalRecords := sqlCondition.Count(model.Department{})
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
