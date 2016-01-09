package gowork

import (
	"reflect"
)

func SliceContains(slice interface{}, v interface{}) bool {

	if slice == nil {
		return false
	}

	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)
		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == v {
				return true
			}
		}
	}
	return false
}

func StringMapToSlice(m map[string]interface{}) interface{} {

	var r reflect.Value

	i := 0
	for _, v := range m {
		if !r.IsValid() {
			val := reflect.Indirect(reflect.ValueOf(v))
			r = reflect.MakeSlice(reflect.SliceOf(val.Type()), len(m), len(m))
		}
		r.Index(i).Set(reflect.ValueOf(v).Elem())
	}
	return r.Interface()
}
