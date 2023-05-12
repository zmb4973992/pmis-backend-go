package model

import (
	"errors"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"time"
)

//不要软删除，因为progress涉及到大量的删除、新增，会产生大量的冗余数据

type BasicModel struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;size:0"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime;size:0"`
	Creator      *int64    `json:"creator"`
	LastModifier *int64    `json:"last_modifier"`
}

type IModel interface {
	TableName() string
}

func (b *BasicModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = idgen.NextId()
	if b.ID == 0 {
		return errors.New("生成id失败")
	}
	return nil
}
