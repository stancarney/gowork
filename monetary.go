package gowork

import (
	"fmt"
	inf "gopkg.in/inf.v0"
	"strings"
	"github.com/stancarney/gowork/currencies"
)

var decZero = inf.NewDec(0, 2)
var maZero = MonetaryAmount{Dec: inf.NewDec(0, 2)}

func Zero(currency currencies.Currency) MonetaryAmount {
	return MonetaryAmount{Dec: decZero, Currency: currency}
}

//Minimal hiding of the underlying Decimal implementation so API changes or a when better implementation comes along it won't cause huge app changes.
type MonetaryAmount struct {
	*inf.Dec
	Currency currencies.Currency
}

func (m MonetaryAmount) Add(ma MonetaryAmount) MonetaryAmount {
	return MonetaryAmount{Dec: new(inf.Dec).Add(m.Dec, ma.Dec), Currency: m.Currency}
}

func (m MonetaryAmount) Multiply(i *inf.Dec) MonetaryAmount {
	return MonetaryAmount{Dec: new(inf.Dec).Mul(m.Dec, i), Currency: m.Currency}
}

func (m MonetaryAmount) Neg() MonetaryAmount {
	return MonetaryAmount{Dec: new(inf.Dec).Neg(m.Dec), Currency: m.Currency}
}

func (m MonetaryAmount) Abs() MonetaryAmount {
	return MonetaryAmount{Dec: new(inf.Dec).Abs(m.Dec), Currency: m.Currency}
}

func (m MonetaryAmount) Round(scale int) string { //TODO:Stan change to return MA?
	return MonetaryAmount{Dec: new(inf.Dec).Round(m.Dec, inf.Scale(scale), inf.RoundHalfUp), Currency: m.Currency}.String()
}

func (m MonetaryAmount) AssumeScale(scale int) string {  //TODO:Stan change to return MA?
	return strings.Replace(m.Round(scale), ".", "", 1)
}

func (m MonetaryAmount) Cmp(ma MonetaryAmount) int {
	return m.Dec.Cmp(ma.Dec)
}

func (m MonetaryAmount) IsZero() bool {
	return m.Cmp(maZero) == 0
}

//UnmarshalText mutates the current MonetaryAmount instance with the provided byte array. 
func (m *MonetaryAmount) UnmarshalText(data []byte) error {
	m.Dec = new(inf.Dec)
	return m.Dec.UnmarshalText(data)
}

func (m MonetaryAmount) String() string {
	return fmt.Sprintf("%s%4s", m.Dec.String(), m.Currency)
}

func MonetaryAmountFromString(value string, cur string) (amount MonetaryAmount, err error) {
	dec, ok := new(inf.Dec).SetString(value)
	if !ok {
		return MonetaryAmount{}, fmt.Errorf("Invalid string input: %s %s", value, cur)
	}

	return MonetaryAmount{Dec: dec, Currency: currencies.Currency(cur)}, nil
}

// MonetaryAmountFromStringPanic is just like MonetaryAmountFromString only it will panic if the string is invalid.
// Useful for things like unit tests where the amount is hardcoded in the test.
func MonetaryAmountFromStringPanic(value string, cur string) MonetaryAmount {
	m, err := MonetaryAmountFromString(value, cur)
	if err != nil {
		panic(err)
	}
	return m
}