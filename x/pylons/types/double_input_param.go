package types

// DoubleInputParam describes the bounds on an item input/output parameter of type float64
type DoubleInputParam struct {
	Key string
	// The minimum legal value of this parameter.
	MinValue FloatString
	// The maximum legal value of this parameter.
	MaxValue FloatString
}

// DoubleInputParamList is a list of DoubleInputParam
type DoubleInputParamList []DoubleInputParam

// Has check if an input is between double input param range
func (dp DoubleInputParam) Has(input float64) bool {
	return input >= dp.MinValue.Float() && input <= dp.MaxValue.Float()
}
