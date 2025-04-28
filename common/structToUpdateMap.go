package common

import (
	"reflect"
)

// StructToMap 将结构体转换为 map[string]interface{}
func StructToMap(obj interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)

	// 如果是指针类型，需要获取指向的值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 获取结构体的字段和对应的值
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// 忽略零值字段（如果需要的话）
		if fieldValue.IsZero() {
			continue
		}

		// 将字段名转换为小写并赋值
		fieldName := lowerFirstChar(field.Name)
		result[fieldName] = fieldValue.Interface()
	}

	return result, nil
}

// lowerFirstChar 将字段首字母转为小写
func lowerFirstChar(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = rune(runes[0] + 32) // 将首字母转为小写
	return string(runes)
}
