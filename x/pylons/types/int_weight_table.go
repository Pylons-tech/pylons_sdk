package types

type IntWeightTable struct {
	WeightRanges []IntWeightRange
}

type IntWeightRange struct {
	Lower  int
	Upper  int
	Weight int
}
