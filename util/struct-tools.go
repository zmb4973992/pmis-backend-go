package util

import (
	"errors"
	"reflect"
)

func StructToMap(structData any) (mapData map[string]any, err error) {
	t := reflect.TypeOf(structData)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("入参不是结构体，请检查")
	}
	v := reflect.ValueOf(structData)
	var data = make(map[string]any)
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data, nil
}
