package util

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

func GetOrganizationIDsInDataScope(userID int64) (organizationIDsInDataScope []int64) {
	//先获取角色
	var roleIDs []int64
	global.DB.Model(&model.UserAndRole{}).Where("user_id = ?", userID).
		Select("role_id").Find(&roleIDs)

	fmt.Println("角色id：", roleIDs)

	//获得所属角色的数据范围类型
	var dataScopeTypes []int
	global.DB.Model(&model.Role{}).Where("id in ?", roleIDs).
		Select("data_scope_type").Find(&dataScopeTypes)
	fmt.Println("数据范围：", dataScopeTypes)

	var tempOrganizationIDs []int64
	//判断数据范围的类型
	switch {
	//如果数据范围的类型为AllOrganization，就返回全部的组织id
	case SliceIncludes(dataScopeTypes, global.AllOrganization):
		global.DB.Model(&model.Organization{}).Select("id").Find(&organizationIDsInDataScope)
		return
	//如果数据范围的类型为HisOrganizationAndInferiors，就返回该条件下的所有组织id
	//并继续向下穿透执行
	case SliceIncludes(dataScopeTypes, global.HisOrganizationAndInferiors):
		tempOrganizationIDs = GetOrganizationIDsWithInferiors(userID)
		fallthrough

	case SliceIncludes(dataScopeTypes, global.HisOrganization):
		var tempOrganizationIDs1 []int64
		global.DB.Model(&model.OrganizationAndUser{}).Where("user_id = ?", userID).
			Select("organization_id").Find(&tempOrganizationIDs1)
		tempOrganizationIDs = append(tempOrganizationIDs, tempOrganizationIDs1...)
		fallthrough

	case SliceIncludes(dataScopeTypes, global.CustomOrganization):
		var tempOrganizationIDs2 []int64
		global.DB.Model(&model.RoleAndOrganization{}).Where("role_id in ?", roleIDs).
			Select("organization_id").Find(&tempOrganizationIDs2)
		tempOrganizationIDs = append(tempOrganizationIDs, tempOrganizationIDs2...)
	}
	organizationIDsInDataScope = RemoveDuplication(tempOrganizationIDs)
	fmt.Println("数据范围内的组织id：", organizationIDsInDataScope)
	return
}

// GetOrganizationIDsWithInferiors 获得所有的组织id(含子组织)
func GetOrganizationIDsWithInferiors(userID int64) (organizationIDs []int64) {
	global.DB.Model(&model.OrganizationAndUser{}).Where("user_id = ?", userID).
		Select("organization_id").Find(&organizationIDs)
	for i := range organizationIDs {
		res := getInferiorOrganizationIDs(organizationIDs[i])
		organizationIDs = append(organizationIDs, res...)
	}
	organizationIDs = RemoveDuplication(organizationIDs)
	return
}

func getInferiorOrganizationIDs(organizationID int64) (organizationIDs []int64) {
	global.DB.Model(&model.Organization{}).Where("superior_id = ?", organizationID).
		Select("id").Find(&organizationIDs)
	for i := range organizationIDs {
		res := getInferiorOrganizationIDs(organizationIDs[i])
		organizationIDs = append(organizationIDs, res...)
	}
	return
}
