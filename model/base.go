package model

import "time"

type BaseModel struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IModel 定义接口，用于sqlCondition传参
type IModel interface {
	TableName() string
}
