package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

//2022/2/3
//这个逻辑比较乱，最后改

type roleAndUser struct{}

func (*roleAndUser) ListByRoleID(roleID int) response.Common {
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

func (*roleAndUser) CreateByRoleID(roleID int, paramIn dto.RoleAndUserCreateOrUpdate) response.Common {
	var paramOut []model.RoleAndUser
	for i := range paramIn.UserIDs {
		var record model.RoleAndUser
		if paramIn.Creator != nil {
			record.Creator = paramIn.Creator
		}

		if paramIn.LastModifier != nil {
			record.LastModifier = paramIn.LastModifier
		}

		record.RoleID = &roleID
		record.UserID = &paramIn.UserIDs[i]
		paramOut = append(paramOut, record)
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (*roleAndUser) UpdateByRoleID(roleID int, paramIn dto.RoleAndUserCreateOrUpdate) response.Common {
	//先删掉原始记录
	err := global.DB.Where("role_id = ?", roleID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
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
		if paramIn.LastModifier != nil {
			record.LastModifier = paramIn.LastModifier
		}

		record.RoleID = &roleID
		record.UserID = &paramIn.UserIDs[i]
		paramOut = append(paramOut, record)
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (*roleAndUser) DeleteByRoleID(roleID int) response.Common {
	err := global.DB.Where("role_id = ?", roleID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	return response.Success()
}

func (*roleAndUser) ListByUserID(userID int) response.Common {
	var roleIDs []int
	err := global.DB.Model(&model.RoleAndUser{}).Where("user_id = ?", userID).Select("role_id").Find(&roleIDs).Error
	if err != nil || len(roleIDs) == 0 {
		return response.Failure(util.ErrorRecordNotFound)
	}

	var roleNames []string
	for _, roleID := range roleIDs {
		var roleName string
		global.DB.Model(&model.Role{}).Where("id = ?", roleID).Select("name").Find(&roleName)
		roleNames = append(roleNames, roleName)
	}

	//构建返回结果
	data := make(map[string]any)
	data["user_id"] = userID
	//data["role_ids"] = roleIDs
	data["role_names"] = roleNames

	return response.SuccessWithData(data)
}

func (*roleAndUser) CreateByUserID(userID int, paramIn dto.RoleAndUserCreateOrUpdate) response.Common {
	var paramOut []model.RoleAndUser
	for i := range paramIn.RoleIDs {
		var record model.RoleAndUser
		if paramIn.Creator != nil {
			record.Creator = paramIn.Creator
		}

		if paramIn.LastModifier != nil {
			record.LastModifier = paramIn.LastModifier
		}

		record.UserID = &userID
		record.RoleID = &paramIn.RoleIDs[i]
		paramOut = append(paramOut, record)
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (*roleAndUser) UpdateByUserID(userID int, paramIn dto.RoleAndUserCreateOrUpdate) response.Common {
	//先删掉原始记录
	err := global.DB.Where("user_id = ?", userID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
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
		if paramIn.LastModifier != nil {
			record.LastModifier = paramIn.LastModifier
		}
		record.UserID = &userID
		record.RoleID = &paramIn.RoleIDs[i]
		paramOut = append(paramOut, record)
	}

	err = global.DB.Create(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (*roleAndUser) DeleteByUserID(userID int) response.Common {
	err := global.DB.Where("user_id = ?", userID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}

	return response.Success()
}
