package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
	"learn-go/util"
)

type roleAndUserDAO struct{}

func (roleAndUserDAO) Get(departmentID int) *dto.DepartmentGetDTO {
	var department model.Department

	err := global.DB.Where("id = ?", departmentID).First(&department).Error
	if err != nil {
		return nil
	}

	var paramOut dto.DepartmentGetDTO
	if department.Name != "" {
		paramOut.Name = department.Name
	}

	if department.Level != "" {
		paramOut.Level = department.Level
	}

	if department.SuperiorID != nil {
		paramOut.SuperiorID = department.SuperiorID
	}

	return &paramOut
}

// Create 这里是只负责新增，不写任何业务逻辑。只要收到参数就创建数据库记录，然后返回错误
func (roleAndUserDAO) Create(param *model.RoleAndUser) error {
	err := global.DB.Create(param).Error
	return err
}

func (roleAndUserDAO) CreateInBatch(param []model.RoleAndUser) error {
	err := global.DB.Create(param).Error
	return err
}

// Update 这里是只负责更新，不写任何业务逻辑。只要收到id和更新参数，然后返回错误
func (roleAndUserDAO) Update(param *model.Department) error {
	//注意，这里就算没有找到记录，也不会报错，只有更新字段出现问题才会报错。详见gorm的update用法
	err := global.DB.Where("id = ?", param.ID).Omit("created_at").Save(param).Error
	return err
}

func (roleAndUserDAO) Delete(paramPairs []util.ParamPair) error {
	db := global.DB
	if len(paramPairs) > 0 {
		for _, parameterPair := range paramPairs {
			db = db.Where(parameterPair.Key, parameterPair.Value)
		}
	}
	err := db.Delete(&model.RoleAndUser{}).Error
	return err
}

// List 中间表为什么还要用dto、不用[]int作为结果？
//为了后期的可扩展性，万一结果格式变了，可以直接改dto
func (roleAndUserDAO) List(paramPairs []util.ParamPair) []dto.RoleAndUserGetDTO {
	db := global.DB

	for _, paramPair := range paramPairs {
		db = db.Where(paramPair.Key, paramPair.Value)
	}

	var list []dto.RoleAndUserGetDTO
	err := db.Model(&model.RoleAndUser{}).Find(&list).Error
	if err != nil {
		return nil
	}

	return list
}

//func (roleAndUserDAO) UserSlice(roleID int) []int {
//	list := RoleAndUserDAO.List(&roleID, nil)
//	var userSlice []int
//	for i := range list {
//		userSlice = append(userSlice, *list[i].UserIDs)
//	}
//	return userSlice
//}
//
//func (roleAndUserDAO) RoleSlice(userID int) []int {
//	list := RoleAndUserDAO.List(nil, &userID)
//	var roleSlice []int
//	for i := range list {
//		roleSlice = append(roleSlice, *list[i].RoleIDs)
//	}
//	return roleSlice
//}
