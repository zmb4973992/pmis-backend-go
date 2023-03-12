package model

import (
	"time"
)

//还是不要软删除了，因为progress涉及到大量的删除、新增，会产生大量的冗余数据

type BasicModel struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime;size:0"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime;size:0"`
	//DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime;size:0"`
	Creator      *int `json:"creator"`
	LastModifier *int `json:"last_modifier"`
	//Deleter      *int `json:"deleter"`
}

// IModel 定义接口，用于sqlCondition传参
type IModel interface {
	TableName() string
}
