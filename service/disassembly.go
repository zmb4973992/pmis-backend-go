package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type DisassemblyGet struct {
	Id int64
}

type DisassemblyTree struct {
	UserId    int64
	ProjectId int64 `json:"project_id" binding:"required"`
}

type DisassemblyCreate struct {
	UserId                int64
	Name                  string  `json:"name" binding:"required"`        //拆解项名称
	Weight                float64 `json:"weight" binding:"required"`      //权重
	SuperiorId            int64   `json:"superior_id" binding:"required"` //上级拆解项id
	ImportedIdFromOldPmis int64   `json:"imported_id_from_old_pmis,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DisassemblyUpdate struct {
	UserId int64
	Id     int64

	Name   *string  `json:"name"`   //拆解项名称
	Weight *float64 `json:"weight"` //权重
}

type DisassemblyDelete struct {
	UserId int64
	Id     int64
}

type DisassemblyGetList struct {
	list.Input
	ProjectId  int64 `json:"project_id"`
	SuperiorId int64 `json:"superior_id"`
	Level      int   `json:"level"`
	LevelGte   *int  `json:"level_gte"`
	LevelLte   *int  `json:"level_lte"`
}

//以下为出参

type DisassemblyOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`

	Name       *string  `json:"name"`        //名称
	ProjectId  *int64   `json:"project_id"`  //所属项目id
	Level      *int     `json:"level"`       //层级
	Weight     *float64 `json:"weight"`      //权重
	SuperiorId *int64   `json:"superior_id"` //上级拆解项id
}

type DisassemblyTreeOutput struct {
	Name     *string                 `json:"name"`
	Id       int64                   `json:"id"`
	Level    int                     `json:"level"`
	Children []DisassemblyTreeOutput `json:"children" gorm:"-"`
}

func (d *DisassemblyGet) Get() (output *DisassemblyOutput, errCode int) {
	err := global.DB.Model(model.Disassembly{}).
		Where("id = ?", d.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}
	return output, util.Success
}

func (d *DisassemblyTree) Tree() (outputs []DisassemblyTreeOutput, errCode int) {
	//根据project_id获取disassembly_id
	var disassemblyId int64
	res := global.DB.Model(model.Disassembly{}).
		Where("project_id = ?", d.ProjectId).
		Where("level = 1").
		Select("id").
		Find(&disassemblyId)
	if res.RowsAffected == 0 {
		return nil, util.ErrorRecordNotFound
	}

	//第一轮查找，查询条件为id
	res = global.DB.Model(model.Disassembly{}).
		Where("id = ?", disassemblyId).
		Find(&outputs)
	if res.RowsAffected == 0 {
		return nil, util.ErrorRecordNotFound
	}

	//第二轮及以后的查找，查询条件为superior_id
	for i := range outputs {
		outputs[i].Children = getDisassemblyTree(outputs[i].Id)
	}

	return outputs, util.Success
}

func getDisassemblyTree(superiorId int64) []DisassemblyTreeOutput {
	var result []DisassemblyTreeOutput
	res := global.DB.Model(model.Disassembly{}).
		Where("superior_id = ?", superiorId).
		Find(&result)
	if res.RowsAffected == 0 {
		return nil
	}
	for i := range result {
		result[i].Children = getDisassemblyTree(result[i].Id)
	}
	return result
}

func (d *DisassemblyCreate) Create() (errCode int) {
	var paramOut model.Disassembly
	if d.UserId > 0 {
		paramOut.Creator = &d.UserId
	}

	paramOut.Name = &d.Name
	paramOut.Weight = &d.Weight
	paramOut.SuperiorId = &d.SuperiorId
	paramOut.ImportedIdFromOldPmis = &d.ImportedIdFromOldPmis

	//根据上级拆解id，找到项目id和层级
	var superiorDisassembly model.Disassembly
	err := global.DB.Where("id = ?", d.SuperiorId).
		First(&superiorDisassembly).Error
	if err != nil {
		return util.ErrorWrongSuperiorInformation
	}

	if superiorDisassembly.ProjectId == nil || superiorDisassembly.Level == nil {
		return util.ErrorWrongSuperiorInformation
	}

	paramOut.ProjectId = superiorDisassembly.ProjectId
	level := *superiorDisassembly.Level + 1
	paramOut.Level = &level

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}

	var operationType model.DictionaryType
	err = global.DB.Where("name = ?", "操作类型").
		First(&operationType).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}

	var create model.DictionaryDetail
	err = global.DB.Where("name = ?", "添加").
		First(&create).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}

	var param OperationLogCreate
	param.Creator = d.UserId
	param.Operator = d.UserId
	param.ProjectId = *superiorDisassembly.ProjectId

	param.Date = time.Now().Format("2006-01-02")
	param.OperationType = create.Id
	param.Detail = "添加了一条项目拆解记录：" + d.Name
	param.Create()

	return util.Success
}

func (d *DisassemblyUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	paramOut["last_modifier"] = d.UserId

	if d.Name != nil {
		if *d.Name != "" {
			paramOut["name"] = d.Name
		} else {
			return util.ErrorInvalidJSONParameters
		}
	}

	if d.Weight != nil {
		if *d.Weight >= 0 {
			paramOut["weight"] = d.Weight
		} else {
			return util.ErrorInvalidJSONParameters
		}
	}

	err := global.DB.Model(&model.Disassembly{}).
		Where("id = ?", d.Id).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	var disassembly model.Disassembly
	err = global.DB.Where("id = ?", d.Id).
		First(&disassembly).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	var operationType model.DictionaryType
	err = global.DB.Where("name = ?", "操作类型").
		First(&operationType).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	var update model.DictionaryDetail
	err = global.DB.Where("name = ?", "修改").
		First(&update).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	var param OperationLogCreate
	param.Creator = d.UserId
	param.Operator = d.UserId
	if disassembly.ProjectId != nil {
		param.ProjectId = *disassembly.ProjectId
	}
	param.Date = time.Now().Format("2006-01-02")
	param.OperationType = update.Id
	param.Detail = "修改了一条项目拆解记录：" + *disassembly.Name
	param.Create()

	//如果修改的字段里包含weight，就要更新相关的进度
	for i := range paramOut {
		if i == "weight" {
			//获取进度类型在字典类型表中的值
			var progressTypeId int64
			err = global.DB.Model(&model.DictionaryType{}).
				Where("name = '进度类型'").
				Select("id").
				First(&progressTypeId).Error
			if err != nil {
				return util.ErrorFailToCalculateSelfAndSuperiorProgress
			}
			//获取进度类型在字典详情表中的值
			var progressDetailIds []int64
			err = global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_id = ?", progressTypeId).
				Select("id").
				Find(&progressDetailIds).Error
			if err != nil {
				return util.ErrorFailToCalculateSelfAndSuperiorProgress
			}
			//更新自身和所有上级的进度
			for _, v := range progressDetailIds {
				//更新自身进度
				err = util.UpdateOwnProgress(d.Id, v, d.UserId)
				if err != nil {
					return util.ErrorFailToCalculateSelfProgress
				}
				//更新所有上级的进度
				err = util.UpdateProgressOfSuperiors(d.Id, v, d.UserId)
				if err != nil {
					return util.ErrorFailToCalculateSuperiorProgress
				}
			}
		}
	}

	return util.Success
}

func (d *DisassemblyDelete) Delete() (errCode int) {
	//先找到所有的上级id(如果放在删除后执行，就找不到上级id了)
	//这里的上级id需要更新进度
	superiorIds := util.GetSuperiorIds(d.Id)

	//先找到所有的下级id(如果放在删除后执行，就找不到上级id了)
	//这里的下级id(包括自己)需要删除进度
	inferiorIds := util.GetInferiorIds(d.Id)
	ToBeDeletedIds := append([]int64{d.Id}, inferiorIds...)

	var disassembly model.Disassembly
	err := global.DB.Where("id = ?", d.Id).
		First(&disassembly).Error
	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	err = global.DB.Where("id in ?", ToBeDeletedIds).
		Delete(&model.Disassembly{}).Error
	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	var operationType model.DictionaryType
	err = global.DB.Where("name = ?", "操作类型").
		First(&operationType).Error
	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	var deleting model.DictionaryDetail
	err = global.DB.Where("name = ?", "删除").
		First(&deleting).Error
	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	var param OperationLogCreate
	param.Creator = d.UserId
	param.Operator = d.UserId
	param.ProjectId = *disassembly.ProjectId
	param.Date = time.Now().Format("2006-01-02")
	param.OperationType = deleting.Id
	param.Detail = "删除了一条项目拆解记录：" + *disassembly.Name
	param.Create()

	//获取进度类型在字典类型表中的值
	var progressType model.DictionaryType
	err = global.DB.
		Where("name = '进度类型'").
		First(&progressType).Error
	if err != nil {
		return util.ErrorFailToCalculateSelfAndSuperiorProgress
	}

	//获取所有进度类型在字典详情表中的值
	var allProgressTypes []model.DictionaryDetail
	err = global.DB.
		Where("dictionary_type_id = ?", progressType.Id).
		Find(&allProgressTypes).Error
	if err != nil {
		return util.ErrorFailToCalculateSelfAndSuperiorProgress
	}

	//删除自身和所有下级的进度
	for _, v := range ToBeDeletedIds {
		err = global.DB.Where("disassembly_id = ?", v).
			Delete(&model.Progress{}).Error
		if err != nil {
			return util.ErrorFailToDeleteRecord
		}
	}

	//更新所有上级的进度
	for i := range allProgressTypes {
		for j := range superiorIds {
			err = util.UpdateOwnProgress(superiorIds[j], allProgressTypes[i].Id, d.UserId)
			if err != nil {
				return util.ErrorFailToCalculateSuperiorProgress
			}
		}
	}

	return util.Success
}

func (d *DisassemblyGetList) GetList() (
	outputs []DisassemblyOutput, errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.Disassembly{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if d.ProjectId > 0 {
		db = db.Where("project_id = ?", d.ProjectId)
	}

	if d.SuperiorId > 0 {
		db = db.Where("superior_id = ?", d.SuperiorId)
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
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Disassembly{}, orderBy)
		if !exists {
			return nil, util.ErrorSortingFieldDoesNotExist, nil
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
	pageSize := global.Config.Paging.DefaultPageSize
	if d.PagingInput.PageSize != nil && *d.PagingInput.PageSize >= 0 &&
		*d.PagingInput.PageSize <= global.Config.Paging.MaxPageSize {
		pageSize = *d.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.Disassembly{}).Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return outputs,
		util.Success,
		&list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		}
}
