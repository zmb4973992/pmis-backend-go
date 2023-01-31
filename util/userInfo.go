package util

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

// GetDepartmentIDs 获取当前用户所属的所有部门id
func GetDepartmentIDs(userID int) []int {
	//设置所属部门id数组
	var tempDepartmentIDs []int
	global.DB.Model(&model.DepartmentAndUser{}).Where("user_id = ?", userID).
		Select("department_id").Find(&tempDepartmentIDs)

	//去重
	tempDepartmentIDs = RemoveDuplication(tempDepartmentIDs)

	//校验level_name是否为部门
	var departmentIDs []int
	for _, departmentID := range tempDepartmentIDs {
		var count int64
		global.DB.Model(&model.Department{}).Where("id = ?", departmentID).
			Where("level_name = ?", "部门").Count(&count)
		if count > 0 {
			departmentIDs = append(departmentIDs, departmentID)
		}
	}
	return departmentIDs
}

// GetBusinessDivisionIDs 获取当前用户所属的所有事业部id
func GetBusinessDivisionIDs(userID int) []int {
	//设置所属事业部id数组
	var tempBusinessDivisionIDs []int
	global.DB.Model(&model.DepartmentAndUser{}).Where("user_id = ?", userID).
		Select("department_id").Find(&tempBusinessDivisionIDs)

	//去重
	tempBusinessDivisionIDs = RemoveDuplication(tempBusinessDivisionIDs)

	//校验level_name是否为事业部
	var businessDivisionIDs []int
	for _, businessDivisionID := range tempBusinessDivisionIDs {
		var count int64
		global.DB.Model(&model.Department{}).Where("id = ?", businessDivisionID).
			Where("level_name = ?", "事业部").Count(&count)
		if count > 0 {
			businessDivisionIDs = append(businessDivisionIDs, businessDivisionID)
		}
	}
	return businessDivisionIDs
}

// GetBiggestRoleName 获取当前用户的、最大权限的角色名称，没有的话就返回最小权限的角色名称
func GetBiggestRoleName(userID int) (biggestRoleName string) {
	//设置角色id数组
	var roleIDs []int
	global.DB.Model(&model.RoleAndUser{}).Where("user_id = ?", userID).
		Select("role_id").Find(&roleIDs)

	//去重
	roleIDs = RemoveDuplication(roleIDs)

	//获取角色的相关信息
	var roles []model.Role
	if len(roleIDs) > 0 {
		global.DB.Where("id in ?", roleIDs).Order("sequence desc").Find(&roles)
	}

	if len(roles) == 0 {
		var role model.Role
		global.DB.Order("sequence").First(&role)
		biggestRoleName = role.Name
	}

	//找到权限最大的角色名称
	biggestRoleName = roles[0].Name
	return biggestRoleName
}
