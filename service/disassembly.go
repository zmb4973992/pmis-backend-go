package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type DisassemblyGet struct {
	ID int
}

type DisassemblyTree struct {
	Creator      int
	LastModifier int
	ProjectID    int `json:"project_id" binding:"required"`
}

type DisassemblyCreate struct {
	Creator      int
	LastModifier int

	Name       string  `json:"name" binding:"required"`        //拆解项名称
	ProjectID  int     `json:"project_id" binding:"required"`  //所属项目id
	Level      int     `json:"level" binding:"required"`       //层级
	Weight     float64 `json:"weight" binding:"required"`      //权重
	SuperiorID int     `json:"superior_id" binding:"required"` //上级拆解项ID
}

type DisassemblyCreateInBatches struct {
	Param []DisassemblyCreate `json:"param"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type DisassemblyUpdate struct {
	LastModifier int
	ID           int

	Name       *string  `json:"name"`        //拆解项名称
	ProjectID  *int     `json:"project_id"`  //所属项目id
	Level      *int     `json:"level"`       //层级
	Weight     *float64 `json:"weight"`      //权重
	SuperiorID *int     `json:"superior_id"` //上级拆解项ID
}

type DisassemblyDelete struct {
	Deleter int
	ID      int
}

type DisassemblyDeleteWithSubitems struct {
	Deleter int
	ID      int
}

type DisassemblyGetList struct {
	dto.ListInput
	NameInclude string `json:"name_include,omitempty"`

	ProjectID  int  `json:"project_id"`
	SuperiorID int  `json:"superior_id"`
	Level      int  `json:"level"`
	LevelGte   *int `json:"level_gte"`
	LevelLte   *int `json:"level_lte"`
}

//以下为出参

type DisassemblyOutput struct {
	Creator      *int `json:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" gorm:"id"`

	Name       *string  `json:"name" gorm:"name"`               //名称
	ProjectID  *int     `json:"project_id" gorm:"project_id"`   //所属项目id
	Level      *int     `json:"level" gorm:"level"`             //层级
	Weight     *float64 `json:"weight" gorm:"weight"`           //权重
	SuperiorID *int     `json:"superior_id" gorm:"superior_id"` //上级拆解项id
}

type DisassemblyTreeOutput struct {
	Name     *string                 `json:"name"`
	ID       int                     `json:"id"`
	Level    int                     `json:"level"`
	Children []DisassemblyTreeOutput `json:"children" gorm:"-"`
}

func (d *DisassemblyGet) Get() response.Common {
	var result DisassemblyOutput
	err := global.DB.Model(model.Disassembly{}).
		Where("id = ?", d.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func treeRecursion(superiorID int) []DisassemblyTreeOutput {
	var result []DisassemblyTreeOutput
	res := global.DB.Model(model.Disassembly{}).
		Where("superior_id = ?", superiorID).Find(&result)
	if res.RowsAffected == 0 {
		return nil
	}
	for i := range result {
		result[i].Children = treeRecursion(result[i].ID)
	}
	return result
}

func (d *DisassemblyTree) Tree() response.Common {
	//根据project_id获取disassembly_id
	var disassemblyID int
	res := global.DB.Model(model.Disassembly{}).Select("id").
		Where("project_id = ?", d.ProjectID).Where("level = 1").
		Find(&disassemblyID)
	if res.RowsAffected == 0 {
		return response.Failure(util.ErrorRecordNotFound)
	}

	//第一轮查找，查询条件为id
	var result []DisassemblyTreeOutput
	res = global.DB.Model(model.Disassembly{}).
		Where("id = ?", disassemblyID).Find(&result)
	if res.RowsAffected == 0 {
		return response.Failure(util.ErrorRecordNotFound)
	}

	//第二轮及以后的查找，查询条件为superior_id
	for i := range result {
		result[i].Children = treeRecursion(result[i].ID)
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

	paramOut.ProjectID = &d.ProjectID

	paramOut.Level = &d.Level

	paramOut.Weight = &d.Weight

	paramOut.SuperiorID = &d.SuperiorID

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (d *DisassemblyCreateInBatches) CreateInBatches() response.Common {
	var paramOut []model.Disassembly
	for i := range d.Param {
		var record model.Disassembly
		if d.Param[i].Creator > 0 {
			record.Creator = &d.Param[i].Creator
		}

		if d.Param[i].LastModifier > 0 {
			record.LastModifier = &d.Param[i].LastModifier
		}

		record.Name = &d.Param[i].Name

		record.Level = &d.Param[i].Level

		record.ProjectID = &d.Param[i].ProjectID

		record.Weight = &d.Param[i].Weight

		record.SuperiorID = &d.Param[i].SuperiorID

		paramOut = append(paramOut, record)
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
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

	if d.ProjectID != nil {
		if *d.ProjectID != 0 {
			paramOut["project_id"] = d.ProjectID
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if d.Level != nil {
		if *d.Level != 0 {
			paramOut["level"] = d.Level
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if d.Weight != nil {
		if *d.Weight != 0 {
			paramOut["weight"] = d.Weight
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if d.SuperiorID != nil {
		if *d.SuperiorID != 0 {
			paramOut["superior_id"] = d.SuperiorID
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.Disassembly{}).Where("id = ?", d.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (d *DisassemblyDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.Disassembly
	global.DB.Where("id = ?", d.ID).Find(&record)
	record.Deleter = &d.Deleter
	err := global.DB.Where("id = ?", d.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (d *DisassemblyDeleteWithSubitems) DeleteWithSubitems() response.Common {
	var ToBeDeletedIDs []int
	ToBeDeletedIDs = append(ToBeDeletedIDs, d.ID)
	//第一轮查找
	var result1 []int
	global.DB.Model(&model.Disassembly{}).Where("superior_id = ?", d.ID).
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
		records[i].Deleter = &d.Deleter
	}
	err := global.DB.Where("id in ?", ToBeDeletedIDs).Delete(&records).Error

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

	if d.ProjectID > 0 {
		db = db.Where("project_id = ?", d.ProjectID)
	}

	if d.SuperiorID > 0 {
		db = db.Where("superior_id = ?", d.SuperiorID)
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
	if d.PagingInput.PageSize > 0 &&
		d.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = d.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

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
