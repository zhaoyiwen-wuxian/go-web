package common

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type ConditionConfig struct {
	ExactFields []string
	LikeFields  []string
	RangeFields []string
}

func BuildConditionsMapWithConfig(input interface{}, config ConditionConfig) map[string]interface{} {
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	conditionMap := make(map[string]interface{})
	exact := make(map[string]interface{})
	like := make(map[string]interface{})
	gt := make(map[string]interface{})
	lt := make(map[string]interface{})

	for _, field := range config.ExactFields {
		if v := getFieldByPath(val, field); v.IsValid() && !isEmpty(v) {
			exact[field] = v.Interface()
		}
	}
	for _, field := range config.LikeFields {
		if v := getFieldByPath(val, field); v.IsValid() && !isEmpty(v) {
			like[field] = v.Interface()
		}
	}
	for _, field := range config.RangeFields {
		v := getFieldByPath(val, field)
		if v.IsValid() && v.Kind() == reflect.Map {
			m := v.Interface().(map[string]interface{})
			if gtVal, ok := m["gt"]; ok {
				gt[field] = gtVal
			}
			if ltVal, ok := m["lt"]; ok {
				lt[field] = ltVal
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

func BuildQueryWithConditions(db *gorm.DB, conditions map[string]interface{}) *gorm.DB {
	exact := map[string]interface{}{}
	gt := map[string]interface{}{}
	lt := map[string]interface{}{}
	like := map[string]interface{}{}

	for k, v := range conditions {
		switch k {
		case "exact":
			exact = v.(map[string]interface{})
		case "gt":
			gt = v.(map[string]interface{})
		case "lt":
			lt = v.(map[string]interface{})
		case "like":
			like = v.(map[string]interface{})
		default:
			exact[k] = v
		}
	}

	return BuildConditions(db, exact, gt, lt, like)
}
