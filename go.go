package gowork

import (
	"reflect"
	"runtime"
	"strings"
)

func FunctionName(i interface{}) string {
	if i == nil {
		return ""
	}
	name := runtime.FuncForPC(reflect.Indirect(reflect.ValueOf(i)).Pointer()).Name()
	segments := strings.Split(name, "/")
	return segments[len(segments) - 1]
}

func CurrentFunctionName(stack int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(stack, pc)
	name := runtime.FuncForPC(pc[0]).Name()
	segments := strings.Split(name, "/")
	return segments[len(segments) - 1]
}

func StructName(i interface{}) string {
	if i == nil {
		return ""
	}
	return reflect.Indirect(reflect.ValueOf(i)).Type().String()
}

func StringValue(i interface{}, fieldName string) string {
	if i == nil {
		return ""
	}

	if s, ok := i.(string); ok {
		return s
	}

	return reflect.Indirect(reflect.ValueOf(i)).FieldByName(fieldName).Interface().(string)
}
