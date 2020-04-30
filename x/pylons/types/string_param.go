package types

// StringParam describes an item input/output parameter of type string
type StringParam struct {
	// The likelihood that this parameter is applied to the output item. Between 0.0 (exclusive) and 1.0 (inclusive).
	Rate  FloatString
	Key   string
	Value string
	// When program is not empty, Value is ignored
	Program string
}

// StringParamList is a list of StringParam
type StringParamList []StringParam
