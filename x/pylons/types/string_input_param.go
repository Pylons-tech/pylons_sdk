package types

// StringInputParam describes the bounds on an item input/output parameter of type int64
type StringInputParam struct {
	Key string
	// The value of the parameter
	Value string
}

// StringInputParamList is a list of StringInputParam
type StringInputParamList []StringInputParam
