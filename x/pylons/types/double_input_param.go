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
