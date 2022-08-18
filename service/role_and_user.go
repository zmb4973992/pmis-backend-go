package service

import (
	"github.com/mitchellh/mapstructure"
	"learn-go/dao"
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
	"learn-go/serializer/response"
	"learn-go/util"
)

// UserService 没有数据、只有方法，所有的数据都放在DTO里
//这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
//所有的增删改查都交给DAO层处理，否则service层会非常庞大
type roleAndUserService struct{}

func (roleAndUserService) Create(paramIn *dto.RoleAndUserCreateDTO) response.Common {
	//对数据进行清洗
	var paramOut model.RoleAndUser
	paramOut.RoleID = paramIn.RoleID
	paramOut.UserID = paramIn.UserID

	err := dao.RoleAndUserDAO.Create(&paramOut)

	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (roleAndUserService) CreateInBatch(paramIn []dto.RoleAndUserCreateDTO) response.Common {
	var paramOut []model.RoleAndUser
	for i := range paramIn {
		var record model.RoleAndUser
		record.RoleID = paramIn[i].RoleID
		record.UserID = paramIn[i].UserID
		paramOut = append(paramOut, record)
	}

	err := dao.RoleAndUserDAO.CreateInBatch(paramOut)
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (roleAndUserService) UpdateUserIDByRoleID(roleID int, paramIn dto.RoleAndUserCreateOrUpdateDTO) response.Common {
	//先删掉原始记录
	err := global.DB.Where("role_id = ?", roleID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	//如果入参是空的切片
	if len(paramIn.UserIDs) == 0 {
		return response.Success()
	}

	//再增加新的记录
	var paramOut []model.RoleAndUser
	for i := range paramIn.UserIDs {
		var record model.RoleAndUser
		record.RoleID = &roleID
		record.UserID = &paramIn.UserIDs[i]
		paramOut = append(paramOut, record)
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (roleAndUserService) UpdateRoleIDByUserID(userID int, paramIn dto.RoleAndUserCreateOrUpdateDTO) response.Common {
	//先删掉原始记录
	err := global.DB.Where("user_id = ?", userID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	//如果入参是空的切片
	if len(paramIn.RoleIDs) == 0 {
		return response.Success()
	}

	//再增加新的记录
	var paramOut []model.RoleAndUser
	for i := range paramIn.RoleIDs {
		var record model.RoleAndUser
		record.UserID = &userID
		record.RoleID = &paramIn.RoleIDs[i]
		paramOut = append(paramOut, record)
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (roleAndUserService) Delete(paramIn dto.RoleAndUserDeleteDTO) response.Common {
	var paramPairs []util.ParamPair

	if paramIn.RoleID != nil {
		paramPairs = append(paramPairs, util.ParamPair{
			Key:   "role_id = ?",
			Value: *paramIn.RoleID,
		})
	}

	if paramIn.UserID != nil {
		paramPairs = append(paramPairs, util.ParamPair{
			Key:   "user_id = ?",
			Value: *paramIn.UserID,
		})
	}

	if len(paramPairs) == 0 {
		return response.Success()
	}

	err := dao.RoleAndUserDAO.Delete(paramPairs)
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	return response.Success()
}

func (roleAndUserService) List(paramIn dto.RoleAndUserListDTO) response.List {
	//生成sql查询条件
	sqlCondition := util.NewSqlCondition()

	//对paramIn进行清洗
	//分页
	if paramIn.Page > 0 {
		sqlCondition.Paging.Page = paramIn.Page
	}
	//如果参数里的pageSize是整数且大于0、小于等于上限：
	maxPagingSize := global.Config.PagingConfig.MaxPageSize
	if paramIn.PageSize > 0 && paramIn.PageSize <= maxPagingSize {
		sqlCondition.Paging.PageSize = paramIn.PageSize
	}

	//这部分是用于where的参数
	if paramIn.RoleID != nil {
		sqlCondition.Equal("role_id", *paramIn.RoleID)
	}

	if paramIn.UserID != nil {
		sqlCondition.Equal("user_id", *paramIn.UserID)
	}

	//这部分是用于order的参数
	orderBy := paramIn.OrderBy
	if orderBy != "" {
		ok := sqlCondition.ValidateColumn(orderBy, model.RoleAndUser{})
		if ok {
			sqlCondition.Sorting.OrderBy = orderBy
		}
	}
	desc := paramIn.Desc
	if desc == true {
		sqlCondition.Sorting.Desc = true
	} else {
		sqlCondition.Sorting.Desc = false
	}

	tempList := sqlCondition.Find(model.RoleAndUser{})
	totalRecords := sqlCondition.Count(model.RoleAndUser{})
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	//这里的tempList是基于model的，不能直接传给前端，要处理成dto才行
	//如果map的字段类型和struct的字段类型不匹配，数据不会同步过来
	var list []dto.RoleAndUserGetDTO
	_ = mapstructure.Decode(&tempList, &list)

	return response.List{
		Data: list,
		Paging: &dto.PagingDTO{
			Page:         sqlCondition.Paging.Page,
			PageSize:     sqlCondition.Paging.PageSize,
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
