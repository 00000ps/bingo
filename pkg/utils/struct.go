package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"unicode"
)

var NilParas = make([]reflect.Value, 0)

func Func(i interface{}) *runtime.Func {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer())
}
func GetFunc(deepth int) *runtime.Func {
	pc, _, _, ok := runtime.Caller(deepth)
	if !ok {
		return nil
	}
	return runtime.FuncForPC(pc)
}
func CallMethod(value reflect.Value, name string, in ...reflect.Value) []reflect.Value {
	method := value.MethodByName(name)
	if len(in) == 0 {
		in = NilParas
	}
	if method.IsValid() {
		return method.Call(in)
	}
	return NilParas
}
func CallMethodGetOneReturn(value reflect.Value, name string, in ...reflect.Value) interface{} {
	rets := CallMethod(value, name, in...)
	if len(rets) > 0 {
		return rets[0].Interface()
	}
	return nil
}
func IsValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}
func GetItemsOfStruct(v reflect.Value) (r []string) {
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			r = append(r, v.Type().Field(i).Name)
			//fmt.Printf("000000=%#v\n", v.Type().Field(i).Name)
			//fmt.Printf("111111=%#v\n", r)
			if v.Field(i).Kind() == reflect.Slice {
				slice := reflect.MakeSlice(v.Field(i).Type(), 1, 1)
				if slice.Index(0).Kind() == reflect.Struct {
					r = append(r, GetItemsOfStruct(slice.Index(0))...)
				}
			} else if v.Field(i).Kind() == reflect.Struct {
				r = append(r, GetItemsOfStruct(v.Field(i))...)
			}
		}
	}
	return
}
func GetKeyInItems(keys []string, key string) string {
	for i := 0; i < len(keys); i++ {
		//fmt.Printf("%#v ----- %#v\n", key, keys[i])
		if strings.ToLower(key) == strings.ToLower(keys[i]) {
			return keys[i]
		}
	}
	fmt.Printf("not found: %#v in %#v (lib.GetKeyInItems)\n", key, keys)
	return key
}
func GetItemWithTag(v reflect.Value, key string) (name string, rv reflect.Value) {
	/*find := false
	for i := 0; i < v.NumField(); i++ {
		sf := v.Type().Field(i)
		tag := sf.Tag.Get("json")
		name = sf.Name

		b := []byte(key)
		if bytes.HasSuffix(b, []byte{0x0}) {
			key = string(b[:len(b)-1])
		}

		if tag != "" && tag == key {
			rv = v.Field(i)
			find = true
			break
		}

		if strings.EqualFold(name, key) {
			rv = v.Field(i)
			find = true
			break
		} else if key == "" {
			rv = v.Field(i)
			find = true
			break
		}
	}
	if !find {
		//fmt.Printf("%#v!!!~~~~~%#v!!!\n", key, name)
		return "", rv
	}
	return name, rv*/
	//fmt.Printf("%#v!!!~~~~~%#v!!!\n", key, v)
	for i := 0; i < v.NumField(); i++ {
		sf := v.Type().Field(i)
		tag := sf.Tag.Get("json")
		name = sf.Name
		//fmt.Printf("%#v!!!~~~~~%#v!!!\n", key, v.Field(i))

		b := []byte(key)
		if bytes.HasSuffix(b, []byte{0x0}) {
			key = string(b[:len(b)-1])
		}

		if tag != "" && tag == key {
			rv = v.Field(i)
			return
		}

		if strings.EqualFold(name, key) {
			rv = v.Field(i)
			//fmt.Printf("%#v!!!~~~~~%#v!!!\n", key, rv.NumField())
			return
		} else if key == "" {
			rv = v.Field(i)
			return
		}
	}
	//fmt.Printf("%#v!!!~~~~~%#v!!!\n", key, name)
	name = ""
	return
}
func ReflectSetValue(v reflect.Value, key string, val interface{}) error {
	if !v.CanSet() {
		// fmt.Printf("cannot set reflect.Value: key=%s, val=%#v set=%#v\n", key, val, v)
		return fmt.Errorf("cannot set reflect.Value: key=%s, val=%#v set=%#v\n", key, val, v)
	}
	if key == "" {
		v.Set(reflect.ValueOf(val))
		return nil
	}
	name, rv := GetItemWithTag(v, key)
	if name == "" {
		// log.Debug("(check pipe interface) failed to get item of %s in ReflectSetValue", key)
		return nil
		//return fmt.Errorf("(check pipe interface) failed to get item of %s in ReflectSetValue", key)
	}

	ReflectSetItem(rv, key, val)
	return nil
}
func ReflectSetItem(v reflect.Value, key string, val interface{}) error {
	if key == "" {
		// fmt.Printf("set item failed, key is blank")
		return fmt.Errorf("set item failed, key is blank")
	}

	if v.IsValid() && v.CanSet() {
		//fmt.Printf("set Succ: %s=%#v (IsValid=%#v; CanSet=%#v)\n", key, val, v.IsValid(), v.CanSet())
		v.Set(reflect.ValueOf(val))
	} else {
		// fmt.Printf("fail: %s=%#v (IsValid=%#v; CanSet=%#v)\n", key, val, v.IsValid(), v.CanSet())
		return fmt.Errorf("set Fail: %s=%#v (IsValid=%#v; CanSet=%#v)", key, val, v.IsValid(), v.CanSet())
	}
	return nil
}
func SetStruct(v interface{}, key string, value interface{}) error {
	pv := reflect.ValueOf(v)
	if pv.Kind() != reflect.Ptr || pv.IsNil() {
		return fmt.Errorf("unmarshal need Ptr data")
	}
	elem := pv.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("unmarshal need Struct data")
	}
	return ReflectSetValue(elem, key, value)
}

type mapData struct {
	data map[interface{}]interface{}
	lock sync.RWMutex
}

func (d *mapData) Has(k interface{}) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	_, ok := d.data[k]
	return ok
}
func (d *mapData) Get(k interface{}) (interface{}, bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	v, ok := d.data[k]
	return v, ok
}

func (d *mapData) Set(k interface{}, c interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.data[k] = c
}
func (d *mapData) Bind(md map[interface{}]interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.data = md
}
