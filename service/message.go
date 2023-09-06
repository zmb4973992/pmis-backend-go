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

type MessageGet struct {
	ID int64
}

type MessageCreate struct {
	UserID     int64
	Recipients []int64 `binding:"required"`
	Title      string  `binding:"required"`
	Content    string  `binding:"required"`
}

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type MessageUpdate struct {
	UserID int64
	ID     int64
}

type MessageDelete struct {
	UserID int64
	ID     int64
}

type MessageGetList struct {
	list.Input
	UserID int64
	IsRead bool `json:"is_read"`
}

type MessageOutput struct {
	Creator      *int64 `json:"creator"`
	LastModifier *int64 `json:"last_modifier"`
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Datetime     string `json:"datetime"`
}

func (m *MessageGet) Get() (output *MessageOutput, errCode int) {
	err := global.DB.Model(model.Message{}).
		Where("id = ?", m.ID).
		First(&output).Error
	if err != nil {
		return nil, util.ErrorRecordNotFound
	}

	output.Datetime = output.Datetime[:10] + " " + output.Datetime[11:19]

	return output, util.Success
}

func (m *MessageCreate) Create() (errCode int) {
	var message model.Message

	message.Creator = &m.UserID
	message.Title = m.Title
	message.Content = m.Content
	message.Datetime = time.Now()

	err := global.DB.Create(&message).Error
	if err != nil {
		return util.ErrorFailToCreateRecord
	}

	m.Recipients = util.RemoveDuplication(m.Recipients)
	var messageAndUsers []model.MessageAndUser
	for i := range m.Recipients {
		var messageAndUser model.MessageAndUser
		messageAndUser.Creator = &m.UserID
		messageAndUser.MessageID = message.ID
		messageAndUser.UserID = m.Recipients[i]
		messageAndUsers = append(messageAndUsers, messageAndUser)
	}

	global.DB.CreateInBatches(messageAndUsers, 100)

	return util.Success
}

func (m *MessageUpdate) Update() (errCode int) {
	paramOut := make(map[string]any)

	paramOut["last_modifier"] = m.UserID

	paramOut["is_read"] = true

	err := global.DB.Model(&model.MessageAndUser{}).
		Where("message_id = ?", m.ID).
		Where("user_id = ?", m.UserID).
		Updates(paramOut).Error
	if err != nil {
		return util.ErrorFailToUpdateRecord
	}

	return util.Success
}

func (m *MessageDelete) Delete() (errCode int) {
	global.DB.Where("message_id = ?", m.ID).
		Where("user_id = ?", m.UserID).
		Delete(&model.MessageAndUser{})

	return util.Success
}

func (m *MessageGetList) GetList() (
	outputs []MessageOutput, errCode int, paging *list.PagingOutput) {
	db := global.DB.Model(&model.Message{}).
		Joins("join (select distinct message_id,user_id,is_read from message_and_user where user_id = ?) as temp1 on message.id = temp1.message_id", m.UserID)

	if m.IsRead {
		db = db.Where("is_read = ?", true)
	} else {
		db = db.Where("is_read = ?", false)
	}

	// 顺序：where -> count -> Order -> limit -> offset -> outputs

	// count
	var count int64
	db.Count(&count)

	//Order
	orderBy := m.SortingInput.OrderBy
	desc := m.SortingInput.Desc
	//如果排序字段为空
	if orderBy == "" {
		//如果要求降序排列
		if desc == true {
			db = db.Order("id desc")
		}
	} else { //如果有排序字段
		//先看排序字段是否存在于表中
		exists := util.FieldIsInModel(&model.Message{}, orderBy)
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
	if m.PagingInput.Page > 0 {
		page = m.PagingInput.Page
	}
	pageSize := global.Config.DefaultPageSize
	if m.PagingInput.PageSize != nil && *m.PagingInput.PageSize >= 0 &&
		*m.PagingInput.PageSize <= global.Config.MaxPageSize {
		pageSize = *m.PagingInput.PageSize
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	//offset
	offset := (page - 1) * pageSize
	db = db.Offset(offset)

	//outputs
	db.Find(&outputs)

	if len(outputs) == 0 {
		return nil, util.ErrorRecordNotFound, nil
	}

	numberOfRecords := int(count)
	numberOfPages := util.GetNumberOfPages(numberOfRecords, pageSize)

	for i := range outputs {
		outputs[i].Datetime = outputs[i].Datetime[:10] + " " + outputs[i].Datetime[11:19]
	}

	return outputs,
		util.Success,
		&list.PagingOutput{
			Page:            page,
			PageSize:        pageSize,
			NumberOfPages:   numberOfPages,
			NumberOfRecords: numberOfRecords,
		}
}
