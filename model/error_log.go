package model

import "time"

type ErrorLog struct {
	BaseModel
	Detail        *string    //详情
	Date          *time.Time //日期
	MajorCategory *string    //大类
	MinorCategory *string    //小类
	IsResolved    *string    //是否已解决
}

// TableName 修改表名
func (ErrorLog) TableName() string {
	return "error_log"
}
