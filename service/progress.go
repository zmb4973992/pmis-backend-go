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
	SnowID int64
}

type ProgressCreate struct {
	Creator      int64
	LastModifier int64

	DisassemblySnowID int64    `json:"disassembly_snow_id" binding:"required"`
	Date              string   `json:"date" binding:"required"`
	Type              int64    `json:"type" binding:"required"`
	Value             *float64 `json:"value" binding:"required"`
	Remarks           string   `json:"remarks,omitempty"`
}

type ProgressUpdate struct {
	LastModifier int64
	SnowID       int64

	Date    *string  `json:"date"`
	Type    *int64   `json:"type"`
	Value   *float64 `json:"value"`
	Remarks *string  `json:"remarks"`
}

type ProgressDelete struct {
	SnowID int64
}

type ProgressGetList struct {
	list.Input

	DisassemblySnowID int64    `json:"disassembly_snow_id" binding:"required"`
	DateGte           string   `json:"date_gte,omitempty"`
	DateLte           string   `json:"date_lte,omitempty"`
	Type              int64    `json:"type,omitempty"`
	TypeIn            []int64  `json:"type_in"`
	ValueGte          *float64 `json:"value_gte"`
	ValueLte          *float64 `json:"value_lte"`
	DataSource        int64    `json:"data_source"`
}

//以下为出参

type ProgressOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	SnowID       int64  `json:"snow_id"`

	DisassemblySnowID   *int64             `json:"-"`
	DisassemblyExternal *DisassemblyOutput `json:"disassembly" gorm:"-"`

	Date               *string                 `json:"date"`
	Type               *int64                  `json:"-"`
	TypeExternal       *DictionaryDetailOutput `json:"type" gorm:"-"`
	Value              *float64                `json:"value"`
	Remarks            *string                 `json:"remarks"`
	DataSource         *string                 `json:"-"`
	DataSourceExternal *DictionaryDetailOutput `json:"data_source" gorm:"-"`
}

func (p *ProgressGet) Get() response.Common {
	var result ProgressOutput
	err := global.DB.Model(model.Progress{}).
		Where("snow_id = ?", p.SnowID).First(&result).Error
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
				Where("snow_id = ?", *result.Type).
				Limit(1).Find(&record)
			if res.RowsAffected > 0 {
				result.TypeExternal = &record
			}
		}

		if result.DataSource != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("snow_id = ?", *result.DataSource).
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

	if p.Creator > 0 {
		paramOut.Creator = &p.Creator
	}

	if p.LastModifier > 0 {
		paramOut.LastModifier = &p.LastModifier
	}

	paramOut.DisassemblySnowID = &p.DisassemblySnowID

	date, err := time.Parse("2006-01-02", p.Date)
	if err != nil {
		return response.Failure(util.ErrorInvalidDateFormat)
	}
	paramOut.Date = &date
	paramOut.Type = &p.Type
	paramOut.Value = p.Value

	//找到"人工填写"的dictionary_item值
	var dataSource model.DictionaryDetail
	err = global.DB.Where("name = '人工填写'").First(&dataSource).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOut.DataSource = &dataSource.SnowID

	if p.Remarks != "" {
		paramOut.Remarks = &p.Remarks
	}

	//计算有修改值的字段数，分别进行不同处理
	tempParamOut, err := util.StructToMap(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	paramOutForCounting := util.MapCopy(tempParamOut,
		"Creator", "LastModifier", "Deleter", "CreateAt", "UpdatedAt", "DeletedAt")

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
		DisassemblySnowID: &p.DisassemblySnowID,
		Date:              &date,
		Type:              &p.Type,
	})

	if res.Error != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	if res.RowsAffected == 0 {
		return response.Failure(util.ErrorDuplicateRecord)
	}

	//更新所有上级的进度
	err = util.UpdateProgressOfSuperiors(p.DisassemblySnowID, p.Type)

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
	}

	return response.Success()
}

func (p *ProgressUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if p.LastModifier > 0 {
		paramOut["last_modifier"] = p.LastModifier
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
	err := global.DB.Where("name = '人工填写'").First(&dataSource).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	paramOut["data_source"] = dataSource.SnowID

	//计算有修改值的字段数，分别进行不同处理
	//data_source是自动添加的，也需要排除在外
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt", "DataSource")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	//找到待更新的这条记录
	var progress model.Progress
	err = global.DB.Where("snow_id = ?", p.SnowID).First(&progress).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
	}

	//如果修改了date或type，意味着可能有重复记录，需要进行判断
	if p.Date != nil || p.Type != nil {
		//从数据库找出相同拆解id、相同日期、相同类型的记录
		var tempProgressSnowIDs []int64
		tempDate, err1 := time.Parse("2006-01-02", *p.Date)
		if err1 != nil {
			return response.Failure(util.ErrorInvalidDateFormat)
		}
		global.DB.Model(&model.Progress{}).Where(&model.Progress{
			DisassemblySnowID: progress.DisassemblySnowID,
			Date:              &tempDate,
			Type:              p.Type,
		}).Select("snow_id").Find(&tempProgressSnowIDs)
		//如果数据库有记录、且待修改的progressID不在数据库记录的progressIDs里面，说明是新的记录
		//则不允许修改
		if len(tempProgressSnowIDs) > 0 && !util.IsInSlice(p.SnowID, tempProgressSnowIDs) {
			return response.Failure(util.ErrorDuplicateRecord)
		}
	}

	//更新记录
	err = global.DB.Model(&model.Progress{}).Where("snow_id = ?", p.SnowID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	//更新所有上级的进度
	if progress.DisassemblySnowID != nil {
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
			var progressTypeSnowIDs []int64
			global.DB.Model(&model.DictionaryDetail{}).
				Where("dictionary_type_snow_id = ?", progressTypeIDInDictionaryType.SnowID).
				Select("snow_id").Find(&progressTypeSnowIDs)

			for _, v := range progressTypeSnowIDs {
				err = util.UpdateProgressOfSuperiors(*progress.DisassemblySnowID, v)
				if err != nil {
					global.SugaredLogger.Errorln(err)
					return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
				}
			}
		} else { //如果没有传入type值(意味着记录的type值不变)，则只更新原来的进度类型
			err = util.UpdateProgressOfSuperiors(*progress.DisassemblySnowID, *progress.Type)
			if err != nil {
				global.SugaredLogger.Errorln(err)
				return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
			}
		}
	}

	return response.Success()
}

func (p *ProgressDelete) Delete() response.Common {
	//先找到记录，这样参数才能获得值、触发钩子函数，再删除记录
	var progress model.Progress
	err := global.DB.Where("snow_id = ?", p.SnowID).Find(&progress).Delete(&progress).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	//更新所有上级的进度
	if progress.DisassemblySnowID != nil && progress.Type != nil {
		err = util.UpdateProgressOfSuperiors(*progress.DisassemblySnowID, *progress.Type)
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
	db = db.Where("disassembly_snow_id = ?", p.DisassemblySnowID)

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
			db = db.Order("snow_id desc")
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
	db.Model(&model.Progress{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	//查拆解信息
	if p.DisassemblySnowID > 0 {
		var record DisassemblyOutput
		res := global.DB.Model(&model.Disassembly{}).
			Where("snow_id = ?", p.DisassemblySnowID).Limit(1).Find(&record)
		if res.RowsAffected > 0 {
			for i := range data {
				data[i].DisassemblyExternal = &record
			}
		}
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
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("snow_id = ?", *data[i].Type).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].TypeExternal = &record
				}
			}

			if data[i].DataSource != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("snow_id = ?", *data[i].DataSource).Limit(1).Find(&record)
				if res.RowsAffected > 0 {
					data[i].DataSourceExternal = &record
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
		Message: util.GetMessage(util.Success),
	}
}
