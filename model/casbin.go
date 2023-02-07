package model

import "pmis-backend-go/global"

// casbin的表由系统自建，这里只导入规则

// 各字段含义见：https://blog.csdn.net/github_34457546/article/details/108608686

type CasbinRule struct {
	ID    int
	PType *string `gorm:"column:ptype;type:nvarchar(100)"`
	V0    *string `gorm:"type:nvarchar(100)"`
	V1    *string `gorm:"type:nvarchar(100)"`
	V2    *string `gorm:"type:nvarchar(100)"`
	V3    *string `gorm:"type:nvarchar(100)"`
	V4    *string `gorm:"type:nvarchar(100)"`
	V5    *string `gorm:"type:nvarchar(100)"`
	V6    *string `gorm:"type:nvarchar(100)"`
	V7    *string `gorm:"type:nvarchar(100)"`
}

func (*CasbinRule) TableName() string {
	return "casbin_rule"
}

func stringToPointer(str string) *string {
	return &str
}

func generateCasbinRules() error {
	rules := []CasbinRule{
		{PType: stringToPointer("p"), V0: stringToPointer("管理员"), V1: stringToPointer("*"), V2: stringToPointer("GET|POST|PUT|DELETE|PATCH")},
		{PType: stringToPointer("p"), V0: stringToPointer("公司级"), V1: stringToPointer("/api/user/*"), V2: stringToPointer("GET|POST|PUT|DELETE|PATCH")},
	}
	for _, rule := range rules {
		err := global.DB.FirstOrCreate(&CasbinRule{}, &rule).Error
		if err != nil {
			return err
		}
	}
	return nil
}
