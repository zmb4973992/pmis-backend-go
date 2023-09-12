package model

// RoleAndOrganization 组织机构和数据权限(范围)的中间表
// 当角色的数据范围为自定义时，系统会到这个表来查询，具体能访问哪些组织的数据
type RoleAndOrganization struct {
	BasicModel
	RoleId         *int64 `gorm:"nut null;"`
	OrganizationId *int64 `gorm:"nut null;"` //等同于组织id，用来定义可以查看哪些组织的数据
}

// TableName 修改表名
func (r *RoleAndOrganization) TableName() string {
	return "role_and_organization"
}
