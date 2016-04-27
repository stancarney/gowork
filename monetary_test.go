package gowork

import (
	"github.com/stretchr/testify/require"
	inf "gopkg.in/inf.v0"
	"testing"
	"github.com/stancarney/gowork/currencies"
)

func TestMonetaryAmountFromString_Valid(t *testing.T) {

	//setup

	//execute
	amount, err := MonetaryAmountFromString("100.99", "CAD")

	//verify
	require.Nil(t, err)
	require.Equal(t, "100.99 CAD", amount.String())
}

func TestMonetaryAmountFromString_InvalidAmount(t *testing.T) {

	//setup

	//execute
	_, err := MonetaryAmountFromString("aa", "CAD")

	//verify
	require.NotNil(t, err)
	require.Equal(t, "Invalid string input: aa CAD", err.Error())
}

func TestRound_RoundHalfUp_Up(t *testing.T) {

	//setup
	d := inf.NewDec(100995, 3)
	ma := MonetaryAmount{Dec: d, Currency: currencies.CAD}

	//execute
	r := ma.Round(2)

	//verify
	require.Equal(t, "101.00 CAD", r)
}

func TestRound_RoundHalfUp_Down(t *testing.T) {

	//setup
	d := inf.NewDec(100994, 3)
	ma := MonetaryAmount{Dec: d, Currency: currencies.CAD}

	//execute
	r := ma.Round(2)

	//verify
	require.Equal(t, "100.99 CAD", r)
}

func TestAdd_Positive(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1", "CAD")
	two, _ := MonetaryAmountFromString("2", "CAD")

	//execute
	r := one.Add(two)


	//verify
	require.Equal(t, "3 CAD", r.String())
}

func TestAdd_Negative(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1", "CAD")
	ntwo, _ := MonetaryAmountFromString("-2", "CAD")

	//execute
	r := one.Add(ntwo)


	//verify
	require.Equal(t, "-1 CAD", r.String())
}

func TestNeg_Positive(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1", "CAD")

	//execute
	r := one.Neg()


	//verify
	require.Equal(t, "-1 CAD", r.String())
}

func TestNeg_Negative(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("-1", "CAD")

	//execute
	r := one.Neg()


	//verify
	require.Equal(t, "1 CAD", r.String())
}

func TestAssumeScale(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1.00", "CAD")

	//execute
	r := one.AssumeScale(2)


	//verify
	require.Equal(t, "100 CAD", r)
}

func TestMarshallText(t *testing.T) {

	//setup
	one, _ := MonetaryAmountFromString("1", "CAD")

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
	require.Equal(t, "2    ", v.String()) //TODO:Stan this is the result of how String is currently working. Needs to be resolved.
}
