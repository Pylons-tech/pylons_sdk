package types

// ItemModifyParams describes the fields that needs to be modified
type ItemModifyParams struct {
	Doubles     DoubleParamList
	Longs       LongParamList
	Strings     StringParamList
	TransferFee int64
}

// ItemModifyOutput describes what is modified from item input
type ItemModifyOutput struct {
	ItemInputRef int
	Doubles      DoubleParamList
	Longs        LongParamList
	Strings      StringParamList
	TransferFee  int64
}

// SerializeModifyItemType describes the serialized format of ModifyItemType
type SerializeModifyItemType struct {
	ItemInputRef *int `json:",omitempty"`
	Doubles      DoubleParamList
	Longs        LongParamList
	Strings      StringParamList
	TransferFee  int64
}

// ItemOutput models the continuum of valid outcomes for item generation in recipes
type ItemOutput struct {
	Doubles     DoubleParamList
	Longs       LongParamList
	Strings     StringParamList
	TransferFee int64
}

// NewItemModifyOutput returns ItemOutput that is modified from item input
func NewItemModifyOutput(ItemInputRef int, ModifyParams ItemModifyParams) ItemModifyOutput {
	return ItemModifyOutput{
		ItemInputRef: ItemInputRef,
		Doubles:      ModifyParams.Doubles,
		Longs:        ModifyParams.Longs,
		Strings:      ModifyParams.Strings,
		TransferFee:  ModifyParams.TransferFee,
	}
}

// NewItemOutput returns new ItemOutput generated from recipe
func NewItemOutput(Doubles DoubleParamList, Longs LongParamList, Strings StringParamList, TransferFee int64) ItemOutput {
	return ItemOutput{
		Doubles:     Doubles,
		Longs:       Longs,
		Strings:     Strings,
		TransferFee: TransferFee,
	}
}

// SerializeItemOutput describes the item output in serialize format
type SerializeItemOutput struct {
	ModifyItem  SerializeModifyItemType
	Doubles     DoubleParamList
	Longs       LongParamList
	Strings     StringParamList
	TransferFee int64
}
