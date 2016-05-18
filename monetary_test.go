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
	require.Equal(t, "100.99 CAD", amount.StringWithCurrency())
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
	require.Equal(t, "101.00 CAD", r.StringWithCurrency())
}

func TestRound_RoundHalfUp_Down(t *testing.T) {

	//setup
	d := inf.NewDec(100994, 3)
	ma := MonetaryAmount{Dec: d, Currency: currencies.CAD}

	//execute
	r := ma.Round(2)

	//verify
	require.Equal(t, "100.99 CAD", r.StringWithCurrency())
}

func TestAdd_Positive(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1", "CAD")
	two := MonetaryAmountFromStringPanic("2", "CAD")

	//execute
	r := one.Add(two)


	//verify
	require.Equal(t, "3 CAD", r.StringWithCurrency())
}

func TestAdd_Negative(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1", "CAD")
	ntwo := MonetaryAmountFromStringPanic("-2", "CAD")

	//execute
	r := one.Add(ntwo)


	//verify
	require.Equal(t, "-1 CAD", r.StringWithCurrency())
}

func TestAdd_ToNoCurrency(t *testing.T) {

	//setup
	one := Zero
	ntwo := MonetaryAmountFromStringPanic("1", "CAD")

	//execute
	r := one.Add(ntwo)


	//verify
	require.Equal(t, "1.00 CAD", r.StringWithCurrency())
}

func TestAdd_WithNoCurrency(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1", "CAD")
	ntwo := Zero

	//execute
	r := one.Add(ntwo)


	//verify
	require.Equal(t, "1.00 CAD", r.StringWithCurrency())
}

func TestAdd_NoCurrency(t *testing.T) {

	//setup
	one := Zero
	ntwo := Zero

	//execute
	r := one.Add(ntwo)


	//verify
	require.Equal(t, "0.00", r.StringWithCurrency())
}

func TestAdd_MNotZeroNoCurrency(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1", "")
	ntwo := MonetaryAmountFromStringPanic("1", "CAD")

	//execute
	require.Panics(t, func() {
		one.Add(ntwo)
	}, "Currency mismatch should occur")


	//verify
}

func TestAdd_MANotZeroNoCurrency(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1", "CAD")
	ntwo := MonetaryAmountFromStringPanic("1", "")

	//execute
	require.Panics(t, func() {
		one.Add(ntwo)
	}, "Currency mismatch should occur")


	//verify
}

func TestAdd_CurrencyMismatch(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1", "CAD")
	ntwo := MonetaryAmountFromStringPanic("1", "USD")

	//execute
	require.Panics(t, func() {
		one.Add(ntwo)
	}, "Currency mismatch should occur")

	//verify
}

func TestAdd_MMisconfigured(t *testing.T) {

	//setup
	one := MonetaryAmount{Currency: currencies.CAD}
	ntwo := MonetaryAmountFromStringPanic("1", "CAD")

	//execute
	require.Panics(t, func() {
		one.Add(ntwo)
	}, "Currency mismatch should occur")

	//verify
}

func TestAdd_MaMisconfigured(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1", "CAD")
	ntwo := MonetaryAmount{Currency: currencies.CAD}

	//execute
	require.Panics(t, func() {
		one.Add(ntwo)
	}, "Currency mismatch should occur")

	//verify
}

func TestAdd_BothMisconfigured(t *testing.T) {

	//setup
	one := MonetaryAmount{Currency: currencies.CAD}
	ntwo := MonetaryAmount{Currency: currencies.CAD}

	//execute
	require.Panics(t, func() {
		one.Add(ntwo)
	}, "Currency mismatch should occur")

	//verify
}

func TestAdd_NeitherCurrency(t *testing.T) {

	//setup
	one := MonetaryAmount{Dec: inf.NewDec(100, 2)}
	ntwo := MonetaryAmount{Dec: inf.NewDec(101, 2)}

	//execute
	require.Panics(t, func() {
		one.Add(ntwo)
	}, "Currency mismatch should occur")

	//verify
}

func TestNeg_Positive(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1", "CAD")

	//execute
	r := one.Neg()


	//verify
	require.Equal(t, "-1 CAD", r.StringWithCurrency())
}

func TestNeg_Negative(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("-1", "CAD")

	//execute
	r := one.Neg()


	//verify
	require.Equal(t, "1 CAD", r.StringWithCurrency())
}

func TestAssumeScale_2(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1.00", "CAD")

	//execute
	r := one.AssumeScale(2)


	//verify
	require.Equal(t, int64(100), r)
}

func TestAssumeScale_3(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("1.999", "CAD")

	//execute
	r := one.AssumeScale(3)


	//verify
	require.Equal(t, int64(1999), r)
}

func TestAssumeScale_Repeating(t *testing.T) {

	//setup
	one := MonetaryAmountFromStringPanic("3.3333333333333333", "CAD")

	//execute
	r := one.AssumeScale(2)


	//verify
	require.Equal(t, int64(333), r)
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
	require.Equal(t, "2", v.StringWithCurrency())
}
