package gowork

import (
	"fmt"
	inf "gopkg.in/inf.v0"
	"github.com/stancarney/gowork/currencies"
)

var decZero = inf.NewDec(0, 2)
var Zero = MonetaryAmount{Dec: decZero}

//Minimal hiding of the underlying Decimal implementation so API changes or a when better implementation comes along it won't cause huge app changes.
type MonetaryAmount struct {
	*inf.Dec
	Currency currencies.Currency
}

// Add will add the provided MonetaryAmount to the current MonetaryAmount and return a new MonetaryAmount with the summed value.
// It will panic if the currencies differ.
func (m MonetaryAmount) Add(ma MonetaryAmount) MonetaryAmount {
	switch {
	case m.Currency != "" && m.Currency == ma.Currency: //Happy Case.
	case ma.Currency == "" && ma.Dec == decZero && m.Currency != "": //m.Currency is good so we will just use that one.
	case m.Currency == "" && m.Dec == decZero && ma.Currency != "": //m.Currency is "" so we need to set it.
		m.Currency = ma.Currency
	case m.Dec == decZero && ma.Dec == decZero: //both sides are uninitialized but that is ok
		return Zero
	default:
		panic(fmt.Sprintf("Currency mismatch: %s != %s", m.Currency, ma.Currency))
	}

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

func (m MonetaryAmount) Round(scale int) MonetaryAmount {
	return MonetaryAmount{Dec: new(inf.Dec).Round(m.Dec, inf.Scale(scale), inf.RoundHalfUp), Currency: m.Currency}
}

func (m MonetaryAmount) AssumeScale(scale int) int64 {
	return m.Round(scale).UnscaledBig().Int64()
}

// Cmp compares the provided MonetaryAmount to the current MonetaryAmount and return -1, 0, 1 based on the underlying inf.Dec.Cmp function. 
// Ignores Currency.
func (m MonetaryAmount) Cmp(ma MonetaryAmount) int {
	return m.Dec.Cmp(ma.Dec)
}

func (m MonetaryAmount) IsZero() bool {
	var maZero = MonetaryAmount{Dec: inf.NewDec(0, 2), Currency: m.Currency}
	return m.Cmp(maZero) == 0
}

//UnmarshalText mutates the current MonetaryAmount instance with the provided byte array. 
func (m *MonetaryAmount) UnmarshalText(data []byte) error {
	m.Dec = new(inf.Dec)
	return m.Dec.UnmarshalText(data)
}

func (m MonetaryAmount) StringWithCurrency() string {
	c := string(m.Currency)
	if c != "" {
		c = " " + c
	}
	return fmt.Sprintf("%s%s", m.Dec.String(), c)
}

func (m MonetaryAmount) String() string {
	return m.Dec.String()
}

func NewMonetaryAmount(currency currencies.Currency) MonetaryAmount {
	if currency == currencies.Currency("") {
		panic("Currency cannot be empty")
	}

	return MonetaryAmount{Dec: decZero, Currency: currency}
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
