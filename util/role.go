package util

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

func GetOrganizationIDsForDataScope(userID int) (organizationIDsForDataScope []int) {
	//先获取角色
	var roleIDs []int
	global.DB.Model(&model.UserAndRole{}).Where("user_id = ?", userID).
		Select("role_id").Find(&roleIDs)

	//获得所属角色的数据范围类型
	var dataScopeType []int
	global.DB.Model(&model.Role{}).Where("id in ?", roleIDs).
		Select("data_scope_type").Find(&dataScopeType)

	var tempOrganizationIDs []int
	//判断数据范围的类型
	switch {
	//如果数据范围的类型为AllOrganization，就返回全部的组织id
	case SliceIncludes(dataScopeType, model.AllOrganization):
		global.DB.Model(&model.Organization{}).Select("id").Find(&organizationIDsForDataScope)
		return
		//如果数据范围的类型为HisOrganizationAndInferiors，就返回该条件下的所有组织id
		//并继续向下穿透执行
	case SliceIncludes(dataScopeType, model.HisOrganizationAndInferiors):
		tempOrganizationIDs = GetOrganizationIDsWithInferiors(userID)
		fallthrough
	case SliceIncludes(dataScopeType, model.HisOrganization):
		var tempOrganizationIDs1 []int
		global.DB.Model(&model.OrganizationAndUser{}).Where("user_id = ?", userID).
			Select("organization_id").Find(&tempOrganizationIDs1)
		tempOrganizationIDs = append(tempOrganizationIDs, tempOrganizationIDs1...)
		fallthrough
	case SliceIncludes(dataScopeType, model.CustomOrganization):
		var tempOrganizationIDs2 []int
		global.DB.Model(&model.RoleAndOrganizationForDataScope{}).Where("role_id in ?", roleIDs).
			Select("organization_id_for_data_scope").Find(&tempOrganizationIDs2)
		tempOrganizationIDs = append(tempOrganizationIDs, tempOrganizationIDs2...)
	}
	organizationIDsForDataScope = RemoveDuplication(tempOrganizationIDs)
	return
}

// GetOrganizationIDsWithInferiors 获得所有的组织id(含子组织)
func GetOrganizationIDsWithInferiors(userID int) (organizationIDs []int) {
	global.DB.Model(&model.OrganizationAndUser{}).Where("user_id = ?", userID).
		Select("organization_id").Find(&organizationIDs)
	for i := range organizationIDs {
		res := getInferiorOrganizationIDs(organizationIDs[i])
		organizationIDs = append(organizationIDs, res...)
	}
	organizationIDs = RemoveDuplication(organizationIDs)
	return
}

func getInferiorOrganizationIDs(organizationID int) (organizationIDs []int) {
	global.DB.Model(&model.Organization{}).Where("superior_id = ?", organizationID).
		Select("id").Find(&organizationIDs)
	for i := range organizationIDs {
		res := getInferiorOrganizationIDs(organizationIDs[i])
		organizationIDs = append(organizationIDs, res...)
	}
	return
}
