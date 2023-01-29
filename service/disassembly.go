package service

import (
	"github.com/mitchellh/mapstructure"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

// disassemblyService 没有数据、只有方法，所有的数据都放在DTO里
// 这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
// 所有的增删改查都交给DAO层处理，否则service层会非常庞大
type disassemblyService struct{}

func (disassemblyService) Get(disassemblyID int) response.Common {
	var result dto.DisassemblyOutput
	err := global.DB.Model(model.Disassembly{}).
		Where("id = ?", disassemblyID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (disassemblyService) Tree(paramIn dto.DisassemblyTree) response.Common {
	//根据project_id获取disassembly_id
	var disassemblyID int
	err := global.DB.Model(model.Disassembly{}).Select("id").
		Where("project_id = ?", paramIn.ProjectID).Where("level = 1").
		First(&disassemblyID).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	//第一轮查找
	var result1 []dto.DisassemblyTreeOutput
	err = global.DB.Model(model.Disassembly{}).
		Where("id = ?", disassemblyID).First(&result1).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	//第二轮查找
	var result2 []dto.DisassemblyTreeOutput
	global.DB.Model(model.Disassembly{}).
		Where("superior_id = ?", disassemblyID).Find(&result2)
	for index2 := range result2 {
		//第三轮查找
		var result3 []dto.DisassemblyTreeOutput
		global.DB.Model(model.Disassembly{}).
			Where("superior_id = ?", result2[index2].ID).Find(&result3)
		//第四轮查找
		for index3 := range result3 {
			var result4 []dto.DisassemblyTreeOutput
			global.DB.Model(model.Disassembly{}).
				Where("superior_id = ?", result3[index3].ID).Find(&result4)
			for index4 := range result4 {
				var result5 []dto.DisassemblyTreeOutput
				global.DB.Model(model.Disassembly{}).
					Where("superior_id = ?", result4[index4].ID).Find(&result5)
				result4[index4].Children = append(result4[index4].Children, result5...)
			}
			result3[index3].Children = append(result3[index3].Children, result4...)
		}
		result2[index2].Children = append(result2[index2].Children, result3...)
	}
	result1[0].Children = append(result1[0].Children, result2...)
	return response.SuccessWithData(result1)
}

func (disassemblyService) Create(paramIn *dto.DisassemblyCreateOrUpdate) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut model.Disassembly
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}

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

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (disassemblyService) CreateInBatches(paramIn []dto.DisassemblyCreateOrUpdate) response.Common {
	//对dto进行清洗，生成dao层需要的model
	var paramOut []model.Disassembly
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	for i := range paramIn {
		var record model.Disassembly
		if paramIn[i].Creator != nil {
			record.Creator = paramIn[i].Creator
		}

		if paramIn[i].LastModifier != nil {
			record.LastModifier = paramIn[i].LastModifier
		}

		if *paramIn[i].Name != "" {
			record.Name = paramIn[i].Name
		}

		if *paramIn[i].Level != -1 {
			record.Level = paramIn[i].Level
		}

		if *paramIn[i].ProjectID != -1 {
			record.ProjectID = paramIn[i].ProjectID
		}

		if *paramIn[i].Weight != -1 {
			record.Weight = paramIn[i].Weight
		}

		if *paramIn[i].SuperiorID != -1 {
			record.SuperiorID = paramIn[i].SuperiorID
		}

		paramOut = append(paramOut, record)
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
func (disassemblyService) Update(paramIn *dto.DisassemblyCreateOrUpdate) response.Common {
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
	err := global.DB.Where("id = ?", paramOut.ID).Omit(fieldsToBeOmittedWhenUpdating...).Save(&paramOut).Error
	//拿到dao层的返回结果，进行处理
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (disassemblyService) Delete(disassemblyID int) response.Common {
	err := global.DB.Delete(&model.Disassembly{}, disassemblyID).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (disassemblyService) DeleteWithSubitems(disassemblyID int) response.Common {
	var ToBeDeletedIDs []int
	ToBeDeletedIDs = append(ToBeDeletedIDs, disassemblyID)
	//第一轮查找
	var result1 []int
	global.DB.Model(&model.Disassembly{}).Where("superior_id = ?", disassemblyID).
		Select("id").Find(&result1)
	//第二轮查找
	if len(result1) > 0 {
		ToBeDeletedIDs = append(ToBeDeletedIDs, result1...)
		var result2 []int
		global.DB.Model(&model.Disassembly{}).Where("superior_id IN ?", result1).
			Select("id").Find(&result2)
		//第三轮查找
		if len(result2) > 0 {
			ToBeDeletedIDs = append(ToBeDeletedIDs, result2...)
			var result3 []int
			global.DB.Model(&model.Disassembly{}).Where("superior_id IN ?", result2).
				Select("id").Find(&result3)
			//第四轮查找
			if len(result3) > 0 {
				ToBeDeletedIDs = append(ToBeDeletedIDs, result3...)
				var result4 []int
				global.DB.Model(&model.Disassembly{}).Where("superior_id IN ?", result3).
					Select("id").Find(&result4)
				ToBeDeletedIDs = append(ToBeDeletedIDs, result4...)
			}
		}
	}
	err := global.DB.Delete(&model.Disassembly{}, ToBeDeletedIDs).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (disassemblyService) List(paramIn dto.DisassemblyList) response.List {
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
		ok := sqlCondition.FieldIsInModel(model.Disassembly{}, orderBy)
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
	totalPages := util.GetTotalNumberOfPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var list []dto.DisassemblyOutput
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
