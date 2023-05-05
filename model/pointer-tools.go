package model

func IntToPointer(number int) *int {
	return &number
}

func BoolToPointer(param bool) *bool {
	return &param
}
