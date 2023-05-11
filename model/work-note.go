package model

import "time"

// WorkNote 工作备注
type WorkNote struct {
	BasicModel
	ProjectID  *int64     //项目ID
	Date       *time.Time `gorm:"type:date"` //日期
	Category   *string    //类型
	Subject    *string    //主题
	Content    *string    //内容
	Attachment *string    //附件
}

// TableName 修改表名
func (*WorkNote) TableName() string {
	return "work_note"
}
