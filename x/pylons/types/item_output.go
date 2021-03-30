package types

// NewItemModifyOutput returns ItemOutput that is modified from item input
func NewItemModifyOutput(ID string, ItemInputRef string, ModifyParams ItemModifyParams) ItemModifyOutput {
	return ItemModifyOutput{
		ID:           ID,
		ItemInputRef: ItemInputRef,
		Doubles:      ModifyParams.Doubles,
		Longs:        ModifyParams.Longs,
		Strings:      ModifyParams.Strings,
		TransferFee:  ModifyParams.TransferFee,
	}
}

// SetTransferFee set generate item's transfer fee
func (mit *ItemModifyOutput) SetTransferFee(transferFee int64) {
	mit.TransferFee = transferFee
}

// NewItemOutput returns new ItemOutput generated from recipe
func NewItemOutput(ID string, Doubles DoubleParamList, Longs LongParamList, Strings StringParamList, TransferFee int64) ItemOutput {
	return ItemOutput{
		ID:          ID,
		Doubles:     Doubles,
		Longs:       Longs,
		Strings:     Strings,
		TransferFee: TransferFee,
	}
}

// SetTransferFee set generate item's transfer fee
func (io *ItemOutput) SetTransferFee(transferFee int64) {
	io.TransferFee = transferFee
}
