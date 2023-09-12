package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/util"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

type OperationLogGet struct {
	Id     int64
	UserId int64
}

type OperationLogCreate struct {
	Creator int64

	Operator      int64  `json:"operator,omitempty"`
	ProjectId     int64  `json:"project_id,omitempty"`
	Date          string `json:"date,omitempty"`
	OperationType int64  `json:"operation_type,omitempty"`
	Detail        string `json:"detail,omitempty"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type OperationLogDelete struct {
	OperationLogId int64
	UserId         int64
}

type OperationLogGetList struct {
	list.Input
	ProjectId int64 `json:"project_id,omitempty"`
}

//以下为出参

type OperationLogOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	Id           int64  `json:"id"`
	//连接关联表的id，只用来给gorm查询，不在json中显示
	ProjectId *int64 `json:"-"`
	Operator  *int64 `json:"-"`
	//连接dictionary_item表的id，只用来给gorm查询，不在json中显示
	OperationType *int64 `json:"-"`
	//关联表的详情，不需要gorm查询，需要在json中显示
	ProjectExternal  *ProjectOutput `json:"project" gorm:"-"`
	OperatorExternal *UserOutput    `json:"operator" gorm:"-"`
	//dictionary_item表的详情，不需要gorm查询，需要在json中显示
	OperationTypeExternal *DictionaryDetailOutput `json:"operation_type" gorm:"-"`
	//其他属性
	Date   *string `json:"date"`
	Detail *string `json:"detail"`
}

func (o *OperationLogGet) Get() (output *OperationLogOutput, errCode int) {
	err := global.DB.Model(model.OperationLog{}).
		Where("id = ?", o.Id).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	//查询关联表的详情
	{
		//查项目信息
		if output.ProjectId != nil {
			var record ProjectOutput
			res := global.DB.Model(&model.Project{}).
				Where("id = ?", *output.ProjectId).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.ProjectExternal = &record
			}
		}
		//查用户信息
		if output.Operator != nil {
			var record UserOutput
			res := global.DB.Model(&model.User{}).
				Where("id = ?", *output.Operator).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.OperatorExternal = &record
			}
		}

	}

	//查询dictionary_item表的详情
	{
		if output.OperationType != nil {
			var record DictionaryDetailOutput
			res := global.DB.Model(&model.DictionaryDetail{}).
				Where("id = ?", *output.OperationType).
				Limit(1).
				Find(&record)
			if res.RowsAffected > 0 {
				output.OperationTypeExternal = &record
			}
		}
	}

	//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
	//需要取年月日(即前9位)
	{
		if output.Date != nil {
			temp := *output.Date
			*output.Date = temp[:10]
		}
	}

	return output, util.Success
}

func (o *OperationLogCreate) Create() (errCode int) {
	var paramOut model.OperationLog

	if o.Creator > 0 {
		paramOut.Creator = &o.Creator
	}

	if o.Operator > 0 {
		paramOut.Operator = &o.Operator
	}

	//连接关联表的id
	{
		if o.ProjectId > 0 {
			paramOut.ProjectId = &o.ProjectId
		}
	}

	//连接dictionary_item表的id
	{
		if o.OperationType > 0 {
			paramOut.OperationType = &o.OperationType
		}
	}

	//日期
	{
		if o.Date != "" {
			date, err := time.Parse("2006-01-02", o.Date)
			if err != nil {
				return util.ErrorInvalidDateFormat
			}
			paramOut.Date = &date
		}

	}

	if o.Detail != "" {
		paramOut.Detail = &o.Detail
	}

	err := global.DB.Create(&paramOut).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}

	return util.Success
}

func (o *OperationLogDelete) Delete() (errCode int) {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录
	var record model.OperationLog
	err := global.DB.Where("id = ?", o.OperationLogId).
		Find(&record).
		Delete(&record).Error

	if err != nil {
		return util.ErrorFailToDeleteRecord
	}

	return util.Success
}

func (c *OperationLogGetList) GetList() (
	outputs []OperationLogOutput, errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.OperationLog{})
	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	//where
	if c.ProjectId > 0 {
		db = db.Where("project_id = ?", c.ProjectId)
	}

	//count
	var count int64
	db.Count(&count)

	//order
	orderBy := c.SortingInput.OrderBy
	desc := c.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.OperationLog{}, orderBy)
		if !exists {
			return nil, util.ErrorSortingFieldDoesNotExist, nil
		}

		//如果要求降序排列
		if desc == true {
			db = db.Order(orderBy + " desc")
		} else { //如果没有要求排序方式
			db = db.Order(orderBy)
		}
	}

	//limit
	page := 1
	if c.PagingInput.Page > 0 {
		page = c.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if c.PagingInput.PageSize != nil && *c.PagingInput.PageSize >= 0 &&
		*c.PagingInput.PageSize <= global.Config.MaxPageSize {

		pageSize = *c.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Model(&model.OperationLog{}).Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	for i := range outputs {
		//查询关联表的详情
		{
			//查项目信息
			if outputs[i].ProjectId != nil {
				var record ProjectOutput
				res := global.DB.Model(&model.Project{}).
					Where("id = ?", *outputs[i].ProjectId).
					Limit(1).
					Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].ProjectExternal = &record
				}
			}
			//查用户信息
			if outputs[i].Operator != nil {
				var record UserOutput
				res := global.DB.Model(&model.User{}).
					Where("id = ?", *outputs[i].Operator).
					Limit(1).
					Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].OperatorExternal = &record
				}
			}
		}

		//查dictionary_item表的详情
		{
			if outputs[i].OperationType != nil {
				var record DictionaryDetailOutput
				res := global.DB.Model(&model.DictionaryDetail{}).
					Where("id = ?", *outputs[i].OperationType).
					Limit(1).
					Find(&record)
				if res.RowsAffected > 0 {
					outputs[i].OperationTypeExternal = &record
				}
			}
		}

		//处理日期，默认格式为这样的字符串：2019-11-01T00:00:00Z
		//需要取年月日(即前9位)
		{
			if outputs[i].Date != nil {
				temp := *outputs[i].Date
				*outputs[i].Date = temp[:10]
			}
		}
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return outputs,
		util.Success,
		&list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		}
}
