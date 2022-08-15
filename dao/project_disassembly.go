package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
)

type projectDisassemblyDAO struct{}

func (projectDisassemblyDAO) Get(projectDisassemblyID int) *dto.ProjectDisassemblyGetDTO {
	var param dto.ProjectDisassemblyGetDTO
	//把基础的拆解信息查出来
	var projectDisassembly model.ProjectDisassembly
	err := global.DB.Where("id = ?", projectDisassemblyID).First(&projectDisassembly).Error
	if err != nil {
		return nil
	}
	//把所有查出的结果赋值给输出变量
	if projectDisassembly.Name != nil {
		param.Name = projectDisassembly.Name
	}
	if projectDisassembly.ProjectID != nil {
		param.ProjectID = projectDisassembly.ProjectID
	}
	if projectDisassembly.Level != nil {
		param.Level = projectDisassembly.Level
	}
	if projectDisassembly.Weight != nil {
		param.Weight = projectDisassembly.Weight
	}
	if projectDisassembly.SuperiorID != nil {
		param.SuperiorID = projectDisassembly.SuperiorID
	}
	return &param
}

// Create 这里是只负责新增，不写任何业务逻辑。只要收到参数就创建数据库记录，然后返回错误
func (projectDisassemblyDAO) Create(param *model.ProjectDisassembly) error {
	err := global.DB.Create(param).Error
	return err
}

// Update 这里是只负责更新，不写任何业务逻辑。只要收到id和更新参数，然后返回错误
func (projectDisassemblyDAO) Update(param *model.ProjectDisassembly) error {
	//注意，这里就算没有找到记录，也不会报错，只有更新字段出现问题才会报错。详见gorm的update用法
	err := global.DB.Where("id = ?", param.ID).Omit("created_at").Save(param).Error
	return err
}

func (projectDisassemblyDAO) Delete(projectDisassemblyID int) error {
	//注意，这里就算没有找到记录，也不会报错。详见gorm的delete用法
	err := global.DB.Delete(&model.ProjectDisassembly{}, projectDisassemblyID).Error
	return err
}
