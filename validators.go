package gowork

import (
	"errors"
	"reflect"
)

func IsNotZeroValue(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if !st.IsValid() {
		return errors.New(CANNOT_BE_ZERO_VALUE)
	}
	return nil
}