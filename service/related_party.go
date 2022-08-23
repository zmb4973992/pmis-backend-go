package service

import (
	"github.com/mitchellh/mapstructure"
	"learn-go/dao"
	"learn-go/dto"
	"learn-go/global"
	"learn-go/model"
	"learn-go/serializer/response"
	"learn-go/util"
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
	result := dao.RelatedPartyDAO.Get(relatedPartyID)
	if result == nil {
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

	if *paramIn.ChineseName == "" {
		paramOut.ChineseName = nil
	} else {
		paramOut.ChineseName = paramIn.ChineseName
	}

	if *paramIn.EnglishName == "" {
		paramOut.EnglishName = nil
	} else {
		paramOut.EnglishName = paramIn.EnglishName
	}

	if *paramIn.Address == "" {
		paramOut.Address = nil
	} else {
		paramOut.Address = paramIn.Address
	}

	if *paramIn.Telephone == "" {
		paramOut.Telephone = nil
	} else {
		paramOut.Telephone = paramIn.Telephone
	}

	if *paramIn.UniformSocialCreditCode == "" {
		paramOut.UniformSocialCreditCode = nil
	} else {
		paramOut.UniformSocialCreditCode = paramIn.UniformSocialCreditCode
	}

	err := dao.RelatedPartyDAO.Create(&paramOut)
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

	if *paramIn.ChineseName == "" {
		paramOut.ChineseName = nil
	} else {
		paramOut.ChineseName = paramIn.ChineseName
	}

	if *paramIn.EnglishName == "" {
		paramOut.EnglishName = nil
	} else {
		paramOut.EnglishName = paramIn.EnglishName
	}

	if *paramIn.Address == "" {
		paramOut.Address = nil
	} else {
		paramOut.Address = paramIn.Address
	}

	if *paramIn.Telephone == "" {
		paramOut.Telephone = nil
	} else {
		paramOut.Telephone = paramIn.Telephone
	}

	if *paramIn.UniformSocialCreditCode == "" {
		paramOut.UniformSocialCreditCode = nil
	} else {
		paramOut.UniformSocialCreditCode = paramIn.UniformSocialCreditCode
	}

	//清洗完毕，开始update
	err = dao.RelatedPartyDAO.Update(&paramOut)
	//拿到dao层的返回结果，进行处理
	if err != nil {
		return response.Failure(util.ErrorFailToUpdateRecord)
	}
	return response.Success()
}

func (relatedPartyService) Delete(relatedPartyID int) response.Common {
	err := dao.RelatedPartyDAO.Delete(relatedPartyID)
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

	if paramIn.IDGte != nil {
		sqlCondition.Lte("id", *paramIn.IDLte)
	}

	if paramIn.ChineseName != nil && *paramIn.ChineseName != "" {
		sqlCondition = sqlCondition.Equal("chinese_name", *paramIn.ChineseName)
	}

	if paramIn.ChineseNameInclude != nil && *paramIn.ChineseNameInclude != "" {
		sqlCondition = sqlCondition.Include("chinese_name", *paramIn.ChineseNameInclude)
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

	tempList := sqlCondition.Find(model.RelatedParty{})
	totalRecords := sqlCondition.Count(model.RelatedParty{})
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
