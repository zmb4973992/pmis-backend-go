package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

type dictionaryTypeService struct{}

func (dictionaryTypeService) Create(paramIn *dto.DictionaryTypeCreateOrUpdateDTO) response.Common {
	var paramOut model.DictionaryType
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}

	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	paramOut.Name = paramIn.Name

	if *paramIn.Sort != -1 {
		paramOut.Sort = paramIn.Sort
	}

	if *paramIn.Remarks != "" {
		paramOut.Remarks = paramIn.Remarks
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (dictionaryTypeService) CreateInBatches(paramIn []dto.DictionaryTypeCreateOrUpdateDTO) response.Common {
	var paramOut []model.DictionaryType
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	for i := range paramIn {
		var record model.DictionaryType
		if paramIn[i].Creator != nil {
			record.Creator = paramIn[i].Creator
		}

		if paramIn[i].LastModifier != nil {
			record.LastModifier = paramIn[i].LastModifier
		}

		record.Name = paramIn[i].Name

		if *paramIn[i].Sort != -1 {
			record.Sort = paramIn[i].Sort
		}

		if *paramIn[i].Remarks != "" {
			record.Remarks = paramIn[i].Remarks
		}

		paramOut = append(paramOut, record)
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

// Update 更新为什么要用dto？首先因为很多数据需要绑定，也就是一定要传参；
// 其次是需要清洗
func (dictionaryTypeService) Update(paramIn *dto.DictionaryTypeCreateOrUpdateDTO) response.Common {
	var paramOut model.DictionaryType
	paramOut.ID = paramIn.ID
	//把dto的数据传递给model，由于下面的结构体字段为指针，所以需要进行处理
	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	paramOut.Name = paramIn.Name

	if *paramIn.Sort != -1 {
		paramOut.Sort = paramIn.Sort
	}

	if *paramIn.Remarks != "" {
		paramOut.Remarks = paramIn.Remarks
	}

	//清洗完毕，开始update
	err := global.DB.Omit("created_at", "creator").
		Save(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (dictionaryTypeService) Delete(dictionaryTypeID int) response.Common {
	err := global.DB.Delete(&model.DictionaryType{}, dictionaryTypeID).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}
