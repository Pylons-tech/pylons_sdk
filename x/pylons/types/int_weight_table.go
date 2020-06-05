package types

// IntWeightTable describes weight loot table that produce int value
type IntWeightTable struct {
	WeightRanges []IntWeightRange
}

// IntWeightRange describes weight range that produce int value
type IntWeightRange struct {
	Lower  int
	Upper  int
	Weight int
}
