package service

import (
	"github.com/mitchellh/mapstructure"
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
type userService struct{}

func (userService) Get(userID int) response.Common {
	result := dao.UserDAO.Get(userID)
	if result == nil {
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (userService) Create(paramIn *dto.UserCreateDTO) response.Common {
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

	//由于涉及到多表的保存，所以这里启用事务
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		//注意，这里没有使用dao层的封装方法，而是使用tx+gorm的原始方法
		err = tx.Create(&paramOut).Error
		if err != nil {
			return err
		}
		//把用户-角色的对应关系添加到role_and_user表
		//如果有角色数据：
		if len(paramIn.Roles) > 0 {
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

		//把用户-部门的对应关系添加到department_and_user表
		//如果有部门数据：
		if len(paramIn.Departments) > 0 {
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
		return nil
	})

	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (userService) Update(paramIn *dto.UserUpdateDTO) response.Common {
	var paramOutForUser model.User

	//先找出原始记录
	err := global.DB.Where("id = ?", paramIn.ID).First(&paramOutForUser).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if *paramIn.FullName == "" {
		paramOutForUser.FullName = nil
	} else {
		paramOutForUser.FullName = paramIn.FullName
	}
	if *paramIn.EmailAddress == "" {
		paramOutForUser.EmailAddress = nil
	} else {
		paramOutForUser.EmailAddress = paramIn.EmailAddress
	}
	paramOutForUser.IsValid = paramIn.IsValid
	if *paramIn.MobilePhoneNumber == "" {
		paramOutForUser.MobilePhoneNumber = nil
	} else {
		paramOutForUser.MobilePhoneNumber = paramIn.MobilePhoneNumber
	}
	if *paramIn.EmployeeNumber == "" {
		paramOutForUser.EmployeeNumber = nil
	} else {
		paramOutForUser.EmployeeNumber = paramIn.EmployeeNumber
	}

	//由于涉及到多表的保存，这里对一对多关系字段不作处理，都交给下面的事务
	err = global.DB.Transaction(
		func(tx *gorm.DB) error {
			//注意，这里没有使用dao层的封装方法，而是使用tx+gorm的原始方法
			err := tx.Where("id = ?", paramIn.ID).Omit("created_at").Save(&paramOutForUser).Error
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
						record.UserID = &paramOutForUser.ID
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
						record.UserID = &paramOutForUser.ID
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

func (userService) Delete(userID int) response.Common {
	//新建一个dao.User结构体的实例
	err := dao.UserDAO.Delete(userID)
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (userService) List(paramIn dto.UserListDTO) response.List {
	//生成sql查询条件
	sqlCondition := util.NewSqlCondition()

	//这部分是用于where的参数
	if paramIn.Page > 0 {
		sqlCondition.Paging.Page = paramIn.Page
	}

	//如果参数里的pageSize是整数且大于0、小于等于上限：
	maxPagingSize := global.Config.PagingConfig.MaxPageSize
	if paramIn.PageSize > 0 && paramIn.PageSize <= maxPagingSize {
		sqlCondition.Paging.PageSize = paramIn.PageSize
	}

	if id := paramIn.ID; id > 0 {
		sqlCondition.Equal("id", id)
	}
	if paramIn.IDGte != nil {
		sqlCondition.Gte("id", *paramIn.IDGte)
	}

	if paramIn.IDLte != nil {
		sqlCondition.Lte("id", *paramIn.IDLte)
	}

	if paramIn.IsValid != nil {
		sqlCondition.Equal("is_valid", *paramIn.IsValid)
	}

	if paramIn.Username != nil && *paramIn.Username != "" {
		sqlCondition = sqlCondition.Equal("username", *paramIn.Username)
	}

	if paramIn.UsernameInclude != nil && *paramIn.UsernameInclude != "" {
		sqlCondition = sqlCondition.Include("username", *paramIn.UsernameInclude)
	}

	//这部分是用于order的参数
	orderBy := paramIn.OrderBy
	if orderBy != "" {
		ok := sqlCondition.ValidateColumn(orderBy, model.User{})
		if ok {
			sqlCondition.Sorting.OrderBy = orderBy
		}
	}
	desc := paramIn.Desc
	if desc == true {
		sqlCondition.Sorting.Desc = true
	} else {
		sqlCondition.Sorting.Desc = false
	}

	tempList := sqlCondition.Find(model.User{})
	totalRecords := sqlCondition.Count(model.User{})
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	//这里的tempList是基于model的，不能直接传给前端，要处理成dto才行
	//如果map的字段类型和struct的字段类型不匹配，数据不会同步过来
	var list []dto.UserGetDTO
	_ = mapstructure.Decode(&tempList, &list)

	//处理字段类型不匹配、或者有特殊格式要求的字段
	for i := range tempList {
		userID := tempList[i]["id"]
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
		list[i].Roles = roleNames

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
		list[i].Departments = departmentNames
	}

	return response.List{
		Data: list,
		Paging: &dto.PagingDTO{
			Page:         sqlCondition.Paging.Page,
			PageSize:     sqlCondition.Paging.PageSize,
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
