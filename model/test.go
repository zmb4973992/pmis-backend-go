package model

import "time"

type Test struct {
	Date1 time.Time `gorm:"type:datetime"`
	Date2 time.Time `gorm:"type:datetime"`
	Date3 time.Time `gorm:"type:datetime;scale:0"`
}

func (t *Test) TableName() string {
	return "test"
}
