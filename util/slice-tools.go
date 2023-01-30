package util

import "reflect"

// byte 是 uint8 的别名,rune 是 int32 的别名
type typeForSliceComparing interface {
	bool | string | int | int64 | int32 | int16 | int8 |
		uint | uint64 | uint32 | uint16 | uint8 |
		float64 | float32
}

// 校验单个内容是否包含在切片中

func IsInSlice[T typeForSliceComparing](element T, slice []T) bool {
	for _, v := range slice {
		if element == v {
			return true
		}
	}
	return false
}

// 校验切片是否包含单个内容

func SliceIncludes(slice any, element any) bool {
	tempSlice := reflect.ValueOf(slice)
	if tempSlice.Kind() == reflect.Slice {
		for i := 0; i < tempSlice.Len(); i++ {
			if tempSlice.Index(i).Interface() == element {
				return true
			}
		}
	}
	return false
}

// deprecated
// 建议用新方法SliceIncludes
func SliceIncludesOld[T typeForSliceComparing](slice []T, element T) bool {
	for _, value := range slice {
		if element == value {
			return true
		}
	}
	return false
}

// 校验两个切片的值是否相等（不看顺序）

func SlicesAreSame[T typeForSliceComparing](slice1 []T, slice2 []T) bool {
	//如果任意切片的长度为0
	if len(slice1) == 0 || len(slice2) == 0 {
		return false
	}
	//如果两个切片的长度不相等
	if len(slice1) != len(slice2) {
		return false
	}
	//对两个切片进行双重遍历比较
	for _, element := range slice1 {
		res := IsInSlice(element, slice2)
		if res == false {
			return false
		}
	}
	return true
}
