package util

import (
	"fmt"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

func GetOrganizationSnowIDsInDataScope(userSnowID int64) (organizationSnowIDsInDataScope []int64) {
	//先获取角色
	var roleSnowIDs []int64
	global.DB.Model(&model.UserAndRole{}).Where("user_snow_id = ?", userSnowID).
		Select("role_snow_id").Find(&roleSnowIDs)

	fmt.Println("角色snow_id：", roleSnowIDs)

	//获得所属角色的数据范围类型
	var dataScopeTypes []int
	global.DB.Model(&model.Role{}).Where("snow_id in ?", roleSnowIDs).
		Select("data_scope_type").Find(&dataScopeTypes)
	fmt.Println("数据范围：", dataScopeTypes)

	var tempOrganizationSnowIDs []int64
	//判断数据范围的类型
	switch {
	//如果数据范围的类型为AllOrganization，就返回全部的组织id
	case SliceIncludes(dataScopeTypes, global.AllOrganization):
		global.DB.Model(&model.Organization{}).Select("snow_id").Find(&organizationSnowIDsInDataScope)
		return
	//如果数据范围的类型为HisOrganizationAndInferiors，就返回该条件下的所有组织id
	//并继续向下穿透执行
	case SliceIncludes(dataScopeTypes, global.HisOrganizationAndInferiors):
		tempOrganizationSnowIDs = GetOrganizationSnowIDsWithInferiors(userSnowID)
		fallthrough

	case SliceIncludes(dataScopeTypes, global.HisOrganization):
		var tempOrganizationSnowIDs1 []int64
		global.DB.Model(&model.OrganizationAndUser{}).Where("user_snow_id = ?", userSnowID).
			Select("organization_snow_id").Find(&tempOrganizationSnowIDs1)
		tempOrganizationSnowIDs = append(tempOrganizationSnowIDs, tempOrganizationSnowIDs1...)
		fallthrough

	case SliceIncludes(dataScopeTypes, global.CustomOrganization):
		var tempOrganizationSnowIDs2 []int64
		global.DB.Model(&model.RoleAndOrganization{}).Where("role_snow_id in ?", roleSnowIDs).
			Select("organization_snow_id").Find(&tempOrganizationSnowIDs2)
		tempOrganizationSnowIDs = append(tempOrganizationSnowIDs, tempOrganizationSnowIDs2...)
	}
	organizationSnowIDsInDataScope = RemoveDuplication(tempOrganizationSnowIDs)
	fmt.Println("数据范围内的组织snow_id：", organizationSnowIDsInDataScope)
	return
}

// GetOrganizationSnowIDsWithInferiors 获得所有的组织id(含子组织)
func GetOrganizationSnowIDsWithInferiors(userSnowID int64) (organizationSnowIDs []int64) {
	global.DB.Model(&model.OrganizationAndUser{}).Where("user_snow_id = ?", userSnowID).
		Select("organization_snow_id").Find(&organizationSnowIDs)
	for i := range organizationSnowIDs {
		res := getInferiorOrganizationSnowIDs(organizationSnowIDs[i])
		organizationSnowIDs = append(organizationSnowIDs, res...)
	}
	organizationSnowIDs = RemoveDuplication(organizationSnowIDs)
	return
}

func getInferiorOrganizationSnowIDs(organizationSnowID int64) (organizationSnowIDs []int64) {
	global.DB.Model(&model.Organization{}).Where("superior_snow_id = ?", organizationSnowID).
		Select("snow_id").Find(&organizationSnowIDs)
	for i := range organizationSnowIDs {
		res := getInferiorOrganizationSnowIDs(organizationSnowIDs[i])
		organizationSnowIDs = append(organizationSnowIDs, res...)
	}
	return
}
