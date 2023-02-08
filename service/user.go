package service

import (
	"github.com/mojocn/base64Captcha"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"pmis-backend-go/util/jwt"
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
	Deleter int
	ID      int
}

type UserGetList struct {
	dto.ListInput
	IsValid         *bool  `json:"is_valid"`
	UsernameInclude string `json:"username_include,omitempty"`
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

func (u *UserLogin) UserLogin() response.Common {
	var record model.User

	//根据入参的用户名，从数据库取出记录赋值给user
	err := global.DB.Where("username=?", u.Username).First(&record).Error

	//如果没有找到记录
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorInvalidUsernameOrPassword)
	}

	//如果密码错误
	if !util.CheckPassword(u.Password, record.Password) {
		return response.Fail(util.ErrorInvalidUsernameOrPassword)
	}

	//账号密码都正确时，生成token
	token, err := jwt.GenerateToken(record.ID)
	if err != nil {
		return response.Fail(util.ErrorFailToGenerateToken)
	}

	return response.SucceedWithData(
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
		return response.Fail(util.ErrorRecordNotFound)
	}
	return response.SucceedWithData(result)
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
		return response.Fail(util.ErrorFailToEncrypt)
	}

	paramOut.Password = encryptedPassword

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
		return response.Fail(util.ErrorFailToCreateRecord)
	}
	return response.Succeed()
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
			paramOut["full_name"] = nil
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
	paramOutForCounting := util.MapCopy(paramOut, "last_modifier")

	if len(paramOutForCounting) == 0 {
		return response.Fail(util.ErrorFieldsToBeUpdatedNotFound)
	}

	err := global.DB.Model(&model.User{}).Where("id = ?", u.ID).
		Updates(paramOut).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToUpdateRecord)
	}

	return response.Succeed()
}

func (u *UserDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.User
	global.DB.Where("id = ?", u.ID).Find(&record)
	record.Deleter = &u.Deleter
	err := global.DB.Where("id = ?", u.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Fail(util.ErrorFailToDeleteRecord)
	}
	return response.Succeed()
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
			return response.FailForList(util.ErrorSortingFieldDoesNotExist)
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
	if u.PagingInput.PageSize > 0 &&
		u.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = u.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []UserOutput
	db.Model(&model.User{}).Find(&data)

	if len(data) == 0 {
		return response.FailForList(util.ErrorRecordNotFound)
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
