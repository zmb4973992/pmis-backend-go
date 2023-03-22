package service

import (
	"pmis-backend-go/dto"
	"pmis-backend-go/global"
	"pmis-backend-go/model"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/util"
	"time"
)

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type OperationLogGet struct {
	ID int
}

type OperationLogDelete struct {
	ID int
}

type OperationLogGetList struct {
	dto.ListInput
	UserID int `json:"user_id,omitempty"`
}

//以下为出参

type OperationLogOutput struct {
	Creator      *int       `json:"creator" gorm:"creator"`
	LastModifier *int       `json:"last_modifier" gorm:"last_modifier"`
	ID           int        `json:"id" gorm:"id"`
	UserID       *int       `json:"user_id" gorm:"user_id"`             //操作人id
	IP           *string    `json:"ip" gorm:"ip"`                       //IP
	Location     *string    `json:"location" gorm:"location"`           //所在地
	Method       *string    `json:"method" gorm:"method"`               //请求方式
	Path         *string    `json:"path" gorm:"path"`                   //请求路径
	Remarks      *string    `json:"remarks" gorm:"remarks"`             //备注
	ResponseCode *int       `json:"response_code" gorm:"response_code"` //响应码
	StartTime    *time.Time `json:"start_time" gorm:"start_time"`       //发起时间
	TimeElapsed  *int       `json:"time_elapsed" gorm:"time_elapsed"`   //处理耗时（毫秒）
	UserAgent    *string    `json:"user_agent" gorm:"user_agent"`       //浏览器标识
}

func (o *OperationLogGet) Get() response.Common {
	var result OperationLogOutput
	err := global.DB.Model(model.OperationLog{}).
		Where("id = ?", o.ID).First(&result).Error
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorRecordNotFound)
	}
	return response.SuccessWithData(result)
}

func (o *OperationLogDelete) Delete() response.Common {
	//先找到记录，然后把deleter赋值给记录方便传给钩子函数，再删除记录，详见：
	var record model.OperationLog
	global.DB.Where("id = ?", o.ID).Find(&record)
	err := global.DB.Where("id = ?", o.ID).Delete(&record).Error

	if err != nil {
		global.SugaredLogger.Errorln(err)
		return response.Failure(util.ErrorFailToDeleteRecord)
	}
	return response.Success()
}

func (o *OperationLogGetList) GetList() response.List {
	db := global.DB.Model(&model.OperationLog{})
	// 顺序：where -> count -> Order -> limit -> offset -> data

	//where
	if o.UserID > 0 {
		db = db.Where("user_id = ?", o.UserID)
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
		exists := util.FieldIsInModel(&model.OperationLog{}, orderBy)
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
	if o.PagingInput.PageSize >= 0 &&
		o.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = o.PagingInput.PageSize
	}
	db = db.Limit(pageSize)

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//data
	var data []OperationLogOutput
	db.Model(&model.OperationLog{}).Find(&data)

	if len(data) == 0 {
		return response.FailureForList(util.ErrorRecordNotFound)
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	return response.List{
		Data: data,
		Paging: &dto.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		},
		Code:    util.Success,
		Message: util.GetMessage(util.Success),
	}
}
