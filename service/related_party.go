package service

import (
	"github.com/mitchellh/mapstructure"
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
)

/*
Service层没有数据结构、只有方法，所有的数据结构都放在DTO里
入参为id、DTO，出参为response。
这里的方法从controller拿来id或初步处理的入参dto，重点是处理业务逻辑。
所有的增删改查都交给DAO层处理，否则service层会非常庞大。
生成出参response后，交给controller展示。
*/

type relatedPartyService struct{}

func (relatedPartyService) Get(relatedPartyID int) response.Common {
	var result dto.RelatedPartyGetDTO
	err := global.DB.Model(&model.RelatedParty{}).
		Where("id = ?", relatedPartyID).First(&result).Error
	if err != nil {
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (relatedPartyService) Create(paramIn *dto.RelatedPartyCreateOrUpdateDTO) response.Common {
	//对model进行清洗，生成dao层需要的model
	var paramOut model.RelatedParty
	if paramIn.Creator != nil {
		paramOut.Creator = paramIn.Creator
	}

	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.ChineseName != "" {
		paramOut.ChineseName = paramIn.ChineseName
	}

	if *paramIn.EnglishName != "" {
		paramOut.EnglishName = paramIn.EnglishName
	}

	if *paramIn.Address != "" {
		paramOut.Address = paramIn.Address
	}

	if *paramIn.Telephone != "" {
		paramOut.Telephone = paramIn.Telephone
	}

	if *paramIn.UniformSocialCreditCode != "" {
		paramOut.UniformSocialCreditCode = paramIn.UniformSocialCreditCode
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToCreateRecord)
	}
	return response.Success()
}

func (relatedPartyService) Update(paramIn *dto.RelatedPartyCreateOrUpdateDTO) response.Common {
	var paramOut model.RelatedParty
	//先找出原始记录
	err := global.DB.Where("id = ?", paramIn.ID).First(&paramOut).Error
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}

	//对dto进行清洗，生成dao层需要的model
	if paramIn.LastModifier != nil {
		paramOut.LastModifier = paramIn.LastModifier
	}

	if *paramIn.ChineseName != "" {
		paramOut.ChineseName = paramIn.ChineseName
	}

	if *paramIn.EnglishName != "" {
		paramOut.EnglishName = paramIn.EnglishName
	}

	if *paramIn.Address != "" {
		paramOut.Address = paramIn.Address
	}

	if *paramIn.Telephone != "" {
		paramOut.Telephone = paramIn.Telephone
	}

	if *paramIn.UniformSocialCreditCode != "" {
		paramOut.UniformSocialCreditCode = paramIn.UniformSocialCreditCode
	}

	//清洗完毕，开始update
	err = global.DB.Where("id = ?", paramOut.ID).Omit("created_at", "creator").Save(&paramOut).Error
	//拿到dao层的返回结果，进行处理
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (relatedPartyService) Delete(relatedPartyID int) response.Common {
	err := global.DB.Delete(&model.RelatedParty{}, relatedPartyID).Error
	if err != nil {
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (relatedPartyService) List(paramIn dto.RelatedPartyListDTO) response.List {
	//生成sql查询条件
	sqlCondition := util.NewSqlCondition()
	//对paramIn进行清洗
	//这部分是用于where的参数
	if paramIn.Page > 0 {
		sqlCondition.Paging.Page = paramIn.Page
	}

	//如果参数里的pageSize是整数且大于0、小于等于100：
	if paramIn.PageSize > 0 && paramIn.PageSize <= 100 {
		sqlCondition.Paging.PageSize = paramIn.PageSize
	}

	if id := paramIn.ID; id > 0 {
		sqlCondition.Equal("id", id)
	}

	if paramIn.IDGte != nil {
		sqlCondition.Gte("id", *paramIn.IDGte)
	}

	if paramIn.IDLte != nil {
		sqlCondition.Lte("id", *paramIn.IDLte)
	}

	if paramIn.ChineseName != nil && *paramIn.ChineseName != "" {
		sqlCondition = sqlCondition.Equal("chinese_name", *paramIn.ChineseName)
	}

	if paramIn.ChineseNameLike != nil && *paramIn.ChineseNameLike != "" {
		sqlCondition = sqlCondition.Like("chinese_name", *paramIn.ChineseNameLike)
	}

	if paramIn.EnglishNameLike != nil && *paramIn.EnglishNameLike != "" {
		sqlCondition = sqlCondition.Like("english_name", *paramIn.EnglishNameLike)
	}

	//这部分是用于order的参数
	orderBy := paramIn.OrderBy
	if orderBy != "" {
		ok := sqlCondition.ValidateColumn(orderBy, model.RelatedParty{})
		if ok {
			sqlCondition.Sorting.OrderBy = orderBy
		}
	}
	desc := paramIn.Desc
	if desc == true {
		sqlCondition.Sorting.Desc = true
	} else {
		sqlCondition.Sorting.Desc = false
	}

	tempList := sqlCondition.Find(global.DB, model.RelatedParty{})
	totalRecords := sqlCondition.Count(global.DB, model.RelatedParty{})
	totalPages := util.GetTotalPages(totalRecords, sqlCondition.Paging.PageSize)

	if len(tempList) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	var list []dto.RelatedPartyGetDTO
	_ = mapstructure.Decode(&tempList, &list)

	//处理字段类型不匹配、或者有特殊格式要求的字段
	//for k := range tempList {
	//	a := tempList[k]["date"].(*time.Time).Format("2006-01-02")
	//	list[k].Date = &a
	//}

	return response.List{
		Data: list,
		Paging: &dto.PagingDTO{
			Page:         sqlCondition.Paging.Page,
			PageSize:     sqlCondition.Paging.PageSize,
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
