package model

import (
	"gorm.io/gorm"
	"pmis-backend-go/global"
	"time"
)

type Role struct {
	BaseModel
	Name     string  //角色名称
	Sequence int     //顺序值，权限越大的值越大，用来比较用的
	User     []*User `gorm:"many2many:user_role"`
	//这里是声名外键关系，并不是实际字段。不建议用gorm的多对多的设定，不好修改
	//User []RoleAndUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 修改表名
func (*Role) TableName() string {
	return "role"
}

func (d *Role) BeforeDelete(tx *gorm.DB) error {
	if d.ID > 0 {
		//如果有删除人的id，则记录下来
		if d.Deleter != nil && *d.Deleter > 0 {
			err := tx.Model(&Role{}).Where("id = ?", d.ID).
				Update("deleter", d.Deleter).Error
			if err != nil {
				return err
			}
		}
		//删除相关的子表记录
		err = tx.Model(&RoleAndUser{}).Where("role_id = ?", d.ID).
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
