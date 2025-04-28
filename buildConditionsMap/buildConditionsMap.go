package buildconditionsmap

import (
	"reflect"
	"strings"
)

type ConditionConfig struct {
	TableFields map[string][]string    // 表名 -> 字段名映射
	ExactFields []string               // 精确匹配字段名（如 "name"）
	LikeFields  []string               // 模糊匹配字段名（如 "name"）
	RangeFields []string               // 范围字段名（字段值应为 map[string]interface{}，支持 gt/lt）
	FieldValues map[string]interface{} // 动态传参场景，如前端传 {"name": "Alice"}
}

// 返回结构体字段+字段值组合的通用条件 map
func BuildConditionsMapWithConfig(input interface{}, config ConditionConfig) map[string]interface{} {
	conditionMap := make(map[string]interface{})
	exact := make(map[string]interface{})
	like := make(map[string]interface{})
	gt := make(map[string]interface{})
	lt := make(map[string]interface{})

	// 处理结构体字段（如 struct 查询）
	if input != nil {
		val := reflect.ValueOf(input)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		for _, field := range config.ExactFields {
			tableField := getMappedField(field, config.TableFields)
			if v := getFieldByPath(val, field); v.IsValid() && !isEmpty(v) {
				exact[tableField] = v.Interface()
			}
		}
		for _, field := range config.LikeFields {
			tableField := getMappedField(field, config.TableFields)
			if v := getFieldByPath(val, field); v.IsValid() && !isEmpty(v) {
				like[tableField] = v.Interface()
			}
		}
		for _, field := range config.RangeFields {
			v := getFieldByPath(val, field)
			if v.IsValid() && v.Kind() == reflect.Map {
				m := v.Interface().(map[string]interface{})
				tableField := getMappedField(field, config.TableFields)
				if gtVal, ok := m["gt"]; ok {
					gt[tableField] = gtVal
				}
				if ltVal, ok := m["lt"]; ok {
					lt[tableField] = ltVal
				}
			}
		}
	}

	// 处理动态 map 输入（如前端 JSON）
	for field, value := range config.FieldValues {
		tableField := getMappedField(field, config.TableFields)
		if contains(config.ExactFields, field) {
			exact[tableField] = value
		} else if contains(config.LikeFields, field) {
			like[tableField] = value
		}
	}

	// RangeFields 中如果来自 FieldValues 的支持
	for _, field := range config.RangeFields {
		if v, ok := config.FieldValues[field]; ok {
			if m, ok := v.(map[string]interface{}); ok {
				tableField := getMappedField(field, config.TableFields)
				if gtVal, ok := m["gt"]; ok {
					gt[tableField] = gtVal
				}
				if ltVal, ok := m["lt"]; ok {
					lt[tableField] = ltVal
				}
			}
		}
	}

	if len(exact) > 0 {
		conditionMap["exact"] = exact
	}
	if len(like) > 0 {
		conditionMap["like"] = like
	}
	if len(gt) > 0 {
		conditionMap["gt"] = gt
	}
	if len(lt) > 0 {
		conditionMap["lt"] = lt
	}
	return conditionMap
}

// 获取字段对应的完整表字段（例如 "users.name"）
func getMappedField(field string, tableFields map[string][]string) string {
	for table, fields := range tableFields {
		for _, f := range fields {
			if f == field {
				return table + "." + f
			}
		}
	}
	return field
}

// 获取结构体字段值
func getFieldByPath(v reflect.Value, path string) reflect.Value {
	parts := strings.Split(path, ".")
	for _, part := range parts {
		v = v.FieldByName(part)
		if !v.IsValid() {
			return reflect.Value{}
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}
	return v
}

// 判断值是否为空
func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}

func contains(list []string, key string) bool {
	for _, s := range list {
		if s == key {
			return true
		}
	}
	return false
}
