package service

import (
	"github.com/mitchellh/mapstructure"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

// dictionaryItemService 没有数据、只有方法，所有的数据都放在DTO里
//这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
//所有的增删改查都交给DAO层处理，否则service层会非常庞大
type dictionaryItemService struct{}

func (dictionaryItemService) Get(dictionaryTypeID int) response.Common {
	var result []dto.DictionaryItemOutputDTO
	err := global.DB.Model(model.DictionaryItem{}).
		Where("dictionary_type_id = ?", dictionaryTypeID).Find(&result).Error
	if err != nil || len(result) == 0 {
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (dictionaryItemService) Create(paramIn *dto.DictionaryItemCreateOrUpdateDTO) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut model.DictionaryItem
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}

	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	paramOut.DictionaryTypeID = paramIn.DictionaryTypeID

	paramOut.Name = paramIn.Name

	if *paramIn.Sort != -1 {
		paramOut.Sort = paramIn.Sort
	}

	if *paramIn.Remarks != "" {
		paramOut.Remarks = paramIn.Remarks
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (dictionaryItemService) CreateInBatches(paramIn []dto.DictionaryItemCreateOrUpdateDTO) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut []model.DictionaryItem
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	for i := range paramIn {
		var record model.DictionaryItem
		if paramIn[i].Creator != nil {
			record.Creator = paramIn[i].Creator
		}

		if paramIn[i].LastModifier != nil {
			record.LastModifier = paramIn[i].LastModifier
		}

		record.DictionaryTypeID = paramIn[i].DictionaryTypeID
		record.Name = paramIn[i].Name

		if *paramIn[i].Sort != -1 {
			record.Sort = paramIn[i].Sort
		}

		if *paramIn[i].Remarks != "" {
			record.Remarks = paramIn[i].Remarks
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
func (dictionaryItemService) Update(paramIn *dto.DisassemblyCreateOrUpdateDTO) response.Common {
	var paramOut model.Disassembly
	paramOut.ID = paramIn.ID
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.Name != "" {
		paramOut.Name = paramIn.Name
	}

	if *paramIn.Level != -1 {
		paramOut.Level = paramIn.Level
	}

	if *paramIn.ProjectID != -1 {
		paramOut.ProjectID = paramIn.ProjectID
	}

	if *paramIn.Weight != -1 {
		paramOut.Weight = paramIn.Weight
	}

	if *paramIn.SuperiorID != -1 {
		paramOut.SuperiorID = paramIn.SuperiorID
	}

	//清洗完毕，开始update
	err := global.DB.Where("id = ?", paramOut.ID).Omit("created_at", "creator").Save(&paramOut).Error
	//拿到dao层的返回结果，进行处理
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (dictionaryItemService) Delete(disassemblyID int) response.Common {
	err := global.DB.Delete(&model.Disassembly{}, disassemblyID).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (dictionaryItemService) List(paramIn dto.DisassemblyListDTO) response.List {
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

	tempList := sqlCondition.Find(global.DB, model.Disassembly{})
	totalRecords := sqlCondition.Count(global.DB, model.Disassembly{})
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var list []dto.DisassemblyOutputDTO
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
