package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type DisassemblyGet struct {
	SnowID int64
}

type DisassemblyTree struct {
	Creator       int64
	LastModifier  int64
	ProjectSnowID int64 `json:"project_snow_id" binding:"required"`
}

type DisassemblyCreate struct {
	Creator      int64
	LastModifier int64

	Name           string  `json:"name" binding:"required"`             //拆解项名称
	Weight         float64 `json:"weight" binding:"required"`           //权重
	SuperiorSnowID int64   `json:"superior_snow_id" binding:"required"` //上级拆解项ID
}

type DisassemblyCreateInBatches struct {
	Param []DisassemblyCreate `json:"param"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DisassemblyUpdate struct {
	LastModifier int64
	SnowID       int64

	Name   *string  `json:"name"`   //拆解项名称
	Weight *float64 `json:"weight"` //权重
}

type DisassemblyDelete struct {
	SnowID int64
}

type DisassemblyDeleteWithInferiors struct {
	SnowID int64
}

type DisassemblyGetList struct {
	list.Input
	NameInclude string `json:"name_include,omitempty"`

	ProjectSnowID  int64 `json:"project_snow_id"`
	SuperiorSnowID int64 `json:"superior_snow_id"`
	Level          int   `json:"level"`
	LevelGte       *int  `json:"level_gte"`
	LevelLte       *int  `json:"level_lte"`
}

//以下为出参

type DisassemblyOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	SnowID       int64  `json:"snow_id"`

	Name           *string  `json:"name"`             //名称
	ProjectSnowID  *int64   `json:"project_snow_id"`  //所属项目SnowID
	Level          *int     `json:"level"`            //层级
	Weight         *float64 `json:"weight"`           //权重
	SuperiorSnowID *int64   `json:"superior_snow_id"` //上级拆解项SnowID
}

type DisassemblyTreeOutput struct {
	Name     *string                 `json:"name"`
	SnowID   int64                   `json:"snow_id"`
	Level    int                     `json:"level"`
	Children []DisassemblyTreeOutput `json:"children" gorm:"-"`
}

func (d *DisassemblyGet) Get() response.Common {
	var result DisassemblyOutput
	err := global.DB.Model(model.Disassembly{}).
		Where("snow_id = ?", d.SnowID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func treeRecursion(superiorSnowID int64) []DisassemblyTreeOutput {
	var result []DisassemblyTreeOutput
	res := global.DB.Model(model.Disassembly{}).
		Where("superior_snow_id = ?", superiorSnowID).Find(&result)
	if res.RowsAffected == 0 {
		return nil
	}
	for i := range result {
		result[i].Children = treeRecursion(result[i].SnowID)
	}
	return result
}

func (d *DisassemblyTree) Tree() response.Common {
	//根据project_id获取disassembly_id
	var disassemblySnowID int64
	res := global.DB.Model(model.Disassembly{}).Select("snow_id").
		Where("project_snow_id = ?", d.ProjectSnowID).Where("level = 1").
		Find(&disassemblySnowID)
	if res.RowsAffected == 0 {
		return response.Failure(util.ErrorRecordNotFound)
	}

	//第一轮查找，查询条件为id
	var result []DisassemblyTreeOutput
	res = global.DB.Model(model.Disassembly{}).
		Where("snow_id = ?", disassemblySnowID).Find(&result)
	if res.RowsAffected == 0 {
		return response.Failure(util.ErrorRecordNotFound)
	}

	//第二轮及以后的查找，查询条件为superior_id
	for i := range result {
		result[i].Children = treeRecursion(result[i].SnowID)
	}

	return response.SuccessWithData(result)
}

func (d *DisassemblyCreate) Create() response.Common {
	var paramOut model.Disassembly
	if d.Creator > 0 {
		paramOut.Creator = &d.Creator
	}

	if d.LastModifier > 0 {
		paramOut.LastModifier = &d.LastModifier
	}

	paramOut.Name = &d.Name
	paramOut.Weight = &d.Weight
	paramOut.SuperiorSnowID = &d.SuperiorSnowID

	//根据上级拆解id，找到项目id和层级
	var superiorDisassembly model.Disassembly
	err := global.DB.Where("snow_id = ?", d.SuperiorSnowID).First(&superiorDisassembly).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	if superiorDisassembly.ProjectSnowID == nil || superiorDisassembly.Level == nil {
		return response.Failure(util.ErrorWrongSuperiorInformation)
	}

	paramOut.ProjectSnowID = superiorDisassembly.ProjectSnowID
	level := *superiorDisassembly.Level + 1
	paramOut.Level = &level

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (d *DisassemblyUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if d.LastModifier > 0 {
		paramOut["last_modifier"] = d.LastModifier
	}

	if d.Name != nil {
		if *d.Name != "" {
			paramOut["name"] = d.Name
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if d.Weight != nil {
		if *d.Weight >= 0 {
			paramOut["weight"] = d.Weight
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Disassembly{}).Where("snow_id = ?", d.SnowID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	//如果修改的字段里包含weight，就要更新相关的进度
	for i := range paramOut {
		if i == "weight" {
			//获取进度类型在字典类型表中的值
			var progressTypeSnowID int64
			err = global.DB.Model(&model.DictionaryType{}).
				Where("name = '进度类型'").Select("snow_id").First(&progressTypeSnowID).Error
			if err != nil {
				return response.Failure(util.ErrorFailToCalculateSelfAndSuperiorProgress)
			}
			//获取进度类型在字典详情表中的值
			var progressItemSnowIDs []int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_snow_id = ?", progressTypeSnowID).
				Select("snow_id").Find(&progressItemSnowIDs).Error
			if err != nil {
				return response.Failure(util.ErrorFailToCalculateSelfAndSuperiorProgress)
			}
			//更新自身和所有上级的进度
			for _, v := range progressItemSnowIDs {
				//更新自身进度
				err = util.UpdateSelfProgress(d.SnowID, v)
				if err != nil {
					global.SugaredLogger.Errorln(err)
					return response.Failure(util.ErrorFailToCalculateSelfProgress)
				}
				//更新所有上级的进度
				err = util.UpdateProgressOfSuperiors(d.SnowID, v)
				if err != nil {
					global.SugaredLogger.Errorln(err)
					return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
				}
			}
		}
	}

	return response.Success()
}

func (d *DisassemblyDelete) Delete() response.Common {
	//先找到记录，这样参数才能获得值、触发钩子函数，再删除记录
	var record model.Disassembly
	err := global.DB.Where("snow_id = ?", d.SnowID).Find(&record).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	//获取进度类型在字典类型表中的值
	var progressTypeSnowID int
	err = global.DB.Model(&model.DictionaryType{}).
		Where("name = '进度类型'").Select("snow_id").First(&progressTypeSnowID).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCalculateSelfAndSuperiorProgress)
	}

	//获取进度类型在字典详情表中的值
	var progressItemSnowIDs []int64
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_snow_id = ?", progressTypeSnowID).
		Select("snow_id").Find(&progressItemSnowIDs).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCalculateSelfAndSuperiorProgress)
	}

	//更新所有上级的进度(删除自身进度的任务放在删除的钩子函数里)
	for _, v := range progressItemSnowIDs {
		err = util.UpdateProgressOfSuperiors(d.SnowID, v)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
		}
	}

	return response.Success()
}

func (d *DisassemblyDeleteWithInferiors) DeleteWithInferiors() response.Common {
	inferiorSnowIDs := util.GetInferiorSnowIDs(d.SnowID)
	ToBeDeletedSnowIDs := append([]int64{d.SnowID}, inferiorSnowIDs...)

	//先找到记录，这样参数才能获得值、触发钩子函数，再删除记录
	var disassemblies []model.Disassembly
	err := global.DB.Where("snow_id in ?", ToBeDeletedSnowIDs).
		Find(&disassemblies).Delete(&disassemblies).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	return response.Success()
}

func (d *DisassemblyGetList) GetList() response.List {
	db := global.DB.Model(&model.Disassembly{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if d.NameInclude != "" {
		db = db.Where("name like ?", "%"+d.NameInclude+"%")
	}

	if d.ProjectSnowID > 0 {
		db = db.Where("project_snow_id = ?", d.ProjectSnowID)
	}

	if d.SuperiorSnowID > 0 {
		db = db.Where("superior_snow_id = ?", d.SuperiorSnowID)
	}

	if d.Level > 0 {
		db = db.Where("level = ?", d.Level)
	}

	if d.LevelGte != nil && *d.LevelGte >= 0 {
		db = db.Where("level >= ?", d.LevelGte)
	}

	if d.LevelLte != nil && *d.LevelLte >= 0 {
		db = db.Where("level <= ?", d.LevelLte)
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := d.SortingInput.OrderBy
	desc := d.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("snow_id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Disassembly{}, orderBy)
		if !exists {
			return response.FailureForList(util.ErrorSortingFieldDoesNotExist)
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
	if d.PagingInput.Page > 0 {
		page = d.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if d.PagingInput.PageSize != nil && *d.PagingInput.PageSize >= 0 &&
		*d.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *d.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []DisassemblyOutput
	db.Model(&model.Disassembly{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return response.List{
		Data: data,
		Paging: &list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
