package types

// LongInputParam describes the bounds on an item input/output parameter of type int64
type LongInputParam struct {
	Key      string
	MinValue int
	MaxValue int
}

// LongInputParamList is a list of LongInputParam
type LongInputParamList []LongInputParam
