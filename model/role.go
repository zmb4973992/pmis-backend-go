package model

import (
	"gorm.io/gorm"
	"pmis-backend-go/global"
)

type Role struct {
	BasicModel
	Name     string //角色名称
	Sequence int    //顺序值，权限越大的值越大，用来比较用的
}

// TableName 修改表名
func (*Role) TableName() string {
	return "role"
}

func (d *Role) BeforeDelete(tx *gorm.DB) error {
	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []RoleAndUser
	err = tx.Where("role_id = ?", d.ID).
		Find(&records).Delete(&records).Error
	if err != nil {
		return err
	}
	return nil
}

func generateRoles() error {
	roles := []Role{
		{Name: "管理员", Sequence: 10000},
		{Name: "公司级", Sequence: 9000},
		{Name: "事业部级", Sequence: 8000},
		{Name: "部门级", Sequence: 7000},
		{Name: "项目级", Sequence: 6000},
	}
	for _, role := range roles {
		err := global.DB.FirstOrCreate(&Role{}, role).Error
		if err != nil {
			return err
		}
	}
	return nil
}
