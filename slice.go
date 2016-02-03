package gowork

import (
	"reflect"
	"fmt"
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
		val := reflect.Indirect(reflect.ValueOf(v))
		if !r.IsValid() {
			r = reflect.MakeSlice(reflect.SliceOf(val.Type()), len(m), len(m))
		}
		r.Index(i).Set(val)
		i++
	}

	if !r.IsValid() {
		return nil
	}

	return r.Interface()
}

// DumpSlice is used to dump the contents of a slice out to std out and is useful in testing.
func DumpSlice(slice interface{}) {
	t := reflect.TypeOf(slice)
	switch t.Kind() {
	case reflect.Slice:
		fmt.Printf("\nType: %s\n", t)
		s := reflect.ValueOf(slice)
		for i := 0; i < s.Len(); i++ {
			fmt.Printf("%d: %s\n", i, s.Index(i))
		}
	default:
		panic(fmt.Sprintf("not a slice, %s", t))
	}
}
