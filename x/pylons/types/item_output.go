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

// // MarshalJSON is a custom marshal function
// func (io *ItemOutput) MarshalJSON() ([]byte, error) {
// 	sio := SerializeItemOutput{
// 		ModifyItem: SerializeModifyItemType{
// 			ItemInputRef: nil,
// 			Doubles:      io.ModifyItem.Doubles,
// 			Longs:        io.ModifyItem.Longs,
// 			Strings:      io.ModifyItem.Strings,
// 			TransferFee:  io.ModifyItem.TransferFee,
// 		},
// 		Doubles:     io.Doubles,
// 		Longs:       io.Longs,
// 		Strings:     io.Strings,
// 		TransferFee: io.TransferFee,
// 	}
// 	if io.ModifyItem.ItemInputRef != -1 {
// 		sio.ModifyItem.ItemInputRef = &io.ModifyItem.ItemInputRef
// 	}
// 	return json.Marshal(sio)
// }

// // UnmarshalJSON is a custom unmarshal function
// func (io *ItemOutput) UnmarshalJSON(data []byte) error {
// 	sio := SerializeItemOutput{}
// 	err := json.Unmarshal(data, &sio)
// 	if err != nil {
// 		return err
// 	}
// 	if sio.ModifyItem.ItemInputRef == nil {
// 		io.ModifyItem.ItemInputRef = -1
// 	} else {
// 		io.ModifyItem.ItemInputRef = *sio.ModifyItem.ItemInputRef
// 	}
// 	io.ModifyItem.Doubles = sio.ModifyItem.Doubles
// 	io.ModifyItem.Longs = sio.ModifyItem.Longs
// 	io.ModifyItem.Strings = sio.ModifyItem.Strings
// 	io.TransferFee = sio.ModifyItem.TransferFee
// 	io.Doubles = sio.Doubles
// 	io.Longs = sio.Longs
// 	io.Strings = sio.Strings
// 	io.TransferFee = sio.TransferFee
// 	return nil
// }
