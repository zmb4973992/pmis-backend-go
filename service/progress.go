package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
	"time"
)

//以下为入参

type ProgressGet struct {
	Id int64
}

type ProgressCreate struct {
	UserId int64

	DisassemblyId int64    `json:"disassembly_id" binding:"required"`
	Date          string   `json:"date" binding:"required"`
	Type          int64    `json:"type" binding:"required"`
	Value         *float64 `json:"value" binding:"required"`
	Remarks       string   `json:"remarks,omitempty"`
}

type ProgressUpdate struct {
	UserId int64
	Id     int64

	Date    *string  `json:"date"`
	Type    *int64   `json:"type"`
	Value   *float64 `json:"value"`
	Remarks *string  `json:"remarks"`
}

type ProgressUpdateByProjectId struct {
	UserId    int64
	ProjectId int64 `json:"project_id,omitempty"`
}

type ProgressDelete struct {
	UserId int64
	Id     int64
}

type ProgressGetList struct {
	list.Input
	ProjectId     int64    `json:"project_id"`
	DisassemblyId int64    `json:"disassembly_id"`
	DateGte       string   `json:"date_gte,omitempty"`
	DateLte       string   `json:"date_lte,omitempty"`
	Type          int64    `json:"type,omitempty"`
	TypeName      string   `json:"type_name,omitempty"`
	TypeIn        []int64  `json:"type_in"`
	ValueGte      *float64 `json:"value_gte"`
	ValueLte      *float64 `json:"value_lte"`
	DataSource    int64    `json:"data_source"`
}

//以下为出参

type ProgressOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`

	DisassemblyId       *int64             `json:"-"`
	DisassemblyExternal *DisassemblyOutput `json:"disassembly" gorm:"-"`

	Date               *string                 `json:"date"`
	Type               *int64                  `json:"-"`
	TypeExternal       *DictionaryDetailOutput `json:"type" gorm:"-"`
	Value              *float64                `json:"value"`
	Remarks            *string                 `json:"remarks"`
	DataSource         *int64                  `json:"-"`
	DataSourceExternal *DictionaryDetailOutput `json:"data_source" gorm:"-"`
}

func (p *ProgressGet) Get() (output *ProgressOutput, errCode int) {
	err := global.DB.Model(model.Progress{}).
		Where("id = ?", p.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	//默认格式为这样的string：2019-11-01T00:00:00Z，需要取年月日(前9位)
	if output.Date != nil {
		temp := *output.Date
		*output.Date = temp[:10]
	}

	//查dictionary_item表
	{
		if output.Type != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.Type).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				output.TypeExternal = &record
			}
		}

		if output.DataSource != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.DataSource).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				output.DataSourceExternal = &record
			}
		}
	}

	return output, util.Success
}

func (p *ProgressCreate) Create() (errCode int) {
	var paramOut model.Progress

	if p.UserId > 0 {
		paramOut.Creator = &p.UserId
	}

	paramOut.DisassemblyId = &p.DisassemblyId

	date, err := time.Parse("2006-01-02", p.Date)
	if err != nil {
		return util.ErrorInvalidDateFormat
	}
	paramOut.Date = &date
	paramOut.Type = &p.Type
	paramOut.Value = p.Value

	//找到"人工填写"的dictionary_item值
	var dataSource model.DictionaryDetail
	err = global.DB.Where("name = '人工填写'").
		First(&dataSource).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}
	paramOut.DataSource = &dataSource.Id

	if p.Remarks != "" {
		paramOut.Remarks = &p.Remarks
	}

	//找到“人工填写”在字典详情表的id
	var dictionaryItem model.DictionaryDetail
	err = global.DB.Where("name = '人工填写'").
		First(&dictionaryItem).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}

	//检查数据库是否已有相同日期、相同类型的记录
	res := global.DB.FirstOrCreate(&paramOut, model.Progress{
		DisassemblyId: &p.DisassemblyId,
		Date:          &date,
		Type:          &p.Type,
	})

	if res.Error != nil {
		return util.ErrorFailToCreateRecord
	}

	if res.RowsAffected == 0 {
		return util.ErrorDuplicateRecord
	}

	var disassembly model.Disassembly
	err = global.DB.Where("id = ?", p.DisassemblyId).
		First(&disassembly).Error
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
	param.Creator = p.UserId
	param.Operator = p.UserId
	param.ProjectId = *disassembly.ProjectId
	param.Date = time.Now().Format("2006-01-02")
	param.OperationType = create.Id
	param.Detail = "添加了" + *disassembly.Name + "的进度"
	param.Create()

	//更新所有上级的进度
	err = util.UpdateProgressOfSuperiors(p.DisassemblyId, p.Type, p.UserId)

	if err != nil {
		return util.ErrorFailToCalculateSuperiorProgress
	}

	return util.Success
}

func (p *ProgressUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	if p.UserId > 0 {
		paramOut["last_modifier"] = p.UserId
	}

	if p.Date != nil {
		if *p.Date != "" {
			var err error
			paramOut["date"], err = time.Parse("2006-01-02", *p.Date)
			if err != nil {
				return util.ErrorInvalidJSONParameters
			}
		} else {
			return util.ErrorInvalidJSONParameters
		}
	}

	if p.Type != nil {
		if *p.Type > 0 {
			paramOut["type"] = p.Type
		} else {
			return util.ErrorInvalidJSONParameters
		}
	}

	if p.Value != nil {
		if *p.Value >= 0 {
			paramOut["value"] = p.Value
		} else {
			return util.ErrorInvalidJSONParameters
		}
	}

	if p.Remarks != nil {
		if *p.Remarks != "" {
			paramOut["remarks"] = *p.Remarks
		} else {
			paramOut["remarks"] = nil
		}
	}

	//找到“人工填写”在字典详情表的id
	var dataSource model.DictionaryDetail
	err := global.DB.Where("name = '人工填写'").
		First(&dataSource).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}
	paramOut["data_source"] = dataSource.Id

	//找到待更新的这条记录
	var progress model.Progress
	err = global.DB.Where("id = ?", p.Id).
		First(&progress).Error
	if err != nil {
		return util.ErrorFailToCalculateSuperiorProgress
	}

	//如果修改了date或type，意味着可能有重复记录，需要进行判断
	if p.Date != nil || p.Type != nil {
		//从数据库找出相同拆解id、相同日期、相同类型的记录
		var tempProgressIds []int64
		tempDate, err1 := time.Parse("2006-01-02", *p.Date)
		if err1 != nil {
			return util.ErrorInvalidDateFormat
		}
		global.DB.Model(&model.Progress{}).
			Where(&model.Progress{
				DisassemblyId: progress.DisassemblyId,
				Date:          &tempDate,
				Type:          p.Type,
			}).
			Select("id").
			Find(&tempProgressIds)
		//如果数据库有记录、且待修改的progressId不在数据库记录的progressIds里面，说明是新的记录
		//则不允许修改
		if len(tempProgressIds) > 0 && !util.IsInSlice(p.Id, tempProgressIds) {
			return util.ErrorDuplicateRecord
		}
	}

	//更新记录
	err = global.DB.Model(&model.Progress{}).
		Where("id = ?", p.Id).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	var disassembly model.Disassembly
	err = global.DB.Where("id = ?", progress.DisassemblyId).
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
	param.Creator = p.UserId
	param.Operator = p.UserId
	if disassembly.ProjectId != nil {
		param.ProjectId = *disassembly.ProjectId
	}
	param.Date = time.Now().Format("2006-01-02")
	param.OperationType = update.Id
	param.Detail = "修改了" + *disassembly.Name + "的进度"
	param.Create()

	//更新所有上级的进度
	if progress.DisassemblyId != nil {
		//如果传入了type值(意味着type值可能从a改成b，同时影响a、b)，就准备更新所有的进度类型
		if p.Type != nil {
			//找出”进度类型“的dictionary_type值
			var progressTypeIdInDictionaryType model.DictionaryType
			err = global.DB.Where("name = '进度类型'").
				First(&progressTypeIdInDictionaryType).
				Error
			if err != nil {
				return util.ErrorFailToCalculateSuperiorProgress
			}

			//找出"进度类型"的dictionary_item值，准备遍历
			var progressTypeIds []int64
			global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_id = ?", progressTypeIdInDictionaryType.Id).
				Select("id").Find(&progressTypeIds)

			for _, v := range progressTypeIds {
				err = util.UpdateProgressOfSuperiors(*progress.DisassemblyId, v, p.UserId)
				if err != nil {
					return util.ErrorFailToCalculateSuperiorProgress
				}
			}
		} else { //如果没有传入type值(意味着记录的type值不变)，则只更新原来的进度类型
			err = util.UpdateProgressOfSuperiors(*progress.DisassemblyId, *progress.Type, p.UserId)
			if err != nil {
				return util.ErrorFailToCalculateSuperiorProgress
			}
		}
	}

	return util.Success
}

func (p *ProgressUpdateByProjectId) UpdateByProjectId() error {
	var progressType model.DictionaryType
	err := global.DB.Where("name = '进度类型'").
		First(&progressType).Error
	if err != nil {
		return err
	}

	var allProgressTypes []model.DictionaryDetail
	err = global.DB.Where("dictionary_type_id = ?", progressType.Id).
		Find(&allProgressTypes).Error

	var disassemblies []model.Disassembly
	global.DB.
		Where("project_id = ?", p.ProjectId).
		Order("level desc").
		Find(&disassemblies)

	for i := range allProgressTypes {
		for j := range disassemblies {
			err = util.UpdateOwnProgress(disassemblies[j].Id, allProgressTypes[i].Id, p.UserId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *ProgressDelete) Delete() (errCode int) {
	//先找到记录，这样参数才能获得值、触发钩子函数，再删除记录
	var progress model.Progress
	err := global.DB.Where("id = ?", p.Id).
		First(&progress).
		Delete(&progress).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	var disassembly model.Disassembly
	err = global.DB.Where("id = ?", progress.DisassemblyId).
		First(&disassembly).Error
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
	param.Creator = p.UserId
	param.Operator = p.UserId
	param.ProjectId = *disassembly.ProjectId
	param.Date = time.Now().Format("2006-01-02")
	param.OperationType = deleting.Id
	param.Detail = "删除了" + *disassembly.Name + "的进度"
	param.Create()

	//更新所有上级的进度
	if progress.DisassemblyId != nil && progress.Type != nil {
		err = util.UpdateProgressOfSuperiors(*progress.DisassemblyId, *progress.Type, p.UserId)
		if err != nil {
			return util.ErrorFailToCalculateSuperiorProgress
		}
	}

	return util.Success
}

func (p *ProgressGetList) GetList() (
	outputs []ProgressOutput, errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.Progress{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where
	if p.ProjectId > 0 {
		var disassemblyId int64
		err := global.DB.Model(&model.Disassembly{}).
			Where("project_id = ?", p.ProjectId).
			Where("superior_id is null").
			Select("id").First(&disassemblyId).Error
		if err != nil {
			return nil, util.ErrorRecordNotFound, nil
		}
		db = db.Where("disassembly_id = ?", disassemblyId)
	}
	if p.DisassemblyId > 0 {
		db = db.Where("disassembly_id = ?", p.DisassemblyId)
	}

	if p.DateGte != "" {
		date, err := time.Parse("2006-01-02", p.DateGte)
		if err != nil {
			return nil, util.ErrorInvalidJSONParameters, nil
		}
		db = db.Where("date >= ?", date)
	}

	if p.DateLte != "" {
		date, err := time.Parse("2006-01-02", p.DateLte)
		if err != nil {
			return nil, util.ErrorInvalidJSONParameters, nil
		}
		db = db.Where("date <= ?", date)
	}

	if p.Type > 0 {
		db = db.Where("type = ?", p.Type)
	}

	if p.TypeName != "" {
		var typeId int64
		err := global.DB.Model(&model.DictionaryDetail{}).
			Where("name = ?", p.TypeName).
			Select("id").
			First(&typeId).Error
		if err != nil {
			return nil, util.ErrorRecordNotFound, nil
		}
		db = db.Where("type = ?", typeId)

	}

	if len(p.TypeIn) > 0 {
		db = db.Where("type in ?", p.TypeIn)
	}

	if p.ValueGte != nil {
		db = db.Where("value >= ?", *p.ValueGte)
	}

	if p.ValueLte != nil {
		db = db.Where("value <= ?", *p.ValueLte)
	}

	if p.DataSource > 0 {
		db = db.Where("data_source = ?", p.DataSource)
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := p.SortingInput.OrderBy
	desc := p.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Progress{}, orderBy)
		if !exists {
			return nil, util.ErrorSortingFieldDoesNotExist, nil
		}
		//如果要求降序排列
		if desc == true {
			db = db.Order(orderBy + " desc")
		} else {
			//如果没有要求排序方式
			db = db.Order(orderBy)
		}
	}

	//limit
	page := 1
	if p.PagingInput.Page > 0 {
		page = p.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if p.PagingInput.PageSize != nil && *p.PagingInput.PageSize >= 0 &&
		*p.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *p.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.Progress{}).Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	//查拆解信息
	if p.DisassemblyId > 0 {
		var record DisassemblyOutput
		res := global.DB.Model(&model.Disassembly{}).
			Where("id = ?", p.DisassemblyId).
			Limit(1).
			Find(&record)
		if res.RowsAffected > 0 {
			for i := range outputs {
				outputs[i].DisassemblyExternal = &record
			}
		}
	}

	//一次性查出需要用的进度类型，避免多次重复查询
	var progressType model.DictionaryType
	err := global.DB.Where("name = '进度类型'").
		First(&progressType).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound, nil
	}

	var planned DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", progressType.Id).
		Where("name = '计划进度'").
		First(&planned).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound, nil
	}

	var forecasted DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", progressType.Id).
		Where("name = '预测进度'").
		First(&forecasted).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound, nil
	}

	var actual DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", progressType.Id).
		Where("name = '实际进度'").
		First(&actual).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound, nil
	}

	var dataSource model.DictionaryType
	err = global.DB.Where("name = '进度的数据来源'").
		First(&dataSource).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound, nil
	}

	var systemCalculation DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", dataSource.Id).
		Where("name = '系统计算'").
		First(&systemCalculation).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound, nil
	}

	var manualFilling DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", dataSource.Id).
		Where("name = '人工填写'").
		First(&manualFilling).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound, nil
	}

	for i := range outputs {

		//处理日期格式
		if outputs[i].Date != nil {
			temp := *outputs[i].Date
			*outputs[i].Date = temp[:10]
		}

		//查dictionary_item表
		{
			if outputs[i].Type != nil {
				if *outputs[i].Type == planned.Id {
					outputs[i].TypeExternal = &planned
				} else if *outputs[i].Type == forecasted.Id {
					outputs[i].TypeExternal = &forecasted
				} else if *outputs[i].Type == actual.Id {
					outputs[i].TypeExternal = &actual
				} else {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *outputs[i].Type).
						Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].TypeExternal = &record
					}
				}
			}

			if outputs[i].DataSource != nil {
				if *outputs[i].DataSource == systemCalculation.Id {
					outputs[i].DataSourceExternal = &systemCalculation
				} else if *outputs[i].DataSource == manualFilling.Id {
					outputs[i].DataSourceExternal = &manualFilling
				} else {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *outputs[i].DataSource).
						Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						outputs[i].DataSourceExternal = &record
					}
				}
			}
		}
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
