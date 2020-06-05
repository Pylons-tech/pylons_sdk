package types

// ItemInput is a wrapper struct for Item for recipes
type ItemInput struct {
	Doubles DoubleInputParamList
	Longs   LongInputParamList
	Strings StringInputParamList
}

// ItemInputList is a list of ItemInputs for convinience
type ItemInputList []ItemInput

// Validate is a function to check ItemInputList is valid
func (iil ItemInputList) Validate() error {
	return nil
}
