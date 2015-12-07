package gowork

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(i interface{}) string {
	if i == nil {
		return ""
	}
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	segments := strings.Split(name, "/")
	return segments[len(segments) - 1]
}

func GetStructName(i interface{}) string {
	if i == nil {
		return ""
	}
	return reflect.Indirect(reflect.ValueOf(i)).Type().String()
}

func GetStringValue(i interface{}, fieldName string) string {
	if i == nil {
		return ""
	}

	if reflect.TypeOf(i).Kind() == reflect.String {
		return i.(string)
	}

	return reflect.Indirect(reflect.ValueOf(i)).FieldByName(fieldName).Interface().(string)
}
