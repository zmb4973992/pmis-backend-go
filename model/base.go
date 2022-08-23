package model

import "time"

type BaseModel struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Creator      *int      `json:"creator"`
	LastModifier *int      `json:"last_modifier"`
}

// IModel 定义接口，用于sqlCondition传参
type IModel interface {
	TableName() string
}
