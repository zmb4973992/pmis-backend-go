package service

import (
	"gorm.io/gorm"
	"learn-go/dao"
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
	"learn-go/serializer/response"
	"learn-go/util"
)

// UserService 没有数据、只有方法，所有的数据都放在DTO里
//这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
//所有的增删改查都交给DAO层处理，否则service层会非常庞大
type roleAndUserService struct{}

func (roleAndUserService) Get(userID int) response.Common {
	result := dao.UserDAO.Get(userID)
	if result == nil {
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (roleAndUserService) Create(paramIn *dto.UserCreateDTO) response.Common {
	//对数据进行清洗
	var paramOut model.User
	paramOut.Username = paramIn.Username
	//对密码进行加密
	encryptedPassword, err := util.EncryptPassword(paramIn.Password)
	if err != nil {
		return response.Failure(util.ErrorFailToEncrypt)
	}
	paramOut.Password = encryptedPassword
	paramOut.IsValid = paramIn.IsValid
	if *paramIn.FullName == "" {
		paramOut.FullName = nil
	} else {
		paramOut.FullName = paramIn.FullName
	}
	if *paramIn.EmailAddress == "" {
		paramOut.EmailAddress = nil
	} else {
		paramOut.EmailAddress = paramIn.EmailAddress
	}
	if *paramIn.MobilePhoneNumber == "" {
		paramOut.MobilePhoneNumber = nil
	} else {
		paramOut.MobilePhoneNumber = paramIn.MobilePhoneNumber
	}
	if *paramIn.EmployeeNumber == "" {
		paramOut.EmployeeNumber = nil
	} else {
		paramOut.EmployeeNumber = paramIn.EmployeeNumber
	}
	//这里对一对多关系字段不作处理，都交给下面的事务

	//生成子表需要的角色数据
	for i := range paramIn.Roles {
		paramOut.Roles = append(paramOut.Roles, model.RoleAndUser{
			RoleID: &paramIn.Roles[i],
		})
	}
	//生成子表需要的部门数据
	for i := range paramIn.Departments {
		paramOut.Departments = append(paramOut.Departments, model.DepartmentAndUser{
			DepartmentID: &paramIn.Departments[i],
		})
	}

	err = dao.UserDAO.Create(&paramOut)

	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (roleAndUserService) Update(paramIn *dto.UserUpdateDTO) response.Common {
	var paramOut model.User

	//先找出原始记录
	err := global.DB.Where("id = ?", paramIn.ID).First(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if *paramIn.FullName == "" {
		paramOut.FullName = nil
	} else {
		paramOut.FullName = paramIn.FullName
	}
	if *paramIn.EmailAddress == "" {
		paramOut.EmailAddress = nil
	} else {
		paramOut.EmailAddress = paramIn.EmailAddress
	}
	paramOut.IsValid = paramIn.IsValid
	if *paramIn.MobilePhoneNumber == "" {
		paramOut.MobilePhoneNumber = nil
	} else {
		paramOut.MobilePhoneNumber = paramIn.MobilePhoneNumber
	}
	if *paramIn.EmployeeNumber == "" {
		paramOut.EmployeeNumber = nil
	} else {
		paramOut.EmployeeNumber = paramIn.EmployeeNumber
	}

	//由于涉及到多表的保存，这里对一对多关系字段不作处理，都交给下面的事务
	err = global.DB.Transaction(
		func(tx *gorm.DB) error {
			//注意，这里没有使用dao层的封装方法，而是使用tx+gorm的原始方法
			err := tx.Where("id = ?", paramIn.ID).Omit("created_at").Save(&paramOut).Error
			if err != nil {
				return err
			}
			//把用户-角色的对应关系添加到role_and_user表
			//如果有角色数据：
			if len(paramIn.Roles) > 0 {
				//获取原始的角色数据：
				var existedRoleIDs []int
				tx.Model(&model.RoleAndUser{}).Select("role_id").Where("user_id = ?", paramIn.ID).Find(&existedRoleIDs)
				//新老数据比较
				ok := util.SlicesAreSame(paramIn.Roles, existedRoleIDs)
				//如果不相同，则开始更新
				if !ok {
					//先把中间表的数据删除
					tx.Where("user_id = ?", paramIn.ID).Delete(&model.RoleAndUser{})
					//然后插入新的中间表数据
					var paramOutForRoleAndUser []model.RoleAndUser
					//这里不能使用v进行循环赋值，因为涉及到指针，会导致所有记录都变成一样的
					for k := range paramIn.Roles {
						var record model.RoleAndUser
						record.UserID = &paramOut.ID
						record.RoleID = &paramIn.Roles[k]
						paramOutForRoleAndUser = append(paramOutForRoleAndUser, record)
					}
					err = tx.Create(&paramOutForRoleAndUser).Error
					if err != nil {
						return err
					}
				}
			}

			//把用户-部门的对应关系添加到department_and_user表
			//如果有部门数据：
			if len(paramIn.Departments) > 0 {
				//获取原始的部门数据：
				var existedDepartmentIDs []int
				tx.Model(&model.DepartmentAndUser{}).Select("department_id").Where("user_id = ?", paramIn.ID).Find(&existedDepartmentIDs)
				//新老数据比较
				ok := util.SlicesAreSame(paramIn.Roles, existedDepartmentIDs)
				//如果不相同，则开始更新
				if !ok {
					//先把中间表的数据删除
					tx.Where("user_id = ?", paramIn.ID).Delete(&model.DepartmentAndUser{})
					//然后插入新的中间表数据
					var paramOutForDepartmentAndUser []model.DepartmentAndUser
					for k := range paramIn.Departments {
						var record model.DepartmentAndUser
						record.UserID = &paramOut.ID
						record.DepartmentID = &paramIn.Departments[k]
						paramOutForDepartmentAndUser = append(paramOutForDepartmentAndUser, record)
					}
					err = tx.Create(&paramOutForDepartmentAndUser).Error
					if err != nil {
						return err
					}
				}
			}
			//事务执行完毕,返回空则自动提交
			return nil
		})

	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (roleAndUserService) Delete(userID int) response.Common {
	//新建一个dao.User结构体的实例
	err := dao.UserDAO.Delete(userID)
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (roleAndUserService) List(param dto.RoleAndUserListDTO) []dto.RoleAndUserGetDTO {
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

	if len(paramPairs) == 0 {
		return nil
	}

	res := dao.RoleAndUserDAO.List(paramPairs)
	return res
}
