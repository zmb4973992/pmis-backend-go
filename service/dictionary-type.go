package service

import (
	"gorm.io/gorm"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type dictionaryTypeService struct{}

func (dictionaryTypeService) Create(paramIn *dto.DictionaryTypeCreateOrUpdate) response.Common {
	var paramOut model.DictionaryType
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}

	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	paramOut.Name = paramIn.Name

	if *paramIn.Sort != -1 {
		paramOut.Sort = paramIn.Sort
	}

	if *paramIn.Remarks != "" {
		paramOut.Remarks = paramIn.Remarks
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (dictionaryTypeService) CreateInBatches(paramIn []dto.DictionaryTypeCreateOrUpdate) response.Common {
	var paramOut []model.DictionaryType
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	for i := range paramIn {
		var record model.DictionaryType
		if paramIn[i].Creator != nil {
			record.Creator = paramIn[i].Creator
		}

		if paramIn[i].LastModifier != nil {
			record.LastModifier = paramIn[i].LastModifier
		}

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
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

// Update 更新为什么要用dto？首先因为很多数据需要绑定，也就是一定要传参；
// 其次是需要清洗
func (dictionaryTypeService) Update(paramIn *dto.DictionaryTypeCreateOrUpdate) response.Common {

	var paramOut model.DictionaryType
	paramOut.ID = paramIn.ID
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	paramOut.Name = paramIn.Name

	if *paramIn.Sort != -1 {
		paramOut.Sort = paramIn.Sort
	}

	if *paramIn.Remarks != "" {
		paramOut.Remarks = paramIn.Remarks
	}

	//清洗完毕，开始update
	err := global.DB.Omit("created_at", "creator").
		Save(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (dictionaryTypeService) Delete(paramIn dto.DictionaryTypeDelete) response.Common {
	//由于删除需要做两件事：软删除+记录删除人，所以需要用事务
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//这里记录删除人，在事务中必须放在前面
		//如果放后面，由于是软删除，系统会找不到这条记录，导致无法更新
		err := tx.Debug().Model(&model.DictionaryType{}).Where("id = ?", paramIn.DictionaryTypeID).
			Update("deleter", *paramIn.Deleter).Error
		if err != nil {
			return err
		}
		//这里删除记录
		err = tx.Delete(&model.DictionaryType{}, paramIn.DictionaryTypeID).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (dictionaryTypeService) List(paramIn dto.DictionaryTypeList) response.List {
	db := global.DB
	//where order limit offset
	if paramIn.NameInclude != "" {
		db = db.Where("name like ?", "%"+paramIn.NameInclude+"%")
	}

	if paramIn.OrderBy != "" {
		db = db.Order(paramIn.OrderBy)
	}

	var count int64
	db.Model(&model.DictionaryType{}).Count(&count)

	var res []string
	db.Model(&model.DictionaryType{}).Select("name").
		Find(&res)

	////生成sql查询条件
	//sqlCondition := util.NewSqlCondition()
	//
	////对paramIn进行清洗
	////这部分是用于where的参数
	//if paramIn.Page > 0 {
	//	sqlCondition.Paging.Page = paramIn.Page
	//}
	////如果参数里的pageSize是整数且大于0、小于等于上限：
	//maxPagingSize := global.Config.PagingConfig.MaxPageSize
	//if paramIn.PageSize > 0 && paramIn.PageSize <= maxPagingSize {
	//	sqlCondition.Paging.PageSize = paramIn.PageSize
	//}
	//
	//if paramIn.ProjectID != nil {
	//	sqlCondition.Equal("project_id", *paramIn.ProjectID)
	//}
	//
	//if paramIn.SuperiorID != nil {
	//	sqlCondition.Equal("superior_id", *paramIn.SuperiorID)
	//}
	//
	//if paramIn.Level != nil {
	//	sqlCondition.Equal("level", *paramIn.Level)
	//}
	//
	//if paramIn.LevelGte != nil {
	//	sqlCondition.Gte("level", *paramIn.LevelGte)
	//}
	//
	//if paramIn.LevelLte != nil {
	//	sqlCondition.Lte("level", *paramIn.LevelLte)
	//}
	//
	////这部分是用于order的参数
	//orderBy := paramIn.OrderBy
	//if orderBy != "" {
	//	ok := sqlCondition.ValidateColumn(orderBy, model.Disassembly{})
	//	if ok {
	//		sqlCondition.Sorting.OrderBy = orderBy
	//	}
	//}
	//desc := paramIn.Desc
	//if desc == true {
	//	sqlCondition.Sorting.Desc = true
	//} else {
	//	sqlCondition.Sorting.Desc = false
	//}
	//
	//tempList := sqlCondition.Find(global.DB, model.Disassembly{})
	//totalRecords := sqlCondition.Count(global.DB, model.Disassembly{})
	//totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)
	//
	//if len(tempList) == 0 {
	//	return response.FailureForList(util.ErrorRecordNotFound)
	//}
	//
	//var list []dto.DisassemblyOutput
	//_ = mapstructure.Decode(&tempList, &list)
	//
	//return response.List{
	//	Data: list,
	//	Paging: &dto.Paging{
	//		Page:         sqlCondition.Paging.Page,
	//		PageSize:     sqlCondition.Paging.PageSize,
	//		TotalPages:   totalPages,
	//		TotalRecords: totalRecords,
	//	},
	//	Code:    util.Success,
	//	Message: util.GetMessage(util.Success),
	//}

	return response.List{
		Data:    res,
		Paging:  nil,
		Code:    int(count),
		Message: "",
	}

}
