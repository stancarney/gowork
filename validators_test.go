package gowork

import (
	"github.com/stretchr/testify/require"
	"testing"
	"github.com/stancarney/gowork/currencies"
)

func TestIsNotZeroValue_IntZero(t *testing.T) {

	//setup
	v := 0

	//execute
	r := IsNotZeroValue(v, "")

	//verify
	require.Error(t, r)
}

func TestIsNotZeroValue_IntNotZero(t *testing.T) {

	//setup
	v := 1

	//execute
	r := IsNotZeroValue(v, "")

	//verify
	require.Nil(t, r)
}

func TestIsNotZeroValue_MonetaryAmountZero(t *testing.T) {

	//setup
	v := MonetaryAmount{}

	//execute
	r := IsNotZeroValue(v, "")

	//verify
	require.Error(t, r)
}

func TestIsNotZeroValue_MonetaryAmountNotZero(t *testing.T) {

	//setup
	v := NewMonetaryAmount(currencies.CAD) // Zero as in $0.00 not as in doesn't exist.

	//execute
	r := IsNotZeroValue(v, "")

	//verify
	require.Nil(t, r)
}

func TestIsNotZeroValue_ErrorResponseZero(t *testing.T) {

	//setup
	v := ErrorResponse{}

	//execute
	r := IsNotZeroValue(v, "")

	//verify
	require.Error(t, r)
}

func TestIsNotZeroValue_ErrorResponsePartialZero(t *testing.T) {

	//setup
	v := ErrorResponse{}
	v.Count = 1

	//execute
	r := IsNotZeroValue(v, "")

	//verify
	require.Nil(t, r)
}