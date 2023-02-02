package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type disassembly struct{}

func (*disassembly) Get(disassemblyID int) response.Common {
	var result dto.DisassemblyOutput
	err := global.DB.Model(model.Disassembly{}).
		Where("id = ?", disassemblyID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}
	return response.SucceedWithData(result)
}

func (*disassembly) Tree(paramIn dto.DisassemblyTree) response.Common {
	//根据project_id获取disassembly_id
	var disassemblyID int
	err := global.DB.Model(model.Disassembly{}).Select("id").
		Where("project_id = ?", paramIn.ProjectID).Where("level = 1").
		First(&disassemblyID).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}

	//第一轮查找
	var result1 []dto.DisassemblyTreeOutput
	err = global.DB.Model(model.Disassembly{}).
		Where("id = ?", disassemblyID).First(&result1).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
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
	return response.SucceedWithData(result1)
}

func (*disassembly) Create(paramIn dto.DisassemblyCreate) response.Common {
	var paramOut model.Disassembly
	if paramIn.Creator > 0 {
		paramOut.Creator = &paramIn.Creator
	}

	if paramIn.LastModifier > 0 {
		paramOut.LastModifier = &paramIn.LastModifier
	}

	paramOut.Name = &paramIn.Name

	paramOut.ProjectID = &paramIn.ProjectID

	paramOut.Level = &paramIn.Level

	paramOut.Weight = &paramIn.Weight

	paramOut.SuperiorID = &paramIn.SuperiorID

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (*disassembly) CreateInBatches(paramIn []dto.DisassemblyCreate) response.Common {
	var paramOut []model.Disassembly
	for i := range paramIn {
		var record model.Disassembly
		if paramIn[i].Creator > 0 {
			record.Creator = &paramIn[i].Creator
		}

		if paramIn[i].LastModifier > 0 {
			record.LastModifier = &paramIn[i].LastModifier
		}

		record.Name = &paramIn[i].Name

		record.Level = &paramIn[i].Level

		record.ProjectID = &paramIn[i].ProjectID

		record.Weight = &paramIn[i].Weight

		record.SuperiorID = &paramIn[i].SuperiorID

		paramOut = append(paramOut, record)
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (*disassembly) Update(paramIn dto.DisassemblyUpdate) response.Common {
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

	if paramIn.ProjectID != nil {
		if *paramIn.ProjectID != 0 {
			paramOut["project_id"] = paramIn.ProjectID
		} else {
			paramOut["project_id"] = nil
		}
	}

	if paramIn.Level != nil {
		if *paramIn.Level != 0 {
			paramOut["level"] = paramIn.Level
		} else {
			paramOut["level"] = nil
		}
	}

	if paramIn.Weight != nil {
		if *paramIn.Weight != 0 {
			paramOut["weight"] = paramIn.Weight
		} else {
			paramOut["weight"] = nil
		}
	}

	if paramIn.SuperiorID != nil {
		if *paramIn.SuperiorID != 0 {
			paramOut["superior_id"] = paramIn.SuperiorID
		} else {
			paramOut["superior_id"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Disassembly{}).Where("id = ?", paramIn.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}
	return response.Succeed()
}

func (*disassembly) Delete(paramIn dto.DisassemblyDelete) response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.Disassembly
	global.DB.Where("id = ?", paramIn.ID).Find(&record)
	record.Deleter = &paramIn.Deleter
	err := global.DB.Where("id = ?", paramIn.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (*disassembly) DeleteWithSubitems(paramIn dto.DisassemblyDelete) response.Common {
	var ToBeDeletedIDs []int
	ToBeDeletedIDs = append(ToBeDeletedIDs, paramIn.ID)
	//第一轮查找
	var result1 []int
	global.DB.Model(&model.Disassembly{}).Where("superior_id = ?", paramIn.ID).
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

	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var records []model.Disassembly
	global.DB.Where("id in ?", ToBeDeletedIDs).Find(&records)
	for i := range records {
		records[i].Deleter = &paramIn.Deleter
	}
	err := global.DB.Where("id in ?", ToBeDeletedIDs).Delete(&records).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (*disassembly) GetList(paramIn dto.DisassemblyList) response.List {
	db := global.DB.Model(&model.Disassembly{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if paramIn.NameInclude != "" {
		db = db.Where("name like ?", "%"+paramIn.NameInclude+"%")
	}

	if paramIn.ProjectID > 0 {
		db = db.Where("project_id = ?", paramIn.ProjectID)
	}

	if paramIn.SuperiorID > 0 {
		db = db.Where("superior_id = ?", paramIn.SuperiorID)
	}

	if paramIn.Level > 0 {
		db = db.Where("level = ?", paramIn.Level)
	}

	if paramIn.LevelGte != nil && *paramIn.LevelGte >= 0 {
		db = db.Where("level >= ?", paramIn.LevelGte)
	}

	if paramIn.LevelLte != nil && *paramIn.LevelLte >= 0 {
		db = db.Where("level <= ?", paramIn.LevelLte)
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := paramIn.SortingInput.OrderBy
	desc := paramIn.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Disassembly{}, orderBy)
		if !exists {
			return response.FailForList(util.ErrorSortingFieldDoesNotExist)
		}
		//如果要求降序排列
		if desc == true {
			db = db.Order(orderBy + " desc")
		} else { //如果没有要求排序方式
			db = db.Order(orderBy)
		}
	}

	//limit
	page := 1
	if paramIn.PagingInput.Page > 0 {
		page = paramIn.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if paramIn.PagingInput.PageSize > 0 &&
		paramIn.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = paramIn.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []dto.DisassemblyOutput
	db.Model(&model.Disassembly{}).Find(&data)

	if len(data) == 0 {
		return response.FailForList(util.ErrorRecordNotFound)
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetTotalNumberOfPages(numberOfRecords, pageSize)

	return response.List{
		Data: data,
		Paging: &dto.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
