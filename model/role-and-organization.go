package model

// RoleAndOrganization 组织机构和数据权限(范围)的中间表
// 当角色的数据范围为自定义时，系统会到这个表来查询，具体能访问哪些组织的数据
type RoleAndOrganization struct {
	BasicModel
	RoleSnowID         *int64 `gorm:"nut null;"`
	OrganizationSnowID *int64 `gorm:"nut null;"` //等同于组织SnowID，用来定义可以查看哪些组织的数据
}

// TableName 修改表名
func (*RoleAndOrganization) TableName() string {
	return "role_and_organization"
}

//func (d *Role) BeforeDelete(tx *gorm.DB) error {
//	//删除相关的子表记录
//	//先find，再delete，可以激活相关的钩子函数
//	var records []UserAndRole
//	err = tx.Where("role_id = ?", d.SnowID).
//		Find(&records).Delete(&records).Error
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func generateRoles() error {
//	roles := []Role{
//		{Name: "管理员"},
//		{Name: "公司级"},
//		{Name: "事业部级"},
//		{Name: "部门级"},
//		{Name: "项目级"},
//	}
//	for _, role := range roles {
//		err := global.DB.FirstOrCreate(&Role{}, role).Error
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
