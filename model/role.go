package model

import (
	"gorm.io/gorm"
)

type Role struct {
	BasicModel
	Name            string //角色名称
	SuperiorID      *int64 //上级角色id
	DataAuthorityID int64  //数据权限id
}

// TableName 修改表名
func (r *Role) TableName() string {
	return "role"
}

func (r *Role) BeforeDelete(tx *gorm.DB) error {
	if r.ID == 0 {
		return nil
	}

	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []UserAndRole
	err = tx.Where(UserAndRole{RoleID: r.ID}).
		Find(&records).Delete(&records).Error
	if err != nil {
		return err
	}
	return nil
}
