package model

// DepartmentAndDataRange 组织机构和数据权限(范围)的中间表
// 此表用来定义一个组织可以查看哪些组织的数据
type DepartmentAndDataRange struct {
	BasicModel
	DepartmentID *int
	DataRangeID  *int //等同于组织机构id，用来定义可以查看哪些组织的数据
}

// TableName 修改表名
func (*DepartmentAndDataRange) TableName() string {
	return "department_and_data_range"
}

//func (d *Role) BeforeDelete(tx *gorm.DB) error {
//	//删除相关的子表记录
//	//先find，再delete，可以激活相关的钩子函数
//	var records []RoleAndUser
//	err = tx.Where("role_id = ?", d.ID).
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
