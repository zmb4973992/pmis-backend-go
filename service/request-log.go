package service

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/list"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type RequestLogGet struct {
	ID int64
}

type RequestLogDelete struct {
	ID int64
}

type RequestLogGetList struct {
	list.Input
	UserID int64 `json:"user_id,omitempty"`
}

//以下为出参

type RequestLogOutput struct {
	Creator      *int64     `json:"creator"`
	LastModifier *int64     `json:"last_modifier"`
	ID           int64      `json:"id"`
	IP           *string    `json:"ip"`            //IP
	Location     *string    `json:"location"`      //所在地
	Method       *string    `json:"method"`        //请求方式
	Path         *string    `json:"path"`          //请求路径
	Remarks      *string    `json:"remarks"`       //备注
	ResponseCode *int       `json:"response_code"` //响应码
	StartTime    *time.Time `json:"start_time"`    //发起时间
	TimeElapsed  *int       `json:"time_elapsed"`  //处理耗时（毫秒）
	UserAgent    *string    `json:"user_agent"`    //浏览器标识
}

func (o *RequestLogGet) Get() response.Common {
	var result RequestLogOutput
	err := global.DB.Model(model.RequestLog{}).
		Where("id = ?", o.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (o *RequestLogDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.RequestLog
	global.DB.Where("id = ?", o.ID).Find(&record)
	err := global.DB.Where("id = ?", o.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (o *RequestLogGetList) GetList() response.List {
	db := global.DB.Model(&model.RequestLog{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if o.UserID > 0 {
		db = db.Where("creator = ?", o.UserID)
	}

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := o.SortingInput.OrderBy
	desc := o.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.RequestLog{}, orderBy)
		if !exists {
			return response.FailureForList(util.ErrorSortingFieldDoesNotExist)
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
	if o.PagingInput.Page > 0 {
		page = o.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if o.PagingInput.PageSize != nil && *o.PagingInput.PageSize >= 0 &&
		*o.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *o.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []RequestLogOutput
	db.Model(&model.RequestLog{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return response.List{
		Data: data,
		Paging: &list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
