package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type dictionaryService struct{}

func (dictionaryService) Get() response.Common {
	var data dto.DictionaryGetDTO

	var projectType []string
	global.DB.Debug().Model(&model.Dictionary{}).
		Where("project_type is not null").
		Select("project_type").Find(&projectType)
	data.ProjectType = projectType

	var province []string
	global.DB.Debug().Model(&model.Dictionary{}).
		Where("province is not null").
		Select("province").Find(&province)
	data.Province = province

	var contractType []string
	global.DB.Model(&model.Dictionary{}).
		Where("contract_type is not null").
		Select("contract_type").Find(&contractType)
	data.ContractType = contractType

	return response.Common{
		Data:    data,
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
