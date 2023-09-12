package service

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
	"strconv"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type UserLogin struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CaptchaId string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

type UserGet struct {
	Id int64
}

type UserCreate struct {
	UserId            int64
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
	UserId            int64
	Id                int64
	FullName          *string `json:"full_name"`           //全名
	EmailAddress      *string `json:"email_address"`       //邮箱地址
	IsValid           *bool   `json:"is_valid"`            //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"` //手机号
	EmployeeNumber    *string `json:"employee_number"`     //工号
}

type UserDelete struct {
	Id int64
}

type UserGetList struct {
	list.Input
	IsValid         *bool  `json:"is_valid"`
	UsernameInclude string `json:"username_include,omitempty"`
	RoleId          int64  `json:"role_id,omitempty"`
}

type UserUpdateRoles struct {
	LastModifier int64

	UserId  int64    `json:"-"`
	RoleIds *[]int64 `json:"role_ids"`
}

type UserUpdateDataAuthority struct {
	LastModifier int64

	UserId          int64 `json:"-"`
	DataAuthorityId int64 `json:"data_authority_id" binding:"required"`
}

//以下为出参

type UserOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`

	Username          string  `json:"username"`            //用户名
	FullName          *string `json:"full_name"`           //全名
	EmailAddress      *string `json:"email_address"`       //邮箱地址
	IsValid           *bool   `json:"is_valid"`            //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"` //手机号
	EmployeeNumber    *string `json:"employee_number"`     //工号
}

func (u *UserLogin) Verify() bool {
	store := base64Captcha.DefaultMemStore
	permitted := store.Verify(u.CaptchaId, u.Captcha, true)
	return permitted
}

func (u *UserLogin) Login() (output any, errCode int) {
	permitted, err := util.LoginByLDAP(u.Username, u.Password)

	if err != nil || !permitted {
		return nil, util.ErrorInvalidUsernameOrPassword
	}

	var user UserOutput
	err = global.DB.Model(model.User{}).
		Where("username = ?", u.Username).
		First(&user).Error
	if err != nil {
		return nil, util.ErrorInvalidUsernameOrPassword
	}

	token, err1 := util.GenerateToken(user.Id)
	if err1 != nil {
		return nil, util.ErrorFailToGenerateToken
	}

	return gin.H{"access_token": token},
		util.Success
}

func (u *UserGet) Get() (output *UserOutput, errCode int) {
	err := global.DB.Model(model.User{}).
		Where("id = ?", u.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	return output, util.Success
}

func (u *UserCreate) Create() (errCode int) {
	var paramOut model.User
	if u.UserId > 0 {
		paramOut.Creator = &u.UserId
	}

	paramOut.Username = u.Username
	encryptedPassword, err := util.Encrypt(u.Password)
	if err != nil {
		return util.ErrorFailToEncrypt
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
		return util.ErrorFailToCreateRecord
	}
	return util.Success
}

func (u *UserUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	if u.UserId > 0 {
		paramOut["last_modifier"] = u.UserId
	}

	if u.FullName != nil {
		if *u.FullName != "" {
			paramOut["full_name"] = u.FullName
		} else {
			return util.ErrorInvalidJSONParameters
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

	err := global.DB.Model(&model.User{}).
		Where("id = ?", u.Id).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}

func (u *UserDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.User
	err := global.DB.Where("id = ?", u.Id).
		Find(&record).
		Delete(&record).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	return util.Success
}

func (u *UserGetList) GetList() (outputs []UserOutput,
	errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.User{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where
	if u.IsValid != nil {
		db = db.Where("is_valid = ?", *u.IsValid)
	}

	if u.UsernameInclude != "" {
		db = db.Where("username like ?", "%"+u.UsernameInclude+"%")
	}

	if u.RoleId > 0 {
		var userIds []int64
		global.DB.Model(&model.UserAndRole{}).
			Where("role_id = ?", u.RoleId).
			Select("user_id").
			Find(&userIds)
		db = db.Where("id in ?", userIds)
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
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.User{}, orderBy)
		if !exists {
			return nil, util.ErrorSortingFieldDoesNotExist, nil
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

	//outputs
	db.Model(&model.User{}).Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return outputs,
		util.Success,
		&list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		}
}

func (u *UserUpdateRoles) Update() (errCode int) {
	if u.RoleIds == nil {
		return util.ErrorInvalidJSONParameters
	}

	if len(*u.RoleIds) == 0 {
		err := global.DB.Where("user_id = ?", u.UserId).
			Delete(&model.UserAndRole{}).Error
		if err != nil {
			return util.ErrorFailToDeleteRecord
		}
		return util.Success
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		//先删掉原始记录
		err := tx.Where("user_id = ?", u.UserId).
			Delete(&model.UserAndRole{}).Error
		if err != nil {
			return ErrorFailToDeleteRecord
		}

		//再增加新的记录
		var paramOut []model.UserAndRole
		for _, roleId := range *u.RoleIds {
			var record model.UserAndRole
			if u.LastModifier > 0 {
				record.LastModifier = &u.LastModifier
			}

			record.UserId = u.UserId
			record.RoleId = roleId
			paramOut = append(paramOut, record)
		}

		err = global.DB.Create(&paramOut).Error
		if err != nil {
			return ErrorFailToCreateRecord
		}

		//更新casbin的rbac分组规则
		var param1 rbacUpdateGroupingPolicyByMember
		param1.Member = strconv.FormatInt(u.UserId, 10)
		for _, roleId := range *u.RoleIds {
			param1.Groups = append(param1.Groups, strconv.FormatInt(roleId, 10))
		}
		err = param1.Update()
		if err != nil {
			return ErrorFailToUpdateRBACGroupingPolicies
		}

		return nil
	})

	switch err {
	case nil:
		return util.Success
	case ErrorFailToCreateRecord:
		return util.ErrorFailToCreateRecord
	case ErrorFailToDeleteRecord:
		return util.ErrorFailToDeleteRecord
	case ErrorFieldsToBeCreatedNotFound:
		return util.ErrorFieldsToBeCreatedNotFound
	case ErrorFailToUpdateRBACGroupingPolicies:
		return util.ErrorFailToUpdateRBACGroupingPolicies
	default:
		return util.ErrorFailToUpdateRecord
	}
}

func (u *UserUpdateDataAuthority) Update() (errCode int) {
	paramOut := make(map[string]any)
	if u.LastModifier > 0 {
		paramOut["last_modifier"] = u.LastModifier
	}

	paramOut["data_authority_id"] = u.DataAuthorityId

	err := global.DB.Model(&model.UserAndDataAuthority{}).
		Where("user_id = ?", u.UserId).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}
