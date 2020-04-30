package types

type DoubleWeightTable struct {
	WeightRanges []DoubleWeightRange
}

type DoubleWeightRange struct {
	Lower  FloatString // This is added due to amino.Marshal does not support float variable
	Upper  FloatString
	Weight int
}
