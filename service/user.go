package service

import (
	"github.com/mojocn/base64Captcha"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"strconv"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type UserLogin struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	CaptchaSnowID string `json:"captcha_snow_id"`
	Captcha       string `json:"captcha"`
}

type UserGet struct {
	SnowID int64
}

type UserCreate struct {
	Creator           int64
	LastModifier      int64
	Username          string `json:"username" binding:"required"`
	Password          string `json:"password" binding:"required"`
	FullName          string `json:"full_name,omitempty"`           //全名
	EmailAddress      string `json:"email_address,omitempty"`       //邮箱地址
	IsValid           *bool  `json:"is_valid"`                      //是否有效
	MobilePhoneNumber string `json:"mobile_phone_number,omitempty"` //手机号
	EmployeeNumber    string `json:"employee_number,omitempty"`     //工号
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type UserUpdate struct {
	LastModifier      int64
	SnowID            int64
	FullName          *string `json:"full_name"`           //全名
	EmailAddress      *string `json:"email_address"`       //邮箱地址
	IsValid           *bool   `json:"is_valid"`            //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"` //手机号
	EmployeeNumber    *string `json:"employee_number"`     //工号
}

type UserDelete struct {
	SnowID int64
}

type UserGetList struct {
	list.Input
	IsValid         *bool  `json:"is_valid"`
	UsernameInclude string `json:"username_include,omitempty"`
	RoleSnowID      int64  `json:"role_snow_id,omitempty"`
}

type UserUpdateRoles struct {
	Creator      int64
	LastModifier int64

	UserSnowID  int64    `json:"-"`
	RoleSnowIDs *[]int64 `json:"role_snow_ids"`
}

//以下为出参

type UserOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	SnowID       int64  `json:"snow_id"`

	Username          string  `json:"username"`            //用户名
	FullName          *string `json:"full_name"`           //全名
	EmailAddress      *string `json:"email_address"`       //邮箱地址
	IsValid           *bool   `json:"is_valid"`            //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"` //手机号
	EmployeeNumber    *string `json:"employee_number"`     //工号
}

func (u *UserLogin) Verify() bool {
	store := base64Captcha.DefaultMemStore
	permitted := store.Verify(u.CaptchaSnowID, u.Captcha, true)
	return permitted
}

func (u *UserLogin) Login() response.Common {
	permitted, err := util.LoginByLDAP(u.Username, u.Password)

	if err != nil {
		return response.Failure(util.ErrorInvalidUsernameOrPassword)
	}

	if !permitted {
		return response.Failure(util.ErrorInvalidUsernameOrPassword)
	}

	var user UserOutput
	err = global.DB.Model(model.User{}).
		Where("username = ?", u.Username).First(&user).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorInvalidUsernameOrPassword)
	}

	token, err1 := util.GenerateToken(user.SnowID)
	if err1 != nil {
		return response.Failure(util.ErrorFailToGenerateToken)
	}

	return response.SuccessWithData(
		map[string]any{
			"access_token": token,
		})
}

func (u *UserGet) Get() response.Common {
	var result UserOutput
	err := global.DB.Model(model.User{}).
		Where("snow_id = ?", u.SnowID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (u *UserCreate) Create() response.Common {
	var paramOut model.User
	if u.Creator > 0 {
		paramOut.Creator = &u.Creator
	}

	if u.LastModifier > 0 {
		paramOut.LastModifier = &u.LastModifier
	}

	paramOut.SnowID = idgen.NextId()

	paramOut.Username = u.Username
	encryptedPassword, err := util.Encrypt(u.Password)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToEncrypt)
	}

	paramOut.Password = &encryptedPassword

	if u.IsValid != nil {
		paramOut.IsValid = u.IsValid
	}

	if u.FullName != "" {
		paramOut.FullName = &u.FullName
	}

	if u.EmailAddress != "" {
		paramOut.EmailAddress = &u.EmailAddress
	}

	if u.MobilePhoneNumber != "" {
		paramOut.MobilePhoneNumber = &u.MobilePhoneNumber
	}
	if u.EmployeeNumber != "" {
		paramOut.EmployeeNumber = &u.EmployeeNumber
	}

	err = global.DB.Create(&paramOut).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (u *UserUpdate) Update() response.Common {
	paramOut := make(map[string]any)

	if u.LastModifier > 0 {
		paramOut["last_modifier"] = u.LastModifier
	}

	if u.FullName != nil {
		if *u.FullName != "" {
			paramOut["full_name"] = u.FullName
		} else {
			return response.Failure(util.ErrorInvalidJSONParameters)
		}
	}

	if u.EmailAddress != nil {
		if *u.EmailAddress != "" {
			paramOut["email_address"] = u.EmailAddress
		} else {
			paramOut["email_address"] = nil
		}
	}

	if u.IsValid != nil {
		paramOut["is_valid"] = u.IsValid
	}

	if u.MobilePhoneNumber != nil {
		if *u.MobilePhoneNumber != "" {
			paramOut["mobile_phone_number"] = u.MobilePhoneNumber
		} else {
			paramOut["mobile_phone_number"] = nil
		}
	}

	if u.EmployeeNumber != nil {
		if *u.EmployeeNumber != "" {
			paramOut["employee_number"] = u.EmployeeNumber
		} else {
			paramOut["employee_number"] = nil
		}
	}

	//计算有修改值的字段数，分别进行不同处理
	paramOutForCounting := util.MapCopy(paramOut, "Creator",
		"LastModifier", "CreateAt", "UpdatedAt")

	if len(paramOutForCounting) == 0 {
		return response.Failure(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.User{}).Where("snow_id = ?", u.SnowID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	return response.Success()
}

func (u *UserDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.User
	global.DB.Where("snow_id = ?", u.SnowID).Find(&record)
	err := global.DB.Where("snow_id = ?", u.SnowID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (u *UserGetList) GetList() response.List {
	db := global.DB.Model(&model.User{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if u.IsValid != nil {
		db = db.Where("is_valid = ?", *u.IsValid)
	}

	if u.UsernameInclude != "" {
		db = db.Where("username like ?", "%"+u.UsernameInclude+"%")
	}

	if u.RoleSnowID > 0 {
		var userSnowIDs []int64
		global.DB.Model(&model.UserAndRole{}).Where("role_snow_id = ?", u.RoleSnowID).
			Select("user_snow_id").Find(&userSnowIDs)
		db = db.Where("snow_id in ?", userSnowIDs)
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := u.SortingInput.OrderBy
	desc := u.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("snow_id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.User{}, orderBy)
		if !exists {
			return response.FailureForList(util.ErrorSortingFieldDoesNotExist)
		}
		//如果要求降序排列
		if desc == true {
			db = db.Order(orderBy + " desc")
		} else { //如果没有要求排序方式
			db = db.Order(orderBy)
		}
	}

	//limit
	page := 1
	if u.PagingInput.Page > 0 {
		page = u.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if u.PagingInput.PageSize != nil && *u.PagingInput.PageSize >= 0 &&
		*u.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *u.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []UserOutput
	db.Model(&model.User{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return response.List{
		Data: data,
		Paging: &list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}

func (u *UserUpdateRoles) Update() response.Common {
	if u.RoleSnowIDs == nil {
		return response.Failure(util.ErrorInvalidJSONParameters)
	}

	if len(*u.RoleSnowIDs) == 0 {
		err := global.DB.Where("user_snow_id = ?", u.UserSnowID).Delete(&model.UserAndRole{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return response.Failure(util.ErrorFailToDeleteRecord)
		}
		return response.Success()
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//先删掉原始记录
		err := tx.Where("user_snow_id = ?", u.UserSnowID).Delete(&model.UserAndRole{}).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToDeleteRecord
		}

		//再增加新的记录
		var paramOut []model.UserAndRole
		for _, roleID := range *u.RoleSnowIDs {
			var record model.UserAndRole
			if u.Creator > 0 {
				record.Creator = &u.Creator
			}
			if u.LastModifier > 0 {
				record.LastModifier = &u.LastModifier
			}

			record.UserSnowID = u.UserSnowID
			record.RoleSnowID = roleID
			record.SnowID = idgen.NextId()
			paramOut = append(paramOut, record)
		}

		for i := range paramOut {
			//计算有修改值的字段数，分别进行不同处理
			tempParamOut, err := util.StructToMap(paramOut[i])
			if err != nil {
				return ErrorFailToUpdateRecord
			}
			paramOutForCounting := util.MapCopy(tempParamOut,
				"Creator", "LastModifier", "CreateAt", "UpdatedAt", "SnowId")

			if len(paramOutForCounting) == 0 {
				return ErrorFieldsToBeCreatedNotFound
			}
		}

		err = global.DB.Create(&paramOut).Error
		if err != nil {
			global.SugaredLogger.Errorln(err)
			return ErrorFailToCreateRecord
		}

		//更新casbin的rbac分组规则
		var param1 rbacUpdateGroupingPolicyByMember
		param1.Member = strconv.FormatInt(u.UserSnowID, 10)
		for _, roleSnowID := range *u.RoleSnowIDs {
			param1.Groups = append(param1.Groups, strconv.FormatInt(roleSnowID, 10))
		}
		err = param1.Update()
		if err != nil {
			return ErrorFailToUpdateRBACGroupingPolicies
		}

		return nil
	})

	switch err {
	case nil:
		return response.Success()
	case ErrorFailToCreateRecord:
		return response.Failure(util.ErrorFailToCreateRecord)
	case ErrorFailToDeleteRecord:
		return response.Failure(util.ErrorFailToDeleteRecord)
	case ErrorFieldsToBeCreatedNotFound:
		return response.Failure(util.ErrorFieldsToBeCreatedNotFound)
	case ErrorFailToUpdateRBACGroupingPolicies:
		return response.Failure(util.ErrorFailToUpdateRBACGroupingPolicies)
	default:
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
}
