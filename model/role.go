package model

import (
	"gorm.io/gorm"
	"pmis-backend-go/global"
)

type Role struct {
	BasicModel
	Name             string //角色名称
	SuperiorID       *int64 //上级角色ID
	DataScopeType    int    //数据范围的类型
	NewDataScopeType string //新的数据范围类型
}

// TableName 修改表名
func (*Role) TableName() string {
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

func generateRoles() error {
	roles := []Role{
		{Name: "所有部门",
			DataScopeType:    global.AllOrganization,
			NewDataScopeType: "all_organization"},
		{Name: "所属部门和子部门",
			DataScopeType:    global.HisOrganizationAndInferiors,
			NewDataScopeType: "own_organization_and_sub_organization"},
		{Name: "所属部门",
			DataScopeType:    global.HisOrganization,
			NewDataScopeType: "own_organization"},
		//{Name: "自定义部门",
		//	DataScopeType:    global.CustomOrganization,
		//	NewDataScopeType: "custom_organization",
		//},
	}
	for _, role := range roles {
		err = global.DB.Where("name = ?", role.Name).
			Where("data_scope_type = ?", role.DataScopeType).
			Where("new_data_scope_type = ?", role.NewDataScopeType).
			FirstOrCreate(&role).Error
		if err != nil {
			return err
		}
	}
	return nil
}
