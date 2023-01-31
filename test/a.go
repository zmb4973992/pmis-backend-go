package main

import (
	"errors"
	"fmt"
)

func main() {
	a, b, err := FindMax([]int{1, 1, 1, 11, 12, 2})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a, b)
}

func FindMax(array []int) (maxValue int, maxIndex int, err error) {
	if len(array) == 0 {
		return 0, 0, errors.New("数组为空")
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
