package model

import (
	"time"
)

type RequestLog struct {
	BasicModel
	IP       *string //IP
	Location *string //所在地
	Method   *string //请求方式
	Path     *string //请求路径
	//URIParams    gin.Params `gorm:"type:nvarchar(500)"` //URI参数
	Remarks      *string    //备注
	ResponseCode *int       //响应码
	StartTime    *time.Time `gorm:"type:datetime"` //发起时间
	TimeElapsed  *int       `gorm:"comment:aaa"`   //用时（毫秒）
	UserAgent    *string    //浏览器标识
}

// TableName 修改数据库的表名
func (r *RequestLog) TableName() string {
	return "request_log"
}
