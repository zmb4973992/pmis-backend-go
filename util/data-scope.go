package util

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

func GetOrganizationIDs(userID int64) (organizationIDs []int64) {
	//先获取角色
	var userAndDataScope model.UserAndDataScope
	err := global.DB.
		Where("user_id = ?", userID).
		First(&userAndDataScope).Error
	if err != nil {
		return nil
	}

	//获得数据范围的信息
	var dataScope model.DataScope
	err = global.DB.
		Where("id = ?", userAndDataScope.DataScopeID).
		First(&dataScope).Error
	if err != nil {
		return nil
	}
	fmt.Println("数据范围名称：", dataScope.Name)

	//如果数据范围是"所有部门"，就返回全部的组织id
	if dataScope.Name == "所有部门" {
		global.DB.Model(&model.Organization{}).
			Select("id").Find(&organizationIDs)
		return
	} else if dataScope.Name == "所属部门和子部门" {
		var tempOrganizationIDs []int64
		global.DB.Model(&model.OrganizationAndUser{}).
			Where("user_id = ?", userID).
			Select("organization_id").
			Find(&tempOrganizationIDs)
		for i := range tempOrganizationIDs {
			res := getSubOrganizationIDs(tempOrganizationIDs[i])
			tempOrganizationIDs = append(tempOrganizationIDs, res...)
		}
		organizationIDs = RemoveDuplication(tempOrganizationIDs)
		return
	} else if dataScope.Name == "所属部门" {
		global.DB.Model(&model.OrganizationAndUser{}).
			Where("user_id = ?", userID).
			Select("organization_id").Find(&organizationIDs)
		return
	} else if dataScope.Name == "无权限" {
		return
	}

	return
}

func getSubOrganizationIDs(organizationID int64) (organizationIDs []int64) {
	global.DB.Model(&model.Organization{}).
		Where("superior_id = ?", organizationID).
		Select("id").
		Find(&organizationIDs)
	for i := range organizationIDs {
		res := getSubOrganizationIDs(organizationIDs[i])
		organizationIDs = append(organizationIDs, res...)
	}
	return
}
