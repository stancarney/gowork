package gowork

import (
	"reflect"
	"fmt"
	"errors"
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

// ChopSlice returns a map of the input slice grouped by the specified field name.
func ChopSlice(slice interface{}, field string) (interface{}, error) {
	if slice == nil {
		return nil, errors.New("Slice is nil")
	}

	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)
		if s.Len() == 0 {
			return nil, errors.New("Slice is empty")
		}

		first := s.Index(0)
		m := reflect.MakeMap(reflect.MapOf(reflect.Indirect(first).FieldByName(field).Type(), reflect.SliceOf(first.Type())))

		for i := 0; i < s.Len(); i++ {
			key := reflect.Indirect(s.Index(i)).FieldByName(field)
			value := m.MapIndex(key)
			if !value.IsValid() {
				value = reflect.MakeSlice(reflect.SliceOf(first.Type()), 0, 1) //TODO:Stan make a better guess of length
			}
			value = reflect.Append(value, s.Index(i))
			m.SetMapIndex(key, value)
		}

		return m.Interface(), nil
	default:
		return nil, errors.New("Argument is not a slice")
	}

	return nil, errors.New("Something bad happened")
}

// ChopSortedSlice is similar to ChopSlice except rather than returning new slices for the map values it just re-slices the input slice, as a result the input slice must be sorted.
// If the input slice is not sorted the results are undefined.
func ChopSortedSlice(slice interface{}, field string) (interface{}, error) {

	if slice == nil {
		return nil, errors.New("Slice is nil")
	}

	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)
		if s.Len() == 0 {
			return nil, errors.New("Slice is empty")
		}

		first := s.Index(0)
		m := reflect.MakeMap(reflect.MapOf(reflect.Indirect(first).FieldByName(field).Type(), reflect.SliceOf(first.Type())))

		start := 0
		for i := 0; i < s.Len(); i++ {
			key := reflect.Indirect(s.Index(i)).FieldByName(field)
			value := m.MapIndex(key)
			if !value.IsValid() {
				start = i
			}
			value = s.Slice(start, i + 1)
			m.SetMapIndex(key, value)
		}

		return m.Interface(), nil
	default:
		return nil, errors.New("Argument is not a slice")
	}

	return nil, errors.New("Something bad happened")
}
