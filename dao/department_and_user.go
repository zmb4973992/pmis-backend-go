package dao

import (
	"learn-go/global"
	"learn-go/model"
)

type departmentAndUserDAO struct{}

// Create 这里是只负责新增，不写任何业务逻辑。只要收到参数就创建数据库记录，然后返回错误
func (departmentAndUserDAO) Create(param *model.DepartmentAndUser) error {
	err := global.DB.Create(param).Error
	return err
}

func (departmentAndUserDAO) CreateBatch(param []model.DepartmentAndUser) error {
	err := global.DB.Create(param).Error
	return err
}
