package types

// LongParam describes the bounds on an item input/output parameter of type int64
type LongParam struct {
	Key string

	// The likelihood that this parameter is applied to the output item. Between 0.0 (exclusive) and 1.0 (inclusive).
	Rate FloatString
	IntWeightTable
	// When program is not empty, IntWeightTable is ignored
	Program string
}

// LongParamList is a list of LongParam
type LongParamList []LongParam
