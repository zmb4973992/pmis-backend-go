package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
)

type projectBreakdownDAO struct{}

func (projectBreakdownDAO) Get(projectBreakdownID int) *dto.ProjectBreakdownGetDTO {
	var param dto.ProjectBreakdownGetDTO
	//把基础的拆解信息查出来
	var projectBreakdown model.ProjectBreakdown
	err := global.DB.Where("id = ?", projectBreakdownID).First(&projectBreakdown).Error
	if err != nil {
		return nil
	}
	//把所有查出的结果赋值给输出变量
	if projectBreakdown.Name != nil {
		param.Name = projectBreakdown.Name
	}
	if projectBreakdown.ProjectID != nil {
		param.ProjectID = projectBreakdown.ProjectID
	}
	if projectBreakdown.Level != nil {
		param.Level = projectBreakdown.Level
	}
	if projectBreakdown.Weight != nil {
		param.Weight = projectBreakdown.Weight
	}
	if projectBreakdown.SuperiorID != nil {
		param.SuperiorID = projectBreakdown.SuperiorID
	}
	return &param
}

// Create 这里是只负责新增，不写任何业务逻辑。只要收到参数就创建数据库记录，然后返回错误
func (projectBreakdownDAO) Create(param *model.ProjectBreakdown) error {
	err := global.DB.Create(param).Error
	return err
}

// Update 这里是只负责更新，不写任何业务逻辑。只要收到id和更新参数，然后返回错误
func (projectBreakdownDAO) Update(param *model.ProjectBreakdown) error {
	//注意，这里就算没有找到记录，也不会报错，只有更新字段出现问题才会报错。详见gorm的update用法
	err := global.DB.Where("id = ?", param.ID).Omit("created_at").Save(param).Error
	return err
}

func (projectBreakdownDAO) Delete(projectBreakdownID int) error {
	//注意，这里就算没有找到记录，也不会报错。详见gorm的delete用法
	err := global.DB.Delete(&model.ProjectBreakdown{}, projectBreakdownID).Error
	return err
}

// List 这里是只负责查询列表，不写任何业务逻辑。
// 查询数据库记录列表，返回dto
// 入参为sql查询条件，结果为数据列表+分页情况
// 已通过sqlCondition实现，废弃
//func (ProjectBreakdownDAO) List(sqlCondition util.SqlCondition) (
//	list []dto.ProjectBreakdownGetDTO, totalPages int, totalRecords int) {
//
//	db := model.DB
//	//select columns
//	if len(sqlCondition.SelectedColumns) > 0 {
//		db = db.Select(sqlCondition.SelectedColumns)
//	}
//	//where
//	for _, paramPair := range sqlCondition.ParamPairs {
//		db = db.Where(paramPair.Key, paramPair.Value)
//	}
//	//orderBy
//	orderBy := sqlCondition.Sorting.OrderBy
//	if orderBy != "" {
//		if sqlCondition.Sorting.Desc == true {
//			db = db.Order(sqlCondition.Sorting.OrderBy + " desc")
//		} else {
//			db = db.Order(sqlCondition.Sorting.OrderBy)
//		}
//	}
//	//count 计算totalRecords
//	var tempTotalRecords int64
//	err := db.Model(&model.ProjectBreakdown{}).Count(&tempTotalRecords).Error
//
//	if err != nil {
//		return nil, 0, 0
//	}
//	totalRecords = int(tempTotalRecords)
//
//	//limit
//	db = db.Limit(sqlCondition.Paging.PageSize)
//	//offset
//	offset := (sqlCondition.Paging.Page - 1) * sqlCondition.Paging.PageSize
//	db = db.Offset(offset)
//
//	//count 计算totalPages
//	totalPages = util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)
//	err = db.Debug().Model(&model.ProjectBreakdown{}).Find(&list).Error
//
//	if err != nil {
//		return nil, 0, 0
//	}
//	return list, totalPages, totalRecords
//}
