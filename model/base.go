package model

import (
	"time"
)

//不要软删除，因为progress涉及到大量的删除、新增，会产生大量的冗余数据

type BasicModel struct {
	ID           int       `json:"id"`
	SnowID       uint64    `gorm:"not null;uniqueIndex;"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;size:0"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime;size:0"`
	Creator      *int      `json:"creator"`
	LastModifier *int      `json:"last_modifier"`
}

type IModel interface {
	TableName() string
}
