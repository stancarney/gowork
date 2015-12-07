package gowork

import (
	"reflect"
	"log"
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
	default:
		log.Println("Not a slice:", reflect.TypeOf(slice).Kind())
	}
	return false
}

