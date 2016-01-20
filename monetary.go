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

var Zero = MonetaryAmount{inf.NewDec(0, 0)}

func (m MonetaryAmount) Add(ma MonetaryAmount) MonetaryAmount {
	return MonetaryAmount{new(inf.Dec).Add(m.Dec, ma.Dec)}
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

func MarshalMonetaryAmount(amount MonetaryAmount) *inf.Dec {
	return amount.Dec
}
