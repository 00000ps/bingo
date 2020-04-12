package utils

import (
	"fmt"
	"reflect"
	"sort"
)

// Abs return the absolute value
func Abs(v interface{}) interface{} {
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

// Max returns the maximum number of slice or array
func Max(v ...interface{}) interface{} { return sorts(true, v...) }

// Min returns the minimum number of slice or array
func Min(v ...interface{}) interface{} { return sorts(false, v...) }
func sorts(max bool, v ...interface{}) interface{} {
	l := len(v)
	if l == 0 {
		return nil
	} else if l == 1 {
		return v[0]
	}

	var s Slice
	for _, vi := range v {
		s = append(s, vi)
	}

	if !sort.IsSorted(s) {
		sort.Sort(s)
	}
	if max {
		return s[l-1]
	}
	return s[0]
}

// Slice is an interface that can be sorting
// note that all elements should be numerous
type Slice []interface{}

func (c Slice) Len() int           { return len(c) }
func (c Slice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Slice) Less(i, j int) bool { return ConvertToFloat64(c[i]) < ConvertToFloat64(c[j]) }

// ConvertToFloat64 used to convert all kinds of number to float64
func ConvertToFloat64(v interface{}) float64 {
	switch vs := reflect.ValueOf(v); vs.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(vs.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(vs.Uint())
	case reflect.Float32, reflect.Float64:
		return float64(vs.Float())
	default:
		panic(fmt.Errorf("unsupported type %s to compare", vs.Kind().String()))
	}
}
