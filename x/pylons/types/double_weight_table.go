package types

// DoubleWeightTable describes weight loot table that produce double value
type DoubleWeightTable struct {
	WeightRanges []DoubleWeightRange
}

// DoubleWeightRange describes weight range that produce double value
type DoubleWeightRange struct {
	Lower  FloatString // This is added due to amino.Marshal does not support float variable
	Upper  FloatString
	Weight int
}
