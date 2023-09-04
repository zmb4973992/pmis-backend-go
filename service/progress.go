package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

//以下为入参

type ProgressGet struct {
	ID int64
}

type ProgressCreate struct {
	UserID int64

	DisassemblyID int64    `json:"disassembly_id" binding:"required"`
	Date          string   `json:"date" binding:"required"`
	Type          int64    `json:"type" binding:"required"`
	Value         *float64 `json:"value" binding:"required"`
	Remarks       string   `json:"remarks,omitempty"`
}

type ProgressUpdate struct {
	UserID int64
	ID     int64

	Date    *string  `json:"date"`
	Type    *int64   `json:"type"`
	Value   *float64 `json:"value"`
	Remarks *string  `json:"remarks"`
}

type ProgressUpdateByProjectID struct {
	UserID    int64
	ProjectID int64 `json:"project_id,omitempty"`
}

type ProgressDelete struct {
	UserID int64
	ID     int64
}

type ProgressGetList struct {
	list.Input
	ProjectID     int64    `json:"project_id"`
	DisassemblyID int64    `json:"disassembly_id"`
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
	ID           int64  `json:"id"`

	DisassemblyID       *int64             `json:"-"`
	DisassemblyExternal *DisassemblyOutput `json:"disassembly" gorm:"-"`

	Date               *string                 `json:"date"`
	Type               *int64                  `json:"-"`
	TypeExternal       *DictionaryDetailOutput `json:"type" gorm:"-"`
	Value              *float64                `json:"value"`
	Remarks            *string                 `json:"remarks"`
	DataSource         *int64                  `json:"-"`
	DataSourceExternal *DictionaryDetailOutput `json:"data_source" gorm:"-"`
}

func (p *ProgressGet) Get() response.Common {
	var result ProgressOutput
	err := global.DB.Model(model.Progress{}).
		Where("id = ?", p.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}

	//默认格式为这样的string：2019-11-01T00:00:00Z，需要取年月日(前9位)
	if result.Date != nil {
		temp := *result.Date
		*result.Date = temp[:10]
	}

	//查dictionary_item表
	{
		if result.Type != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.Type).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.TypeExternal = &record
			}
		}

		if result.DataSource != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *result.DataSource).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.DataSourceExternal = &record
			}
		}
	}

	return response.SuccessWithData(result)
}

func (p *ProgressCreate) Create() response.Common {
	var paramOut model.Progress

	if p.UserID > 0 {
		paramOut.Creator = &p.UserID
	}

	paramOut.DisassemblyID = &p.DisassemblyID

	date, err := time.Parse("2006-01-02", p.Date)
	if err != nil {
		return response.Failure(util.ErrorInvalidDateFormat)
	}
	paramOut.Date = &date
	paramOut.Type = &p.Type
	paramOut.Value = p.Value

	//找到"人工填写"的dictionary_item值
	var dataSource model.DictionaryDetail
	err = global.DB.Where("name = '人工填写'").
		First(&dataSource).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOut.DataSource = &dataSource.ID

	if p.Remarks != "" {
		paramOut.Remarks = &p.Remarks
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"UserID", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
	}

	//找到“人工填写”在字典详情表的id
	var dictionaryItem model.DictionaryDetail
	err = global.DB.Where("name = '人工填写'").First(&dictionaryItem).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	//检查数据库是否已有相同日期、相同类型的记录
	res := global.DB.FirstOrCreate(&paramOut, model.Progress{
		DisassemblyID: &p.DisassemblyID,
		Date:          &date,
		Type:          &p.Type,
	})

	if res.Error != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	if res.RowsAffected == 0 {
		return response.Failure(util.ErrorDuplicateRecord)
	}

	var disassembly model.Disassembly
	err = global.DB.Where("id = ?", p.DisassemblyID).
		First(&disassembly).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	var operationType model.DictionaryType
	err = global.DB.Where("name = ?", "操作类型").
		First(&operationType).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	var create model.DictionaryDetail
	err = global.DB.Where("name = ?", "添加").
		First(&create).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	var param OperationLogCreate
	param.Creator = p.UserID
	param.Operator = p.UserID
	param.ProjectID = *disassembly.ProjectID
	param.Date = time.Now().Format("2006-01-02")
	param.OperationType = create.ID
	param.Detail = "添加了" + *disassembly.Name + "的进度"
	param.Create()

	//更新所有上级的进度
	err = util.UpdateProgressOfSuperiors(p.DisassemblyID, p.Type, p.UserID)

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
	}

	return response.Success()
}

func (p *ProgressUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if p.UserID > 0 {
		paramOut["last_modifier"] = p.UserID
	}

	if p.Date != nil {
		if *p.Date != "" {
			var err error
			paramOut["date"], err = time.Parse("2006-01-02", *p.Date)
			if err != nil {
				return response.Failure(util.ErrorInvalidJSONParameters)
			}
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if p.Type != nil {
		if *p.Type > 0 {
			paramOut["type"] = p.Type
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if p.Value != nil {
		if *p.Value >= 0 {
			paramOut["value"] = p.Value
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
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
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	paramOut["data_source"] = dataSource.ID

	//计算有修改值的字段数，分别进行不同处理
	//data_source是自动添加的，也需要排除在外
	paramOutForCounting := util.MapCopy(paramOut,
		"UserID", "CreateAt", "UpdatedAt", "DataSource")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	//找到待更新的这条记录
	var progress model.Progress
	err = global.DB.Where("id = ?", p.ID).
		First(&progress).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
	}

	//如果修改了date或type，意味着可能有重复记录，需要进行判断
	if p.Date != nil || p.Type != nil {
		//从数据库找出相同拆解id、相同日期、相同类型的记录
		var tempProgressIDs []int64
		tempDate, err1 := time.Parse("2006-01-02", *p.Date)
		if err1 != nil {
			return response.Failure(util.ErrorInvalidDateFormat)
		}
		global.DB.Model(&model.Progress{}).Where(&model.Progress{
			DisassemblyID: progress.DisassemblyID,
			Date:          &tempDate,
			Type:          p.Type,
		}).Select("id").Find(&tempProgressIDs)
		//如果数据库有记录、且待修改的progressID不在数据库记录的progressIDs里面，说明是新的记录
		//则不允许修改
		if len(tempProgressIDs) > 0 && !util.IsInSlice(p.ID, tempProgressIDs) {
			return response.Failure(util.ErrorDuplicateRecord)
		}
	}

	//更新记录
	err = global.DB.Model(&model.Progress{}).
		Where("id = ?", p.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	var disassembly model.Disassembly
	err = global.DB.Where("id = ?", progress.DisassemblyID).
		First(&disassembly).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	var operationType model.DictionaryType
	err = global.DB.Where("name = ?", "操作类型").
		First(&operationType).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	var update model.DictionaryDetail
	err = global.DB.Where("name = ?", "修改").
		First(&update).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	var param OperationLogCreate
	param.Creator = p.UserID
	param.Operator = p.UserID
	if disassembly.ProjectID != nil {
		param.ProjectID = *disassembly.ProjectID
	}
	param.Date = time.Now().Format("2006-01-02")
	param.OperationType = update.ID
	param.Detail = "修改了" + *disassembly.Name + "的进度"
	param.Create()

	//更新所有上级的进度
	if progress.DisassemblyID != nil {
		//如果传入了type值(意味着type值可能从a改成b，同时影响a、b)，就准备更新所有的进度类型
		if p.Type != nil {
			//找出”进度类型“的dictionary_type值
			var progressTypeIDInDictionaryType model.DictionaryType
			err = global.DB.Where("name = '进度类型'").First(&progressTypeIDInDictionaryType).
				Error
			if err != nil {
				return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
			}

			//找出"进度类型"的dictionary_item值，准备遍历
			var progressTypeIDs []int64
			global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_id = ?", progressTypeIDInDictionaryType.ID).
				Select("id").Find(&progressTypeIDs)

			for _, v := range progressTypeIDs {
				err = util.UpdateProgressOfSuperiors(*progress.DisassemblyID, v, p.UserID)
				if err != nil {
					global.SugaredLogger.Errorln(err)
					return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
				}
			}
		} else { //如果没有传入type值(意味着记录的type值不变)，则只更新原来的进度类型
			err = util.UpdateProgressOfSuperiors(*progress.DisassemblyID, *progress.Type, p.UserID)
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
			}
		}
	}

	return response.Success()
}

func (p *ProgressUpdateByProjectID) UpdateByProjectID() error {
	var progressType model.DictionaryType
	err := global.DB.Where("name = '进度类型'").
		First(&progressType).Error
	if err != nil {
		return err
	}

	var allProgressTypes []model.DictionaryDetail
	err = global.DB.Where("dictionary_type_id = ?", progressType.ID).
		Find(&allProgressTypes).Error

	var disassemblies []model.Disassembly
	global.DB.
		Where("project_id = ?", p.ProjectID).
		Order("level desc").
		Find(&disassemblies)

	for i := range allProgressTypes {
		for j := range disassemblies {
			err = util.UpdateOwnProgress(disassemblies[j].ID, allProgressTypes[i].ID, p.UserID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *ProgressDelete) Delete() response.Common {
	//先找到记录，这样参数才能获得值、触发钩子函数，再删除记录
	var progress model.Progress
	err := global.DB.Where("id = ?", p.ID).
		First(&progress).Delete(&progress).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	var disassembly model.Disassembly
	err = global.DB.Where("id = ?", progress.DisassemblyID).
		First(&disassembly).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	var operationType model.DictionaryType
	err = global.DB.Where("name = ?", "操作类型").
		First(&operationType).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	var deleting model.DictionaryDetail
	err = global.DB.Where("name = ?", "删除").
		First(&deleting).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	var param OperationLogCreate
	param.Creator = p.UserID
	param.Operator = p.UserID
	param.ProjectID = *disassembly.ProjectID
	param.Date = time.Now().Format("2006-01-02")
	param.OperationType = deleting.ID
	param.Detail = "删除了" + *disassembly.Name + "的进度"
	param.Create()

	//更新所有上级的进度
	if progress.DisassemblyID != nil && progress.Type != nil {
		err = util.UpdateProgressOfSuperiors(*progress.DisassemblyID, *progress.Type, p.UserID)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
		}
	}

	return response.Success()
}

func (p *ProgressGetList) GetList() response.List {
	db := global.DB.Model(&model.Progress{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if p.ProjectID > 0 {
		var disassemblyID int64
		err := global.DB.Model(&model.Disassembly{}).
			Where("project_id = ?", p.ProjectID).
			Where("superior_id is null").
			Select("id").First(&disassemblyID).Error
		if err != nil {
			return response.FailureForList(util.ErrorRecordNotFound)
		}
		db = db.Where("disassembly_id = ?", disassemblyID)
	}
	if p.DisassemblyID > 0 {
		db = db.Where("disassembly_id = ?", p.DisassemblyID)
	}

	if p.DateGte != "" {
		date, err := time.Parse("2006-01-02", p.DateGte)
		if err != nil {
			return response.FailureForList(util.ErrorInvalidJSONParameters)
		}
		db = db.Where("date >= ?", date)
	}

	if p.DateLte != "" {
		date, err := time.Parse("2006-01-02", p.DateLte)
		if err != nil {
			return response.FailureForList(util.ErrorInvalidJSONParameters)
		}
		db = db.Where("date <= ?", date)
	}

	if p.Type > 0 {
		db = db.Where("type = ?", p.Type)
	}

	if p.TypeName != "" {
		var typeID int64
		err := global.DB.Model(&model.DictionaryDetail{}).
			Where("name = ?", p.TypeName).Select("id").
			First(&typeID).Error
		if err != nil {
			return response.FailureForList(util.ErrorRecordNotFound)
		}
		db = db.Where("type = ?", typeID)

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
			return response.FailureForList(util.ErrorSortingFieldDoesNotExist)
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

	//data
	var data []ProgressOutput
	db.Model(&model.Progress{}).Debug().Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	//查拆解信息
	if p.DisassemblyID > 0 {
		var record DisassemblyOutput
		res := global.DB.Model(&model.Disassembly{}).
			Where("id = ?", p.DisassemblyID).Limit(1).Find(&record)
		if res.RowsAffected > 0 {
			for i := range data {
				data[i].DisassemblyExternal = &record
			}
		}
	}

	//一次性查出需要用的进度类型，避免多次重复查询
	var progressType model.DictionaryType
	err := global.DB.Where("name = '进度类型'").
		First(&progressType).Error
	if err != nil {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var planned DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", progressType.ID).
		Where("name = '计划进度'").
		First(&planned).Error
	if err != nil {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var forecasted DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", progressType.ID).
		Where("name = '预测进度'").
		First(&forecasted).Error
	if err != nil {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var actual DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", progressType.ID).
		Where("name = '实际进度'").
		First(&actual).Error
	if err != nil {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var dataSource model.DictionaryType
	err = global.DB.Where("name = '进度的数据来源'").
		First(&dataSource).Error
	if err != nil {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var systemCalculation DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", dataSource.ID).
		Where("name = '系统计算'").
		First(&systemCalculation).Error
	if err != nil {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var manualFilling DictionaryDetailOutput
	err = global.DB.Model(&model.DictionaryDetail{}).
		Where("dictionary_type_id =?", dataSource.ID).
		Where("name = '人工填写'").
		First(&manualFilling).Error
	if err != nil {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	for i := range data {

		//处理日期格式
		if data[i].Date != nil {
			temp := *data[i].Date
			*data[i].Date = temp[:10]
		}

		//查dictionary_item表
		{
			if data[i].Type != nil {
				if *data[i].Type == planned.ID {
					data[i].TypeExternal = &planned
				} else if *data[i].Type == forecasted.ID {
					data[i].TypeExternal = &forecasted
				} else if *data[i].Type == actual.ID {
					data[i].TypeExternal = &actual
				} else {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *data[i].Type).
						Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						data[i].TypeExternal = &record
					}
				}
			}

			if data[i].DataSource != nil {
				if *data[i].DataSource == systemCalculation.ID {
					data[i].DataSourceExternal = &systemCalculation
				} else if *data[i].DataSource == manualFilling.ID {
					data[i].DataSourceExternal = &manualFilling
				} else {
					var record DictionaryDetailOutput
					res := global.DB.Model(&model.DictionaryDetail{}).
						Where("id = ?", *data[i].DataSource).
						Limit(1).Find(&record)
					if res.RowsAffected > 0 {
						data[i].DataSourceExternal = &record
					}
				}
			}
		}
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
		Message: util.GetErrorDescription(util.Success),
	}
}
