package gowork
import (
	"errors"
	"log"
	"reflect"
)

func IsNotZeroValue(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	log.Println("st.IsValid()", st.IsValid())
	if !st.IsValid() {
		return errors.New(CANNOT_BE_ZERO_VALUE)
	}
	return nil
}