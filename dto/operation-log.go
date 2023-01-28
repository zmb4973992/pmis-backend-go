package dto

import "time"

//以下为入参

type OperationLogCreateOrUpdate struct {
	Base
	UserID              *int       `json:"user_id"`               //操作人id
	IP                  *string    `json:"ip"`                    //IP
	Location            *string    `json:"location"`              //所在地
	Method              *string    `json:"method"`                //请求方式
	Path                *string    `json:"path"`                  //请求路径
	Remarks             *string    `json:"remarks"`               //备注
	ResponseCode        *int       `json:"response_code"`         //响应码
	StartTime           *time.Time `json:"start_time"`            //发起时间
	MilliSecondsElapsed *int       `json:"milli_seconds_elapsed"` //处理耗时（毫秒）
	UserAgent           *string    `json:"user_agent"`            //浏览器标识
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

//以下为出参

type OperationLogOutput struct {
	//Base `mapstructure:",squash"` //这里是嵌套结构体，mapstructure必须加squash，否则无法匹配
	UserID              *int       `json:"user_id"`               //操作人id
	IP                  *string    `json:"ip"`                    //IP
	Location            *string    `json:"location"`              //所在地
	Method              *string    `json:"method"`                //请求方式
	Path                *string    `json:"path"`                  //请求路径
	Remarks             *string    `json:"remarks"`               //备注
	ResponseCode        *int       `json:"response_code"`         //响应码
	StartTime           *time.Time `json:"start_time"`            //发起时间
	MilliSecondsElapsed *int       `json:"milli_seconds_elapsed"` //处理耗时（毫秒）
	UserAgent           *string    `json:"user_agent"`            //浏览器标识
}
