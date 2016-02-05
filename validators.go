package gowork

import (
	"errors"
	"reflect"
)

// IsNotZeroValueDeep is a function used by the validator.V2 GO lib. It should only be used on Structs for performance reasons. 
// Builtins will perform much better using the 'nonzero' validator provided by the lib.
func IsNotZeroValue(v interface{}, param string) error {
	if reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()) {
		return errors.New(CANNOT_BE_ZERO_VALUE)
	}
	return nil
}