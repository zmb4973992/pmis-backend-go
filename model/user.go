package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	BaseModel
	Username          string
	Password          string
	IsValid           *bool   //用户为有效还是禁用
	FullName          *string //全名
	EmailAddress      *string //邮箱地址
	MobilePhoneNumber *string //手机号
	EmployeeNumber    *string //工号
	Role              []*Role `gorm:"many2many:user_role;"`
}

// TableName 将表名改为user
func (*User) TableName() string {
	return "user"
}

type UserRole struct {
	ID        int
	CreatedAt time.Time
	UserID    int `gorm:"primaryKey"`
	RoleID    int `gorm:"primaryKey"`
}

func (d *User) BeforeDelete(tx *gorm.DB) error {
	if d.ID > 0 {
		//如果有删除人的id，则记录下来
		if d.Deleter != nil && *d.Deleter > 0 {
			err := tx.Model(&User{}).Where("id = ?", d.ID).
				Update("deleter", d.Deleter).Error
			if err != nil {
				return err
			}
		}
		//删除相关的子表记录
		err = tx.Model(&RoleAndUser{}).Where("user_id = ?", d.ID).
			Updates(map[string]any{
				"deleted_at": time.Now(),
				"deleter":    d.Deleter,
			}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&DepartmentAndUser{}).Where("user_id = ?", d.ID).
			Updates(map[string]any{
				"deleted_at": time.Now(),
				"deleter":    d.Deleter,
			}).Error
		if err != nil {
			return err
		}

		err = tx.Model(&ProjectAndUser{}).Where("user_id = ?", d.ID).
			Updates(map[string]any{
				"deleted_at": time.Now(),
				"deleter":    d.Deleter,
			}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
