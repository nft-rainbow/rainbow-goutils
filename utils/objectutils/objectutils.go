package objectutils

import (
	"encoding/json"
	"reflect"
)

func AllFieldsNotEmpty(obj interface{}) bool {
	// 获取对象的值
	v := reflect.ValueOf(obj)

	// 确保传入的是结构体类型
	if v.Kind() != reflect.Struct {
		return false
	}

	// 遍历结构体的所有字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// 检查字段类型是否为零值
		if field.IsZero() {
			return false
		}
	}
	return true
}

func AllFieldsEmpty(obj interface{}) bool {
	// 获取对象的值
	v := reflect.ValueOf(obj)

	// 确保传入的是结构体类型
	if v.Kind() != reflect.Struct {
		return false
	}

	// 遍历结构体的所有字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// 检查字段类型是否非零值
		if !field.IsZero() {
			return false
		}
	}
	return true
}

func ConvertToMap(obj interface{}) (map[string]interface{}, error) {
	// Marshal the object to JSON
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a map
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
