package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
	"learn-go/util"
)

/*
dao层的基本原则：
入参为id或model，用于对数据库进行增删改查；
出参为err或dto，用于反馈结果或给其他层使用
*/

// RelatedPartyDAO dao层的结构体没有数据，只是操作数据库进行增删改查，不写业务逻辑
type relatedPartyDAO struct{}

// Get 这里是只负责查询，不写任何业务逻辑。
// 查询数据库记录，返回dto
func (relatedPartyDAO) Get(id int) *dto.RelatedPartyDTO {
	//之所以用dto不用model，是因为model为数据库原表，数据可能包含敏感字段、或未加工，不适合直接传递
	//展现的功能基本都交给dto
	var param dto.RelatedPartyDTO
	err := global.DB.Debug().Model(&model.RelatedParty{}).Where("id = ?", id).First(&param).Error
	if err != nil {
		return nil
	}
	return &param
}

// Create 这里是只负责新增，不写任何业务逻辑。
// 创建数据库记录，返回错误
func (relatedPartyDAO) Create(paramIn *model.RelatedParty) error {
	err := global.DB.Create(paramIn).Error
	return err
}

// Update 这里是只负责修改，不写任何业务逻辑。
// 全量更新，保存数据库记录，返回错误
func (relatedPartyDAO) Update(paramIn *model.RelatedParty) error {
	err := global.DB.Debug().Model(paramIn).Omit("created_at").Save(paramIn).Error
	return err
}

// Delete 这里是只负责删除，不写任何业务逻辑。
// 删除数据库记录，返回错误
func (relatedPartyDAO) Delete(id int) error {
	err := global.DB.Debug().Delete(&model.RelatedParty{}, id).Error
	return err
}

// List 这里是只负责查询列表，不写任何业务逻辑。
// 查询数据库记录列表，返回dto
// 入参为sql查询条件，结果为数据列表+分页情况
// List 入参为sql查询条件，结果为数据列表+分页情况
func (relatedPartyDAO) List(sqlCondition util.SqlCondition) (
	list []dto.RelatedPartyDTO, totalPages int, totalRecords int) {
	db := global.DB
	//select
	if len(sqlCondition.SelectedColumns) > 0 {
		db = db.Select(sqlCondition.SelectedColumns)
	}
	//where
	for _, paramPair := range sqlCondition.ParamPairs {
		db = db.Where(paramPair.Key, paramPair.Value)
	}
	//orderBy
	orderBy := sqlCondition.Sorting.OrderBy
	if orderBy != "" {
		var columns []string
		db.Raw("Select Name FROM SysColumns where id=Object_Id('related_party')").
			Find(&columns)
		if ok := util.IsInSlice(orderBy, columns); ok {
			if sqlCondition.Sorting.Desc == true {
				db = db.Order(sqlCondition.Sorting.OrderBy + " desc")
			} else {
				db = db.Order(sqlCondition.Sorting.OrderBy)
			}
		}
	}
	//count 计算totalRecords
	var tempTotalRecords int64
	err := db.Debug().Model(&model.RelatedParty{}).Count(&tempTotalRecords).Error
	if err != nil {
		return nil, 0, 0
	}
	totalRecords = int(tempTotalRecords)

	//limit
	db = db.Limit(sqlCondition.Paging.PageSize)
	//offset
	offset := (sqlCondition.Paging.Page - 1) * sqlCondition.Paging.PageSize
	db = db.Offset(offset)

	//count 计算totalPages
	totalPages = util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)
	err = db.Debug().Model(&model.RelatedParty{}).Find(&list).Error

	if err != nil {
		return nil, 0, 0
	}
	return list, totalPages, totalRecords
}
