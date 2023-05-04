package model

import "time"

// WorkReview 工作点评
type WorkReview struct {
	BasicModel
	ProjectSnowID  *uint64    //项目SnowID
	Content        *string    //内容
	ExpirationDate *time.Time `gorm:"type:datetime"` //失效日期，有效期截止
}

// TableName 修改表名
func (*WorkReview) TableName() string {
	return "work_review"
}
