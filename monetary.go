package gowork

import (
	"fmt"
	inf "gopkg.in/inf.v0"
	"strings"
)

//Minimal hiding of the underlying Decimal implementation so API changes or a when better implementation comes along it won't cause huge app changes.
type MonetaryAmount struct {
	Dec *inf.Dec
}

var Zero = MonetaryAmount{inf.NewDec(0, 0)}

func (m *MonetaryAmount) Add(ma MonetaryAmount) (nma MonetaryAmount) {
	m.initialize()
	d := new(inf.Dec).Add(m.Dec, ma.Dec)
	nma = MonetaryAmount{d}
	return
}

func (m *MonetaryAmount) Neg() (nma MonetaryAmount) {
	m.initialize()
	d := new(inf.Dec).Neg(m.Dec)
	nma = MonetaryAmount{d}
	return
}

func (m *MonetaryAmount) Round(scale int) string {
	m.initialize()
	return new(inf.Dec).Round(m.Dec, inf.Scale(scale), inf.RoundHalfUp).String()
}

func (m *MonetaryAmount) AssumeScale(scale int) string {
	m.initialize()
	return strings.Replace(m.Round(scale), ".", "", 1)
}

func (m *MonetaryAmount) String() string {
	m.initialize()
	return m.Dec.String()
}

func (m *MonetaryAmount) MarshalText() ([]byte, error) {
	m.initialize()
	return m.Dec.MarshalText()
}

func (m *MonetaryAmount) UnmarshalText(data []byte) error {
	m.initialize()
	return m.Dec.UnmarshalText(data)
}

func (m *MonetaryAmount) initialize() {
	if m.Dec == nil {
		m.Dec = new(inf.Dec)
	}
}

func MonetaryAmountFromString(value string) (amount MonetaryAmount, err error) {
	dec, ok := new(inf.Dec).SetString(value)
	if !ok {
		err = fmt.Errorf("Invalid string input: %v", value)
	}

	amount = MonetaryAmount{dec}
	return
}

func MarshalMonetaryAmount(amount MonetaryAmount) *inf.Dec {
	return amount.Dec
}
