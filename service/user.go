package service

import (
	"github.com/mojocn/base64Captcha"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type UserLogin struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CaptchaID string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

type UserGet struct {
	ID int
}

type UserCreate struct {
	Creator           int
	LastModifier      int
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
	LastModifier      int
	ID                int
	FullName          *string `json:"full_name"`           //全名
	EmailAddress      *string `json:"email_address"`       //邮箱地址
	IsValid           *bool   `json:"is_valid"`            //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number"` //手机号
	EmployeeNumber    *string `json:"employee_number"`     //工号
}

type UserDelete struct {
	ID int
}

type UserGetList struct {
	dto.ListInput
	IsValid         *bool  `json:"is_valid"`
	UsernameInclude string `json:"username_include,omitempty"`
	RoleID          int    `json:"role_id,omitempty"`
}

//以下为出参

type UserOutput struct {
	Creator      *int `json:"creator" gorm:"creator"`
	LastModifier *int `json:"last_modifier" gorm:"last_modifier"`
	ID           int  `json:"id" gorm:"id"`

	Username          string  `json:"username" gorm:"username"`                       //用户名
	FullName          *string `json:"full_name" gorm:"full_name"`                     //全名
	EmailAddress      *string `json:"email_address" gorm:"email_address"`             //邮箱地址
	IsValid           *bool   `json:"is_valid" gorm:"is_valid"`                       //是否有效
	MobilePhoneNumber *string `json:"mobile_phone_number" gorm:"mobile_phone_number"` //手机号
	EmployeeNumber    *string `json:"employee_number" gorm:"employee_number"`         //工号
}

func (u *UserLogin) Verify() bool {
	store := base64Captcha.DefaultMemStore
	permitted := store.Verify(u.CaptchaID, u.Captcha, true)
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

	token, err1 := util.GenerateToken(user.ID)
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
		Where("id = ?", u.ID).First(&result).Error
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

	err := global.DB.Model(&model.User{}).Where("id = ?", u.ID).
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
	global.DB.Where("id = ?", u.ID).Find(&record)
	err := global.DB.Where("id = ?", u.ID).Delete(&record).Error

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

	if u.RoleID > 0 {
		var userIDs []int
		global.DB.Model(&model.UserAndRole{}).Where("role_id = ?", u.RoleID).
			Select("user_id").Find(&userIDs)
		db = db.Where("id in ?", userIDs)
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
		Paging: &dto.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
