package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

// User 没有数据、只有方法，所有的数据都放在DTO里
// 这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
// 所有的增删改查都交给DAO层处理，否则service层会非常庞大
type roleAndUser struct{}

func (*roleAndUser) ListByRoleID(roleID int) response.Common {
	var userIDs []int
	err := global.DB.Model(&model.RoleAndUser{}).Where("role_id = ?", roleID).Select("user_id").Find(&userIDs).Error
	if err != nil || len(userIDs) == 0 {
		return response.Fail(util.ErrorRecordNotFound)
	}

	//构建返回结果
	data := make(map[string]any)
	data["role_id"] = roleID
	data["user_ids"] = userIDs

	return response.SucceedWithData(data)
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
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (*roleAndUser) UpdateByRoleID(roleID int, paramIn dto.RoleAndUserCreateOrUpdate) response.Common {
	//先删掉原始记录
	err := global.DB.Where("role_id = ?", roleID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	//如果入参是空的切片
	if len(paramIn.UserIDs) == 0 {
		return response.Succeed()
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
		return response.Fail(util.ErrorFailToUpdateRecord)
	}
	return response.Succeed()
}

func (*roleAndUser) DeleteByRoleID(roleID int) response.Common {
	err := global.DB.Where("role_id = ?", roleID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}

	return response.Succeed()
}

func (*roleAndUser) ListByUserID(userID int) response.Common {
	var roleIDs []int
	err := global.DB.Model(&model.RoleAndUser{}).Where("user_id = ?", userID).Select("role_id").Find(&roleIDs).Error
	if err != nil || len(roleIDs) == 0 {
		return response.Fail(util.ErrorRecordNotFound)
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

	return response.SucceedWithData(data)
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
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (*roleAndUser) UpdateByUserID(userID int, paramIn dto.RoleAndUserCreateOrUpdate) response.Common {
	//先删掉原始记录
	err := global.DB.Where("user_id = ?", userID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	//如果入参是空的切片
	if len(paramIn.RoleIDs) == 0 {
		return response.Succeed()
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
		return response.Fail(util.ErrorFailToUpdateRecord)
	}
	return response.Succeed()
}

func (*roleAndUser) DeleteByUserID(userID int) response.Common {
	err := global.DB.Where("user_id = ?", userID).Delete(&model.RoleAndUser{}).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}

	return response.Succeed()
}
