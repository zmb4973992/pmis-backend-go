package service

import (
	"learn-go/dao"
	"learn-go/dto"
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

func (relatedPartyService) Create(paramIn *model.RelatedParty) response.Common {
	//对model进行清洗，生成dao层需要的model
	if paramIn.ChineseName != nil && *paramIn.ChineseName == "" {
		paramIn.ChineseName = nil
	}
	if paramIn.EnglishName != nil && *paramIn.EnglishName == "" {
		paramIn.EnglishName = nil
	}
	if paramIn.Address != nil && *paramIn.Address == "" {
		paramIn.Address = nil
	}
	if paramIn.Telephone != nil && *paramIn.Telephone == "" {
		paramIn.Telephone = nil
	}
	if paramIn.UniformSocialCreditCode != nil && *paramIn.UniformSocialCreditCode == "" {
		paramIn.UniformSocialCreditCode = nil
	}
	if paramIn.Code != nil && *paramIn.Code == -1 {
		paramIn.Code = nil
	}

	err := dao.RelatedPartyDAO.Create(paramIn)
	if err != nil {
		return response.Failure(util.ErrorFailToSaveRecord)
	}
	return response.Success()
}

func (relatedPartyService) Update(paramIn *model.RelatedParty) response.Common {
	//var record model.RelatedParty
	//对model进行清洗，生成dao层需要的model
	if paramIn.ChineseName != nil && *paramIn.ChineseName == "" {
		paramIn.ChineseName = nil
	}
	if paramIn.EnglishName != nil && *paramIn.EnglishName == "" {
		paramIn.EnglishName = nil
	}
	if paramIn.Address != nil && *paramIn.Address == "" {
		paramIn.Address = nil
	}
	if paramIn.Telephone != nil && *paramIn.Telephone == "" {
		paramIn.Telephone = nil
	}
	if paramIn.UniformSocialCreditCode != nil && *paramIn.UniformSocialCreditCode == "" {
		paramIn.UniformSocialCreditCode = nil
	}
	if paramIn.Code != nil && *paramIn.Code == -1 {
		paramIn.Code = nil
	}

	//清洗完毕，开始update
	err := dao.RelatedPartyDAO.Update(paramIn)
	//拿到dao层的返回结果，进行处理
	if err != nil {
		return response.Failure(util.ErrorFailToSaveRecord)
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
	if *paramIn.ChineseName != "" {
		sqlCondition = sqlCondition.Equal("chinese_name", *paramIn.ChineseName)
	}
	if *paramIn.ChineseNameInclude != "" {
		sqlCondition = sqlCondition.Include("chinese_name", *paramIn.ChineseNameInclude)
	}

	//这部分是用于order的参数
	column := paramIn.OrderBy
	//allColumns := []string{"id", "telephone", "file"}
	//re := util.IsInSlice(column, allColumns)
	if column != "" {
		sqlCondition.Sorting.OrderBy = column
	}
	desc := paramIn.Desc
	if desc == true {
		sqlCondition.Sorting.Desc = true
	} else {
		sqlCondition.Sorting.Desc = false
	}
	//新建一个dao.User结构体的实例
	list, totalPages, totalRecords := dao.RelatedPartyDAO.List(*sqlCondition)
	if list == nil {
		return response.FailureForList(util.ErrorRecordNotFound)
	}
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
