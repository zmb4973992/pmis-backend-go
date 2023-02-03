package dto

import "time"

//以下为入参
//有些字段不用json tag，因为不从前端读取，而是在controller中处理

//指针字段是为了区分入参为空或0与没有入参的情况，做到分别处理，通常用于update
//如果指针字段为空或0，那么数据库相应字段会改为null；
//如果指针字段没传，那么数据库不会修改该字段

type OperationLogDelete struct {
	Deleter int
	ID      int
}

type OperationLogList struct {
	ListInput
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

type OperationRecordList struct {
	ListInput
	ID         int     `json:"id"`
	ProjectID  *int    `json:"project_id"`
	OperatorID *int    `json:"operator_id"`
	DateGte    *string `json:"date_gte"`
	DateLte    *string `json:"date_lte"`
	Action     *string `json:"action"`
}
