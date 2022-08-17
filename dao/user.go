package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
)

type userDAO struct{}

func (userDAO) Get(userID int) *dto.UserGetDTO {
	var userGetDTO dto.UserGetDTO
	//把基础的账号信息查出来
	var user model.User
	err := global.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil
	}
	//把该userID的所有role_and_user记录查出来
	var roleAndUsers []model.RoleAndUser
	global.DB.Where("user_id = ?", userID).Find(&roleAndUsers)
	//把所有的roleID提取出来，查出相应的角色名称
	var roleNames []string
	for _, roleAndUser := range roleAndUsers {
		var role model.Role
		global.DB.Where("id = ?", roleAndUser.RoleID).First(&role)
		roleNames = append(roleNames, role.Name)
	}
	//把该userID的所有department_and_user记录查出来
	var departmentAndUsers []model.DepartmentAndUser
	global.DB.Where("user_id = ?", userID).Find(&departmentAndUsers)
	//把所有的departmentID提取出来，查出相应的部门名称
	var departmentNames []string
	for _, departmentAndUser := range departmentAndUsers {
		var department model.Department
		global.DB.Where("id = ?", departmentAndUser.DepartmentID).First(&department)
		departmentNames = append(departmentNames, department.Name)
	}
	//把所有查出的结果赋值给输出变量
	userGetDTO.Username = user.Username
	if user.IsValid != nil {
		userGetDTO.IsValid = user.IsValid
	}
	if user.FullName != nil {
		userGetDTO.FullName = user.FullName
	}
	if user.EmailAddress != nil {
		userGetDTO.EmailAddress = user.EmailAddress
	}
	if user.MobilePhoneNumber != nil {
		userGetDTO.MobilePhoneNumber = user.MobilePhoneNumber
	}
	if user.EmployeeNumber != nil {
		userGetDTO.EmployeeNumber = user.EmployeeNumber
	}

	userGetDTO.Roles = roleNames
	userGetDTO.Departments = departmentNames
	return &userGetDTO
}

func (userDAO) Create(param *model.User) error {
	err := global.DB.Create(param).Error
	return err
}

func (userDAO) Delete(userID int) error {
	//注意，这里就算没有找到记录，也不会报错。详见gorm的delete用法
	err := global.DB.Delete(&model.User{}, userID).Error
	return err
}
