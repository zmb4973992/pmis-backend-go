package util

import (
	"errors"
	"reflect"
)

// byte 是 uint8 的别名,rune 是 int32 的别名
type typeForSliceTools interface {
	bool | string | int | int64 | int32 | int16 | int8 |
		uint | uint64 | uint32 | uint16 | uint8 |
		float64 | float32
}

// 校验单个内容是否包含在切片中

func IsInSlice[T typeForSliceTools](element T, array []T) bool {
	for _, v := range array {
		if element == v {
			return true
		}
	}
	return false
}

// 校验切片是否包含单个内容

func SliceIncludes(array any, element any) bool {
	tempSlice := reflect.ValueOf(array)
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
func SliceIncludesOld[T typeForSliceTools](array []T, element T) bool {
	for _, value := range array {
		if element == value {
			return true
		}
	}
	return false
}

// 校验两个切片的值是否相等（不看顺序）

func SlicesAreSame[T typeForSliceTools](array1 []T, array2 []T) bool {
	//如果任意切片的长度为0
	if len(array1) == 0 || len(array2) == 0 {
		return false
	}
	//如果两个切片的长度不相等
	if len(array1) != len(array2) {
		return false
	}
	//对两个切片进行双重遍历比较
	for _, element := range array1 {
		res := IsInSlice(element, array2)
		if res == false {
			return false
		}
	}
	return true
}

// RemoveDuplication 数组去重
func RemoveDuplication[T typeForSliceTools](array []T) []T {
	var result []T
	for i := range array {
		if len(result) == 0 {
			result = append(result, array[i])
		} else {
			for j := range result {
				if array[i] == result[j] {
					break
				}
				if j == len(result)-1 {
					result = append(result, array[i])
				}
			}
		}
	}
	return result
}

// FindMax 传入int数组，给出最大值、最大值的位置和错误信息（如果数组为空则报错）
func FindMax(array []int) (maxValue int, maxIndex int, err error) {
	if len(array) == 0 {
		return 0, 0, errors.New("错误：数组为空")
	}
	maxValue = array[0]
	maxIndex = 0
	for i := 1; i < len(array); i++ {
		if maxValue < array[i] {
			maxValue = array[i]
			maxIndex = i
		}
	}
	return maxValue, maxIndex, nil
}

// 双指针翻转切片
// 入参：[]string{"a","b","c"}   结果：[]string{"c","b","a"}
func reverseSlice[T typeForSliceTools](param []T) []T {
	left, right := 0, len(param)-1
	for left < right {
		param[left], param[right] = param[right], param[left]
		left++
		right--
	}
	return param
}
