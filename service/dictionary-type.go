package service

type dictionaryTypeService struct{}

//
//func (dictionaryTypeService) Get() response.Common {
//	var data dto.DictionaryTypeOutputDTO
//
//	var projectType []string
//	global.DB.Model(&model.Dictionary{}).
//		Where("project_type is not null").
//		Select("project_type").Find(&projectType)
//	data.ProjectType = projectType
//
//	var province []string
//	global.DB.Model(&model.Dictionary{}).
//		Where("province is not null").
//		Select("province").Find(&province)
//	data.Province = province
//
//	var contractType []string
//	global.DB.Model(&model.Dictionary{}).
//		Where("contract_type is not null").
//		Select("contract_type").Find(&contractType)
//	data.ContractType = contractType
//
//	return response.Common{
//		Data:    data,
//		Code:    util.Success,
//		Message: util.GetMessage(util.Success),
//	}
//}
