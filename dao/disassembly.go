package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
)

type disassemblyDAO struct{}

func (disassemblyDAO) Get(disassemblyID int) *dto.DisassemblyGetDTO {
	var param dto.DisassemblyGetDTO
	//把基础的拆解信息查出来
	var disassembly model.Disassembly
	err := global.DB.Where("id = ?", disassemblyID).First(&disassembly).Error
	if err != nil {
		return nil
	}
	//把所有查出的结果赋值给输出变量
	if disassembly.Name != nil {
		param.Name = disassembly.Name
	}
	if disassembly.ProjectID != nil {
		param.ProjectID = disassembly.ProjectID
	}
	if disassembly.Level != nil {
		param.Level = disassembly.Level
	}
	if disassembly.Weight != nil {
		param.Weight = disassembly.Weight
	}
	if disassembly.SuperiorID != nil {
		param.SuperiorID = disassembly.SuperiorID
	}
	return &param
}

// Create 这里是只负责新增，不写任何业务逻辑。只要收到参数就创建数据库记录，然后返回错误
func (disassemblyDAO) Create(param *model.Disassembly) error {
	err := global.DB.Create(param).Error
	return err
}

// Update 这里是只负责更新，不写任何业务逻辑。只要收到id和更新参数，然后返回错误
func (disassemblyDAO) Update(param *model.Disassembly) error {
	//注意，这里就算没有找到记录，也不会报错，只有更新字段出现问题才会报错。详见gorm的update用法
	err := global.DB.Where("id = ?", param.ID).Omit("created_at", "creator").Save(param).Error
	return err
}

func (disassemblyDAO) Delete(disassemblyID int) error {
	//注意，这里就算没有找到记录，也不会报错。详见gorm的delete用法
	err := global.DB.Delete(&model.Disassembly{}, disassemblyID).Error
	return err
}
