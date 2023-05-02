package model

import (
	"gorm.io/gorm"
	"pmis-backend-go/global"
)

const (
	HisOrganization = iota + 1
	HisOrganizationAndInferiors
	AllOrganization
	CustomOrganization
)

type Role struct {
	BasicModel
	Name          string //角色名称
	SuperiorID    *int   //上级角色id
	DataScopeType int    //数据范围的类型
}

// TableName 修改表名
func (*Role) TableName() string {
	return "role"
}

func (d *Role) BeforeDelete(tx *gorm.DB) error {
	//删除相关的子表记录
	//先find，再delete，可以激活相关的钩子函数
	var records []UserAndRole
	err = tx.Where("role_id = ?", d.ID).
		Find(&records).Delete(&records).Error
	if err != nil {
		return err
	}
	return nil
}

func generateRoles() error {
	roles := []Role{
		{Name: "管理员", DataScopeType: AllOrganization},
		{Name: "公司级"},
		{Name: "事业部级"},
		{Name: "部门级"},
		{Name: "项目级"},
		{Name: "所有部门", DataScopeType: AllOrganization},
		{Name: "本部门和子部门", DataScopeType: HisOrganizationAndInferiors},
		{Name: "本部门", DataScopeType: HisOrganization},
		{Name: "自定义部门", DataScopeType: CustomOrganization},
	}
	for _, role := range roles {
		err := global.DB.FirstOrCreate(&Role{}, role).Error
		if err != nil {
			return err
		}
	}
	return nil
}
