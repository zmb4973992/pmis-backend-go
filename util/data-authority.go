package util

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

func GetOrganizationIdsForDataAuthority(userId int64) (organizationIds []int64) {
	//先获取角色
	var userAndDataAuthority model.UserAndDataAuthority
	err := global.DB.
		Where("user_id = ?", userId).
		First(&userAndDataAuthority).Error
	if err != nil {
		return nil
	}

	//获得数据范围的信息
	var dataAuthority model.DataAuthority
	err = global.DB.
		Where("id = ?", userAndDataAuthority.DataAuthorityId).
		First(&dataAuthority).Error
	if err != nil {
		return nil
	}
	//fmt.Println("数据范围名称：", dataAuthority.Name)

	//如果数据范围是"所有部门"，就返回全部的组织id
	if dataAuthority.Name == "所有部门" {
		global.DB.Model(&model.Organization{}).
			Select("id").
			Find(&organizationIds)
		return
	} else if dataAuthority.Name == "所属部门和子部门" {
		var tempOrganizationIds []int64
		global.DB.Model(&model.OrganizationAndUser{}).
			Where("user_id = ?", userId).
			Select("organization_id").
			Find(&tempOrganizationIds)
		for i := range tempOrganizationIds {
			res := getSubOrganizationIds(tempOrganizationIds[i])
			tempOrganizationIds = append(tempOrganizationIds, res...)
		}
		organizationIds = RemoveDuplication(tempOrganizationIds)
		return
	} else if dataAuthority.Name == "所属部门" {
		global.DB.Model(&model.OrganizationAndUser{}).
			Where("user_id = ?", userId).
			Select("organization_id").
			Find(&organizationIds)
		return
	} else if dataAuthority.Name == "无权限" {
		return
	}

	return
}

func getSubOrganizationIds(organizationId int64) (organizationIds []int64) {
	global.DB.Model(&model.Organization{}).
		Where("superior_id = ?", organizationId).
		Select("id").
		Find(&organizationIds)
	for i := range organizationIds {
		res := getSubOrganizationIds(organizationIds[i])
		organizationIds = append(organizationIds, res...)
	}
	return
}
