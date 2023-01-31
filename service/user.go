package service

import (
	"github.com/mitchellh/mapstructure"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

// UserService 没有数据、只有方法，所有的数据都放在DTO里
// 这里的方法从controller拿来初步处理的入参，重点是处理业务逻辑
// 所有的增删改查都交给DAO层处理，否则service层会非常庞大
type userService struct{}

func (userService) Get(userID int) response.Common {
	var result dto.UserOutput
	//把基础的账号信息查出来
	err := global.DB.Model(model.User{}).Where("id = ?", userID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorRecordNotFound)
	}

	return response.SucceedWithData(result)
}

func (userService) Create(paramIn *dto.UserCreate) response.Common {
	//对数据进行清洗
	var paramOut model.User
	paramOut.Username = paramIn.Username
	//对密码进行加密
	encryptedPassword, err := util.EncryptPassword(paramIn.Password)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToEncrypt)
	}
	paramOut.Password = encryptedPassword
	paramOut.IsValid = paramIn.IsValid

	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}

	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.FullName != "" {
		paramOut.FullName = paramIn.FullName
	}
	if *paramIn.EmailAddress != "" {
		paramOut.EmailAddress = paramIn.EmailAddress
	}
	if *paramIn.MobilePhoneNumber != "" {
		paramOut.MobilePhoneNumber = paramIn.MobilePhoneNumber
	}
	if *paramIn.EmployeeNumber != "" {
		paramOut.EmployeeNumber = paramIn.EmployeeNumber
	}

	err = global.DB.Create(&paramOut).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
}

func (userService) Update(paramIn *dto.UserUpdate) response.Common {
	var paramOut model.User

	//先找出原始记录
	err := global.DB.Where("id = ?", paramIn.ID).First(&paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.FullName != "" {
		paramOut.FullName = paramIn.FullName
	}

	if *paramIn.EmailAddress != "" {
		paramOut.EmailAddress = paramIn.EmailAddress
	}

	paramOut.IsValid = paramIn.IsValid

	if *paramIn.MobilePhoneNumber != "" {
		paramOut.MobilePhoneNumber = paramIn.MobilePhoneNumber
	}

	if *paramIn.EmployeeNumber != "" {
		paramOut.EmployeeNumber = paramIn.EmployeeNumber
	}

	err = global.DB.Where("id = ?", paramOut.ID).
		Omit(fieldsToBeOmittedWhenUpdating...).Save(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (userService) Delete(userID int) response.Common {
	err := global.DB.Delete(&model.User{}, userID).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
}

func (userService) List(paramIn dto.UserList) response.List {
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

	if paramIn.IDGte != nil {
		sqlCondition.Gte("id", *paramIn.IDGte)
	}

	if paramIn.IDLte != nil {
		sqlCondition.Lte("id", *paramIn.IDLte)
	}

	if paramIn.IsValid != nil {
		sqlCondition.Equal("is_valid", *paramIn.IsValid)
	}

	if paramIn.UsernameLike != nil && *paramIn.UsernameLike != "" {
		sqlCondition = sqlCondition.Like("username", *paramIn.UsernameLike)
	}

	//这部分是用于order的参数
	orderBy := paramIn.OrderBy
	if orderBy != "" {
		ok := sqlCondition.FieldIsInModel(model.User{}, orderBy)
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

	tempList := sqlCondition.Find(global.DB, model.User{})
	totalRecords := sqlCondition.Count(global.DB, model.User{})
	totalPages := util.GetTotalNumberOfPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailForList(util.ErrorRecordNotFound)
	}

	//这里的tempList是基于model的，不能直接传给前端，要处理成dto才行
	//如果map的字段类型和struct的字段类型不匹配，数据不会同步过来
	var list []dto.UserOutput
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
	}

	return response.List{
		Data: list,
		Paging: &dto.PagingOutput{
			Page:            sqlCondition.Paging.Page,
			PageSize:        sqlCondition.Paging.PageSize,
			NumberOfPages:   totalPages,
			NumberOfRecords: totalRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
