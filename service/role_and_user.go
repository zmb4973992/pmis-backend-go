package service

import (
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

func (roleAndUserService) ListByRoleID(roleID int) response.Common {
	var userIDs []int
	err := global.DB.Model(&model.RoleAndUser{}).Where("role_id = ?", roleID).Select("user_id").Find(&userIDs).Error
	if err != nil || len(userIDs) == 0 {
		return response.Failure(util.ErrorRecordNotFound)
	}

	//构建返回结果
	data := make(map[string]any)
	data["role_id"] = roleID
	data["user_ids"] = userIDs

	return response.SuccessWithData(data)
}

func (roleAndUserService) CreateByRoleID(roleID int, paramIn dto.RoleAndUserCreateOrUpdateDTO) response.Common {
	var paramOut []model.RoleAndUser
	for i := range paramIn.UserIDs {
		var record model.RoleAndUser
		record.RoleID = &roleID
		record.UserID = &paramIn.UserIDs[i]
		paramOut = append(paramOut, record)
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (roleAndUserService) UpdateByRoleID(roleID int, paramIn dto.RoleAndUserCreateOrUpdateDTO) response.Common {
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

func (roleAndUserService) DeleteByRoleID(roleID int) response.Common {
	err := global.DB.Where("role_id = ?", roleID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	return response.Success()
}

func (roleAndUserService) ListByUserID(userID int) response.Common {
	var roleIDs []int
	err := global.DB.Model(&model.RoleAndUser{}).Where("user_id = ?", userID).Select("role_id").Find(&roleIDs).Error
	if err != nil || len(roleIDs) == 0 {
		return response.Failure(util.ErrorRecordNotFound)
	}

	//构建返回结果
	data := make(map[string]any)
	data["user_id"] = userID
	data["role_ids"] = roleIDs

	return response.SuccessWithData(data)
}

func (roleAndUserService) CreateByUserID(userID int, paramIn dto.RoleAndUserCreateOrUpdateDTO) response.Common {
	var paramOut []model.RoleAndUser
	for i := range paramIn.RoleIDs {
		var record model.RoleAndUser
		record.UserID = &userID
		record.RoleID = &paramIn.RoleIDs[i]
		paramOut = append(paramOut, record)
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (roleAndUserService) UpdateByUserID(userID int, paramIn dto.RoleAndUserCreateOrUpdateDTO) response.Common {
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

func (roleAndUserService) DeleteByUserID(userID int) response.Common {
	err := global.DB.Where("user_id = ?", userID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	return response.Success()
}
