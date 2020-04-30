package types

// DoubleParam describes the bounds on an item input/output parameter of type float64
type DoubleParam struct {
	// The likelihood that this parameter is applied to the output item. Between 0.0 (exclusive) and 1.0 (inclusive).
	Rate FloatString
	Key  string
	DoubleWeightTable
	// When program is not empty, DoubleWeightTable is ignored
	Program string
}

// DoubleParamList is a list of DoubleParam
type DoubleParamList []DoubleParam
