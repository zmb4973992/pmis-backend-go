package model

import "pmis-backend-go/global"

// casbin的表由系统自建，这里只导入规则

// 各字段含义见：https://blog.csdn.net/github_34457546/article/details/108608686

type CasbinRule struct {
	ID    int
	PType *string `gorm:"column:ptype"`
	V0    *string
	V1    *string
	V2    *string
	V3    *string
	V4    *string
	V5    *string
	V6    *string
	V7    *string
}

// TableName 修改表名
func (*CasbinRule) TableName() string {
	return "casbin_rule"
}

func stringToPointer(str string) *string {
	return &str
}

func generateCasbinRules() error {
	rules := []CasbinRule{
		{PType: stringToPointer("p"), V0: stringToPointer("管理员"), V1: stringToPointer("*"), V2: stringToPointer("GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD|TRACE")},
		{PType: stringToPointer("p"), V0: stringToPointer("公司级"), V1: stringToPointer("/user/*"), V2: stringToPointer("GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD|TRACE")},
	}
	for _, rule := range rules {
		err := global.DB.FirstOrCreate(&CasbinRule{}, &rule).Error
		if err != nil {
			return err
		}
	}
	return nil
}
