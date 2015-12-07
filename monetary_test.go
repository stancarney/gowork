package gowork

import (
	"github.com/stretchr/testify/require"
	inf "gopkg.in/inf.v0"
	"testing"
)

func TestMonetaryAmountFromString_Valid(t *testing.T) {

	//setup

	//execute
	amount, err := MonetaryAmountFromString("100.99")

	//verify
	require.Nil(t, err)
	require.Equal(t, "100.99", amount.String())
}

func TestMonetaryAmountFromString_Invalid(t *testing.T) {

	//setup

	//execute
	_, err := MonetaryAmountFromString("aa")

	//verify
	require.NotNil(t, err)
	require.Equal(t, "Invalid string input: aa", err.Error())
}

func TestRound_RoundHalfUp_Up(t *testing.T) {

	//setup
	d := inf.NewDec(100995, 3)
	ma := MonetaryAmount{d}

	//execute
	r := ma.Round(2)

	//verify
	require.Equal(t, "101.00", r)
}

func TestRound_RoundHalfUp_Down(t *testing.T) {

	//setup
	d := inf.NewDec(100994, 3)
	ma := MonetaryAmount{d}

	//execute
	r := ma.Round(2)

	//verify
	require.Equal(t, "100.99", r)
}

func TestAdd_Positive(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1")
	two, _ := MonetaryAmountFromString("2")

	//execute
	r := one.Add(two)


	//verify
	require.Equal(t, "3", r.String())
}

func TestAdd_Negative(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1")
	ntwo, _ := MonetaryAmountFromString("-2")

	//execute
	r := one.Add(ntwo)


	//verify
	require.Equal(t, "-1", r.String())
}

func TestNeg_Positive(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1")

	//execute
	r := one.Neg()


	//verify
	require.Equal(t, "-1", r.String())
}

func TestNeg_Negative(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("-1")

	//execute
	r := one.Neg()


	//verify
	require.Equal(t, "1", r.String())
}

func TestAssumeScale(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1.00")

	//execute
	r := one.AssumeScale(2)


	//verify
	require.Equal(t, "100", r)
}

func TestMarshallText(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1")

	//execute
	r, err := one.MarshalText()


	//verify
	require.Nil(t, err)
	require.Equal(t, "1", string(r))
}

func TestUnMarshallText(t *testing.T) {

	//setup
	var v MonetaryAmount
	data := []byte("2")

	//execute
	err := v.UnmarshalText(data)


	//verify
	require.Nil(t, err)
	require.Equal(t, "2", v.String())
}
