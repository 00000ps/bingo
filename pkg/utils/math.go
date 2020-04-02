package utils

import (
	"reflect"
)

// ABS return the abs value
func ABS(v interface{}) interface{} {
	switch v.(type) {
	case int:
		if i := v.(int); i < 0 {
			return i * -1
		}
	case int8:
		if i := v.(int8); i < 0 {
			return i * -1
		}
	case int16:
		if i := v.(int16); i < 0 {
			return i * -1
		}
	case int32:
		if i := v.(int32); i < 0 {
			return i * -1
		}
		return v
	case int64:
		if i := v.(int64); i < 0 {
			return i * -1
		}
	case float32:
		if i := v.(float32); i < 0 {
			return i * -1
		}
	case float64:
		if i := v.(float64); i < 0 {
			return i * -1
		}
	}

	return v
}

func Abs(v interface{}) interface{} {
	switch rv := reflect.ValueOf(v); v.(type) {
	case int, int8, int16, int32, int64:
		if i := rv.Int(); i < 0 {
			return i * -1
		}
	case float32, float64:
		if i := rv.Float(); i < 0 {
			return i * -1
		}
	}

	return v
}
