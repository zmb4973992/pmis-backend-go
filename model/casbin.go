package model

import (
	"pmis-backend-go/global"
)

// 各字段含义见：https://blog.csdn.net/github_34457546/article/details/108608686
// 不继承basicModel了，添加的字段暂时不会处理
// 因为添加的方法是enforcer.addPolicy()，不走gorm。以后再说

type CasbinRule struct {
	ID    int64   `json:"id"`
	PType *string `json:"p_type" gorm:"column:ptype;type:nvarchar(100)"`
	V0    *string `json:"v0" gorm:"type:nvarchar(100)"`
	V1    *string `json:"v1" gorm:"type:nvarchar(100)"`
	V2    *string `json:"v2" gorm:"type:nvarchar(100)"`
	V3    *string `json:"v3" gorm:"type:nvarchar(100)"`
	V4    *string `json:"v4" gorm:"type:nvarchar(100)"`
	V5    *string `json:"v5" gorm:"type:nvarchar(100)"`
	V6    *string `json:"v6" gorm:"type:nvarchar(100)"`
	V7    *string `json:"v7" gorm:"type:nvarchar(100)"`
}

func (c *CasbinRule) TableName() string {
	return "casbin_rule"
}

func stringToPointer(str string) *string {
	return &str
}

func generateCasbinRules() error {
	rules := []CasbinRule{
		{
			PType: stringToPointer("p"),
			V0:    stringToPointer("管理员"),
			V1:    stringToPointer("*"),
			V2:    stringToPointer("GET|POST|PUT|DELETE|PATCH"),
		},
		{
			PType: stringToPointer("p"),
			V0:    stringToPointer("公司级"),
			V1:    stringToPointer("/api/user/*"),
			V2:    stringToPointer("GET|POST|PUT|DELETE|PATCH"),
		},
	}
	for _, rule := range rules {
		err := global.DB.FirstOrCreate(&CasbinRule{}, &rule).Error
		if err != nil {
			return err
		}
	}
	return nil
}
