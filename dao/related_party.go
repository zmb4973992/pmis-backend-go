package dao

import (
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
)

/*
dao层的基本原则：
入参为id或model，用于对数据库进行增删改查；
出参为err或dto，用于反馈结果或给其他层使用
*/

// RelatedPartyDAO dao层的结构体没有数据，只是操作数据库进行增删改查，不写业务逻辑
type relatedPartyDAO struct{}

// Get 这里是只负责查询，不写任何业务逻辑。
// 查询数据库记录，返回dto
func (relatedPartyDAO) Get(id int) *dto.RelatedPartyGetDTO {
	//之所以用dto不用model，是因为model为数据库原表，数据可能包含敏感字段、或未加工，不适合直接传递
	//展现的功能基本都交给dto
	var param dto.RelatedPartyGetDTO
	err := global.DB.Model(&model.RelatedParty{}).Where("id = ?", id).First(&param).Error
	if err != nil {
		return nil
	}
	return &param
}

// Create 这里是只负责新增，不写任何业务逻辑。
// 创建数据库记录，返回错误
func (relatedPartyDAO) Create(paramIn *model.RelatedParty) error {
	err := global.DB.Create(paramIn).Error
	return err
}

// Update 这里是只负责修改，不写任何业务逻辑。
// 全量更新，保存数据库记录，返回错误
func (relatedPartyDAO) Update(paramIn *model.RelatedParty) error {
	err := global.DB.Debug().Model(paramIn).Omit("created_at", "creator").Save(paramIn).Error
	return err
}

// Delete 这里是只负责删除，不写任何业务逻辑。
// 删除数据库记录，返回错误
func (relatedPartyDAO) Delete(id int) error {
	err := global.DB.Debug().Delete(&model.RelatedParty{}, id).Error
	return err
}
