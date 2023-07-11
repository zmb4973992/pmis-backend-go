package model

func IntToPointer(number int) *int {
	return &number
}

func Int64ToPointer(number int64) *int64 {
	return &number
}

func BoolToPointer(param bool) *bool {
	return &param
}

func Float64ToPointer(param float64) *float64 {
	return &param
}

func StringToPointer(param string) *string {
	return &param
}
