package model

// MessageAndUser 组织和用户的中间表
type MessageAndUser struct {
	BasicModel
	MessageId int64 `gorm:"nut null;"`
	UserId    int64 `gorm:"nut null;"`
	IsRead    bool
}

// TableName 修改表名
func (m *MessageAndUser) TableName() string {
	return "message_and_user"
}
