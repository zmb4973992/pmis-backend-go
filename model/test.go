package model

type Test struct {
	BasicModel
	A int64
}

func (t *Test) TableName() string {
	return "test"
}
