package types

import (
	"fmt"
)

// FeeInputParam describes the bounds on an item input/output parameter of type int64
type FeeInputParam struct {
	MinValue int
	MaxValue int
}

func (lp FeeInputParam) String() string {
	return fmt.Sprintf(`
	FeeInputParam{ 
		MinValue: %d,
		MaxValue: %d,
	}`, lp.MinValue, lp.MaxValue)
}

// Has validate if input is between min max range
func (lp FeeInputParam) Has(input int) bool {
	return input >= lp.MinValue && input <= lp.MaxValue
}
