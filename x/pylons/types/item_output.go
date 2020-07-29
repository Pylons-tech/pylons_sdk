package types

import "fmt"

// ItemModifyParams describes the fields that needs to be modified
type ItemModifyParams struct {
	Doubles     DoubleParamList
	Longs       LongParamList
	Strings     StringParamList
	TransferFee int64
}

// ItemModifyOutput describes what is modified from item input
type ItemModifyOutput struct {
	ID           string
	ItemInputRef int
	Doubles      DoubleParamList
	Longs        LongParamList
	Strings      StringParamList
	TransferFee  int64
}

// NewItemModifyOutput returns ItemOutput that is modified from item input
func NewItemModifyOutput(ID string, ItemInputRef int, ModifyParams ItemModifyParams) ItemModifyOutput {
	return ItemModifyOutput{
		ID:           ID,
		ItemInputRef: ItemInputRef,
		Doubles:      ModifyParams.Doubles,
		Longs:        ModifyParams.Longs,
		Strings:      ModifyParams.Strings,
		TransferFee:  ModifyParams.TransferFee,
	}
}

// GetID returns ID of coin output
func (mit ItemModifyOutput) GetID() string {
	return mit.ID
}

func (mit ItemModifyOutput) String() string {
	return fmt.Sprintf(`ItemModifyOutput{
		ID: %s,
		ItemInputRef: %d,
		Doubles: %+v,
		Longs:   %+v,
		Strings: %+v,
		TransferFee: %d,
	}`, mit.ID, mit.ItemInputRef, mit.Doubles, mit.Longs, mit.Strings, mit.TransferFee)
}

// SetTransferFee set generate item's transfer fee
func (mit *ItemModifyOutput) SetTransferFee(transferFee int64) {
	mit.TransferFee = transferFee
}

// GetTransferFee set item's TransferFee
func (mit ItemModifyOutput) GetTransferFee() int64 {
	return mit.TransferFee
}

// ItemOutput models the continuum of valid outcomes for item generation in recipes
type ItemOutput struct {
	ID          string
	Doubles     DoubleParamList
	Longs       LongParamList
	Strings     StringParamList
	TransferFee int64
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

// GetID returns ID of coin output
func (io ItemOutput) GetID() string {
	return io.ID
}

func (io ItemOutput) String() string {
	return fmt.Sprintf(`ItemOutput{
		ID: %s,
		Doubles: %+v,
		Longs:   %+v,
		Strings: %+v,
		TransferFee: %d,
	}`, io.ID, io.Doubles, io.Longs, io.Strings, io.TransferFee)
}

// SetTransferFee set generate item's transfer fee
func (io *ItemOutput) SetTransferFee(transferFee int64) {
	io.TransferFee = transferFee
}

// GetTransferFee set item's TransferFee
func (io ItemOutput) GetTransferFee() int64 {
	return io.TransferFee
}
