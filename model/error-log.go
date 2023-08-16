package model

import "time"

type ErrorLog struct {
	BasicModel
	Detail            *string    //详情
	Datetime          *time.Time `gorm:"type:datetime"` //日期时间
	MainCategory      *string    //主要类别
	SecondaryCategory *string    //次要类别
	IsResolved        *bool      //是否已解决
}

// TableName 修改表名
func (e *ErrorLog) TableName() string {
	return "error_log"
}
