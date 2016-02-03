package gowork

import (
	"fmt"
	inf "gopkg.in/inf.v0"
	"strings"
)

//Minimal hiding of the underlying Decimal implementation so API changes or a when better implementation comes along it won't cause huge app changes.
type MonetaryAmount struct {
	*inf.Dec
}

var Zero = MonetaryAmount{inf.NewDec(0, 2)}

func (m MonetaryAmount) Add(ma MonetaryAmount) MonetaryAmount {
	return MonetaryAmount{new(inf.Dec).Add(m.Dec, ma.Dec)}
}

func (m MonetaryAmount) Multiply(i *inf.Dec) MonetaryAmount {
	return MonetaryAmount{new(inf.Dec).Mul(m.Dec, i)}
}

func (m MonetaryAmount) Neg() MonetaryAmount {
	return MonetaryAmount{new(inf.Dec).Neg(m.Dec)}
}

func (m MonetaryAmount) Round(scale int) string {
	return new(inf.Dec).Round(m.Dec, inf.Scale(scale), inf.RoundHalfUp).String()
}

func (m MonetaryAmount) AssumeScale(scale int) string {
	return strings.Replace(m.Round(scale), ".", "", 1)
}

func (m MonetaryAmount) Cmp(ma MonetaryAmount) int {
	return m.Dec.Cmp(ma.Dec)
}

func (m MonetaryAmount) IsZero() bool {
	return m.Cmp(Zero) == 0
}

//UnmarshalText mutates the current MonetaryAmount instance with the provided byte array. 
func (m *MonetaryAmount) UnmarshalText(data []byte) error {
	m.Dec = new(inf.Dec)
	return m.Dec.UnmarshalText(data)
}

func MonetaryAmountFromString(value string) (amount MonetaryAmount, err error) {
	dec, ok := new(inf.Dec).SetString(value)
	if !ok {
		return Zero, fmt.Errorf("Invalid string input: %v", value)
	}

	return MonetaryAmount{dec}, nil
}
