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
	name := runtime.FuncForPC(reflect.Indirect(reflect.ValueOf(i)).Pointer()).Name()
	segments := strings.Split(name, "/")
	return segments[len(segments) - 1]
}

func GetCurrentFunctionName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	name := runtime.FuncForPC(pc[0]).Name()
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
