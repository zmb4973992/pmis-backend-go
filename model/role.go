package model

import "pmis-backend-go/global"

type Role struct {
	BaseModel
	Name     string //角色名称
	Sequence int    //顺序值，权限越大的值越大，用来比较用的
	//这里是声名外键关系，并不是实际字段。不建议用gorm的多对多的设定，不好修改
	User []RoleAndUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 修改表名
func (Role) TableName() string {
	return "role"
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
