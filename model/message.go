package model

import "time"

type Message struct {
	BasicModel
	//连接其他表的id
	Title    string    //标题
	Content  string    //内容
	Datetime time.Time `gorm:"type:datetime"` //日期时间

}

func (m *Message) TableName() string {
	return "message"
}
