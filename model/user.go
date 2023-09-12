package model

import (
	"gorm.io/gorm"
)

type User struct {
	BasicModel
	Username          string
	Password          *string
	IsValid           *bool   //用户为有效还是禁用
	FullName          *string //全名
	EmailAddress      *string //邮箱地址
	MobilePhoneNumber *string //手机号
	EmployeeNumber    *string //工号
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	if u.Id == 0 {
		return nil
	}

	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []UserAndRole
	err = tx.Where(&UserAndRole{UserId: u.Id}).
		Find(&records).Delete(&records).Error
	if err != nil {
		return err
	}

	var records1 []OrganizationAndUser
	err = tx.Where(&OrganizationAndUser{UserId: u.Id}).
		Find(&records1).Delete(&records1).Error
	if err != nil {
		return err
	}
	return nil
}
