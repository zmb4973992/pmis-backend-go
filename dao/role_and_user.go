package dao

import (
	"fmt"
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
func (roleAndUserDAO) Create(param *model.Department) error {
	err := global.DB.Create(param).Error
	return err
}

// Update 这里是只负责更新，不写任何业务逻辑。只要收到id和更新参数，然后返回错误
func (roleAndUserDAO) Update(param *model.Department) error {
	//注意，这里就算没有找到记录，也不会报错，只有更新字段出现问题才会报错。详见gorm的update用法
	err := global.DB.Where("id = ?", param.ID).Omit("created_at").Save(param).Error
	return err
}

func (roleAndUserDAO) Delete(roleID *int, userID *int) error {
	//注意，这里就算没有找到记录，也不会报错。详见gorm的delete用法
	var sqlCondition util.SqlCondition

	//如果入参都是空，那么不做任何处理，防止勿删
	if roleID == nil && userID == nil {
		return nil
	}

	if roleID != nil {
		sqlCondition.ParamPairs = append(sqlCondition.ParamPairs, util.ParamPair{
			Key:   "role_id",
			Value: roleID,
		})
	}

	if userID != nil {
		sqlCondition.ParamPairs = append(sqlCondition.ParamPairs, util.ParamPair{
			Key:   "user_id",
			Value: userID,
		})
	}

	db := global.DB

	if len(sqlCondition.ParamPairs) > 0 {
		for _, parameterPair := range sqlCondition.ParamPairs {
			db = db.Where(parameterPair.Key, parameterPair.Value)
		}
	}

	err := db.Delete(&model.RoleAndUser{}).Error
	return err
}

func (roleAndUserDAO) List(param *dto.RoleAndUserListDTO) (list []dto.RoleAndUserGetDTO) {
	var paramPairs []util.ParamPair

	if param.RoleID != nil {
		paramPairs = append(paramPairs, util.ParamPair{
			Key:   "role_id",
			Value: param.RoleID,
		})
	}

	if param.UserID != nil {
		paramPairs = append(paramPairs, util.ParamPair{
			Key:   "user_id",
			Value: param.UserID,
		})
	}

	db := global.DB

	if len(paramPairs) > 0 {
		for _, parameterPair := range paramPairs {
			db = db.Where(parameterPair.Key, parameterPair.Value)
		}
	}

	err := db.Find(&list).Error
	if err != nil {
		fmt.Println(err)
	}

}
