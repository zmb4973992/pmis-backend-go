package model

import (
	"github.com/gin-gonic/gin"
	"time"
)

type OperationLog struct {
	BaseModel
	UserID       *int       //操作人id
	IP           *string    //IP
	Location     *string    //所在地
	Method       *string    //请求方式
	Path         *string    //请求路径
	URIParams    gin.Params `gorm:"type:nvarchar(max)"` //URI参数
	Remarks      *string    //备注
	ResponseCode *int       //响应码
	StartTime    *time.Time //发起时间
	TimeElapsed  *int       //处理耗时（毫秒）
	UserAgent    *string    //浏览器标识
}

// TableName 修改数据库的表名
func (*OperationLog) TableName() string {
	return "operation_log"
}
