package dao

import (
	"learn-go/global"
	"learn-go/model"
)

type roleAndUserDAO struct{}

// Create 这里是只负责新增，不写任何业务逻辑。只要收到参数就创建数据库记录，然后返回错误
func (roleAndUserDAO) Create(param *model.RoleAndUser) error {
	err := global.DB.Create(param).Error
	return err
}

func (roleAndUserDAO) CreateBatch(param []model.RoleAndUser) error {
	err := global.DB.Create(param).Error
	return err
}
