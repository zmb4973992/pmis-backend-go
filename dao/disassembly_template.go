package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
)

type disassemblyTemplateDAO struct{}

func (disassemblyTemplateDAO) Get(disassemblyTemplateID int) *dto.DisassemblyTemplateGetDTO {
	var param dto.DisassemblyTemplateGetDTO
	//把基础的拆解信息查出来
	var disassemblyTemplate model.DisassemblyTemplate
	err := global.DB.Where("id = ?", disassemblyTemplateID).First(&disassemblyTemplate).Error
	if err != nil {
		return nil
	}
	//把所有查出的结果赋值给输出变量
	if disassemblyTemplate.Name != nil {
		param.Name = disassemblyTemplate.Name
	}
	if disassemblyTemplate.ProjectID != nil {
		param.ProjectID = disassemblyTemplate.ProjectID
	}
	if disassemblyTemplate.Level != nil {
		param.Level = disassemblyTemplate.Level
	}
	if disassemblyTemplate.Weight != nil {
		param.Weight = disassemblyTemplate.Weight
	}
	if disassemblyTemplate.SuperiorID != nil {
		param.SuperiorID = disassemblyTemplate.SuperiorID
	}
	return &param
}

// Create 这里是只负责新增，不写任何业务逻辑。只要收到参数就创建数据库记录，然后返回错误
func (disassemblyTemplateDAO) Create(param *model.DisassemblyTemplate) error {
	err := global.DB.Create(param).Error
	return err
}

// Update 这里是只负责更新，不写任何业务逻辑。只要收到id和更新参数，然后返回错误
func (disassemblyTemplateDAO) Update(param *model.DisassemblyTemplate) error {
	//注意，这里就算没有找到记录，也不会报错，只有更新字段出现问题才会报错。详见gorm的update用法
	err := global.DB.Where("id = ?", param.ID).Omit("created_at", "creator").Save(param).Error
	return err
}

func (disassemblyTemplateDAO) Delete(disassemblyTemplateID int) error {
	//注意，这里就算没有找到记录，也不会报错。详见gorm的delete用法
	err := global.DB.Delete(&model.DisassemblyTemplate{}, disassemblyTemplateID).Error
	return err
}
