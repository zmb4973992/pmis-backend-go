package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

//以下为入参

type ProgressGet struct {
	ID int
}

type ProgressCreate struct {
	Creator      int
	LastModifier int

	DisassemblyID int      `json:"disassembly_id" binding:"required"`
	Date          string   `json:"date" binding:"required"`
	Type          int      `json:"type" binding:"required"`
	Value         *float64 `json:"value" binding:"required"`
	Remark        string   `json:"remark,omitempty"`
}

type ProgressUpdate struct {
	LastModifier int
	ID           int

	DisassemblyID *int     `json:"disassembly_id"`
	Date          *string  `json:"date"`
	Type          *int     `json:"type"`
	Value         *float64 `json:"value"`
	Remark        *string  `json:"remark"`
}

type ProgressDelete struct {
	ID int
}

type ProgressGetList struct {
	dto.ListInput

	DisassemblyID int     `json:"disassembly_id,omitempty"`
	DateGte       string  `json:"date_gte,omitempty"`
	DateLte       string  `json:"date_lte,omitempty"`
	Date          string  `json:"date,omitempty"`
	Type          string  `json:"type,omitempty"`
	ValueGte      float64 `json:"value_gte,omitempty"`
	ValueLte      float64 `json:"value_lte,omitempty"`
	DataSource    int     `json:"data_source,omitempty"`
	DataSourceIn  []int   `json:"data_source_in"`
}

//以下为出参

type ProgressOutput struct {
	Creator      *int `json:"creator"`
	LastModifier *int `json:"last_modifier"`
	ID           int  `json:"id"`

	DisassemblyID *int     `json:"disassembly_id"`
	Date          *string  `json:"date"`
	Type          *int     `json:"type"`
	Value         *float64 `json:"value"`
	Remark        *string  `json:"remark"`
	DataSource    *string  `json:"data_source"`
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

	paramOut.DisassemblyID = &p.DisassemblyID

	date, err := time.Parse("2006-01-02", p.Date)
	if err != nil {
		return response.Failure(util.ErrorInvalidDateFormat)
	}
	paramOut.Date = &date

	paramOut.Type = &p.Type

	paramOut.Value = p.Value

	if p.Remark != "" {
		paramOut.Remark = &p.Remark
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
	var dictionaryItem model.DictionaryItem
	err = global.DB.Where("name = '人工填写'").First(&dictionaryItem).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	res := global.DB.FirstOrCreate(&paramOut, model.Progress{
		DisassemblyID: &p.DisassemblyID,
		Date:          &date,
		Type:          &p.Type,
		Value:         p.Value,
		DataSource:    &dictionaryItem.ID,
	})

	if res.Error != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}

	if res.RowsAffected == 0 {
		return response.Failure(util.ErrorDuplicateRecord)
	}

	//更新所有上级的进度
	err = util.UpdateProgressOfSuperiors(p.DisassemblyID, p.Type)

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

	if p.DisassemblyID != nil {
		if *p.DisassemblyID > 0 {
			paramOut["disassembly_id"] = p.DisassemblyID
		} else {
			paramOut["disassembly_id"] = nil
		}
	}

	if p.Date != nil {
		if *p.Date != "" {
			var err error
			paramOut["date"], err = time.Parse("2006-01-02", *p.Date)
			if err != nil {
				return response.Failure(util.ErrorInvalidJSONParameters)
			}
		} else {
			paramOut["date"] = nil
		}
	}

	if p.Type != nil {
		if *p.Type > 0 {
			paramOut["type"] = p.Type
		} else {
			paramOut["type"] = nil
		}
	}

	if p.Value != nil {
		if *p.Value >= 0 {
			paramOut["value"] = p.Value
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if p.Remark != nil {
		if *p.Remark != "" {
			paramOut["remark"] = *p.Remark
		} else {
			paramOut["remark"] = nil
		}
	}

	//找到“人工填写”在字典详情表的id
	var dictionaryItem model.DictionaryItem
	err := global.DB.Where("name = '人工填写'").First(&dictionaryItem).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	paramOut["data_source"] = dictionaryItem.ID

	//计算有修改值的字段数，分别进行不同处理
	//data_source是自动添加的，也需要排除在外
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "Deleter", "CreateAt", "UpdatedAt", "DeletedAt", "DataSource")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err = global.DB.Model(&model.Progress{}).Where("id = ?", p.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	//更新所有上级的进度
	var progress model.Progress
	err = global.DB.Where("id = ?", p.ID).First(&progress).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
	}

	if progress.DisassemblyID != nil && progress.Type != nil {
		err = util.UpdateProgressOfSuperiors(*progress.DisassemblyID, *progress.Type)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
		}
	}

	return response.Success()
}

func (p *ProgressDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var progress model.Progress
	global.DB.Where("id = ?", p.ID).Find(&progress)
	err := global.DB.Where("id = ?", p.ID).Delete(&progress).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	//更新所有上级的进度
	if progress.DisassemblyID != nil && progress.Type != nil {
		err = util.UpdateProgressOfSuperiors(*progress.DisassemblyID, *progress.Type)
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Failure(util.ErrorFailToCalculateSuperiorProgress)
		}
	}

	return response.Success()
}

func (p *ProgressGetList) GetList() response.List {
	//db := global.DB.Model(&model.Project{})
	//// 顺序：where -> count -> Order -> limit -> offset -> data
	//
	////where
	//if p.NameInclude != "" {
	//	db = db.Where("name like ?", "%"+p.NameInclude+"%")
	//}
	//
	//if p.DepartmentNameInclude != "" {
	//	var departmentIDs []int
	//	global.DB.Model(&model.Department{}).Where("name like ?", "%"+p.DepartmentNameInclude+"%").
	//		Select("id").Find(&departmentIDs)
	//	if len(departmentIDs) > 0 {
	//		db = db.Where("department_id in ?", departmentIDs)
	//	}
	//}
	//
	//if len(p.DepartmentIDIn) > 0 {
	//	db = db.Where("department_id in ?", p.DepartmentIDIn)
	//}
	//
	//if p.IsShowedByRole {
	//	//先获得最大角色的名称
	//	biggestRoleName := util.GetBiggestRoleName(p.UserID)
	//	if biggestRoleName == "事业部级" {
	//		//获取所在事业部的id数组
	//		businessDivisionIDs := util.GetBusinessDivisionIDs(p.UserID)
	//		//获取归属这些事业部的部门id数组
	//		var departmentIDs []int
	//		global.DB.Model(&model.Department{}).Where("superior_id in ?", businessDivisionIDs).
	//			Select("id").Find(&departmentIDs)
	//		//两个数组进行合并
	//		departmentIDs = append(departmentIDs, businessDivisionIDs...)
	//		//找到部门id在上面两个数组中的记录
	//		db = db.Where("department_id in ?", departmentIDs)
	//	} else if biggestRoleName == "部门级" || biggestRoleName == "项目级" {
	//		//获取用户所属部门的id数组
	//		departmentIDs := util.GetDepartmentIDs(p.UserID)
	//		//找到部门id在上面数组中的记录
	//		db = db.Where("department_id in ?", departmentIDs)
	//	}
	//}
	//
	//// count
	//var count int64
	//db.Count(&count)
	//
	////Order
	//orderBy := p.SortingInput.OrderBy
	//desc := p.SortingInput.Desc
	////如果排序字段为空
	//if orderBy == "" {
	//	//如果要求降序排列
	//	if desc == true {
	//		db = db.Order("id desc")
	//	}
	//} else { //如果有排序字段
	//	//先看排序字段是否存在于表中
	//	exists := util.FieldIsInModel(&model.Project{}, orderBy)
	//	if !exists {
	//		return response.FailureForList(util.ErrorSortingFieldDoesNotExist)
	//	}
	//	//如果要求降序排列
	//	if desc == true {
	//		db = db.Order(orderBy + " desc")
	//	} else { //如果没有要求排序方式
	//		db = db.Order(orderBy)
	//	}
	//}
	//
	////limit
	//page := 1
	//if p.PagingInput.Page > 0 {
	//	page = p.PagingInput.Page
	//}
	//pageSize := global.Config.DefaultPageSize
	//if p.PagingInput.PageSize >= 0 &&
	//	p.PagingInput.PageSize <= global.Config.MaxPageSize {
	//	pageSize = p.PagingInput.PageSize
	//}
	//db = db.Limit(pageSize)
	//
	////offset
	//offset := (page - 1) * pageSize
	//db = db.Offset(offset)
	//
	////data
	//var data []ProjectOutput
	//db.Model(&model.Project{}).Find(&data)
	//
	//if len(data) == 0 {
	//	return response.FailureForList(util.ErrorRecordNotFound)
	//}
	//
	//for i := range data {
	//	//查部门信息
	//	if data[i].DepartmentID != nil {
	//		var record DepartmentOutput
	//		res := global.DB.Model(&model.Department{}).
	//			Where("id=?", *data[i].DepartmentID).Limit(1).Find(&record)
	//		if res.RowsAffected > 0 {
	//			data[i].DepartmentExternal = &record
	//		}
	//	}
	//
	//	//处理日期格式
	//	if data[i].SigningDate != nil {
	//		temp := *data[i].SigningDate
	//		*data[i].SigningDate = temp[:10]
	//	}
	//
	//	if data[i].EffectiveDate != nil {
	//		temp := *data[i].EffectiveDate
	//		*data[i].EffectiveDate = temp[:10]
	//	}
	//
	//	//查dictionary_item表
	//	{
	//		if data[i].Country != nil {
	//			var record DictionaryItemOutput
	//			res := global.DB.Model(&model.DictionaryItem{}).
	//				Where("id = ?", *data[i].Country).Limit(1).Find(&record)
	//			if res.RowsAffected > 0 {
	//				data[i].CountryExternal = &record
	//			}
	//		}
	//
	//		if data[i].Type != nil {
	//			var record DictionaryItemOutput
	//			res := global.DB.Model(&model.DictionaryItem{}).
	//				Where("id = ?", *data[i].Type).Limit(1).Find(&record)
	//			if res.RowsAffected > 0 {
	//				data[i].TypeExternal = &record
	//			}
	//		}
	//
	//		if data[i].Currency != nil {
	//			var record DictionaryItemOutput
	//			res := global.DB.Model(&model.DictionaryItem{}).
	//				Where("id = ?", *data[i].Currency).Limit(1).Find(&record)
	//			if res.RowsAffected > 0 {
	//				data[i].CurrencyExternal = &record
	//			}
	//		}
	//
	//		if data[i].Status != nil {
	//			var record DictionaryItemOutput
	//			res := global.DB.Model(&model.DictionaryItem{}).
	//				Where("id = ?", *data[i].Status).Limit(1).Find(&record)
	//			if res.RowsAffected > 0 {
	//				data[i].StatusExternal = &record
	//			}
	//		}
	//
	//		if data[i].OurSignatory != nil {
	//			var record DictionaryItemOutput
	//			res := global.DB.Model(&model.DictionaryItem{}).
	//				Where("id = ?", *data[i].OurSignatory).Limit(1).Find(&record)
	//			if res.RowsAffected > 0 {
	//				data[i].OurSignatoryExternal = &record
	//			}
	//		}
	//	}
	//}
	//
	//numberOfRecords := int(count)
	//numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)
	//
	//return response.List{
	//	Data: data,
	//	Paging: &dto.PagingOutput{
	//		Page:            page,
	//		PageSize:        pageSize,
	//		NumberOfPages:   numberOfPages,
	//		NumberOfRecords: numberOfRecords,
	//	},
	//	Code:    util.Success,
	//	Message: util.GetMessage(util.Success),
	//}

	return response.List{}
}
