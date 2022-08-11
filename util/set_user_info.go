package util

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
)

func SetUserInfo(c *gin.Context, userID int) {
	c.Set("user_id", userID)
	var user model.User
	//预加载全部关联信息
	global.DB.Where("id = ?", userID).Preload(clause.Associations).First(&user)
	//设置拥有的权限
	var roleNames []string
	for _, role := range user.Roles {
		var roleInfo model.Role
		global.DB.Where("id = ?", role.RoleID).First(&roleInfo)
		roleNames = append(roleNames, roleInfo.Name)
	}
	c.Set("roles", roleNames)
	//设置所属部门
	var departmentNames []string
	for _, department := range user.Departments {
		var departmentInfo model.Department
		global.DB.Where("id = ?", department.DepartmentID).First(&departmentInfo)
		departmentNames = append(departmentNames, departmentInfo.Name)
	}
	c.Set("departments", departmentNames)
}

func GetUserInfo(c *gin.Context) (user *dto.UserGetDTO) {
	return nil
}
