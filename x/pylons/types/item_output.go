package types

import (
	"encoding/json"
)

// ItemOutput models the continuum of valid outcomes for item generation in recipes

type ItemModifyParams struct {
	Doubles DoubleParamList
	Longs   LongParamList
	Strings StringParamList
}

type ModifyItemType struct {
	ItemInputRef int
	Doubles      DoubleParamList
	Longs        LongParamList
	Strings      StringParamList
}

type SerializeModifyItemType struct {
	ItemInputRef *int `json:",omitempty"`
	Doubles      DoubleParamList
	Longs        LongParamList
	Strings      StringParamList
}

type ItemOutput struct {
	ModifyItem ModifyItemType
	Doubles    DoubleParamList
	Longs      LongParamList
	Strings    StringParamList
}

func NewInputRefOutput(ItemInputRef int, ModifyParams ItemModifyParams) ItemOutput {
	return ItemOutput{
		ModifyItem: ModifyItemType{
			ItemInputRef: ItemInputRef,
			Doubles:      ModifyParams.Doubles,
			Longs:        ModifyParams.Longs,
			Strings:      ModifyParams.Strings,
		},
	}
}

func NewItemOutput(Doubles DoubleParamList, Longs LongParamList, Strings StringParamList) ItemOutput {
	return ItemOutput{
		ModifyItem: ModifyItemType{
			ItemInputRef: -1,
		},
		Doubles: Doubles,
		Longs:   Longs,
		Strings: Strings,
	}
}

type SerializeItemOutput struct {
	ModifyItem SerializeModifyItemType
	Doubles    DoubleParamList
	Longs      LongParamList
	Strings    StringParamList
}

func (io *ItemOutput) MarshalJSON() ([]byte, error) {
	sio := SerializeItemOutput{
		ModifyItem: SerializeModifyItemType{
			ItemInputRef: nil,
			Doubles:      io.ModifyItem.Doubles,
			Longs:        io.ModifyItem.Longs,
			Strings:      io.ModifyItem.Strings,
		},
		Doubles: io.Doubles,
		Longs:   io.Longs,
		Strings: io.Strings,
	}
	if io.ModifyItem.ItemInputRef != -1 {
		sio.ModifyItem.ItemInputRef = &io.ModifyItem.ItemInputRef
	}
	return json.Marshal(sio)
}

func (io *ItemOutput) UnmarshalJSON(data []byte) error {
	sio := SerializeItemOutput{}
	err := json.Unmarshal(data, &sio)
	if err != nil {
		return err
	}
	if sio.ModifyItem.ItemInputRef == nil {
		io.ModifyItem.ItemInputRef = -1
	} else {
		io.ModifyItem.ItemInputRef = *sio.ModifyItem.ItemInputRef
	}
	io.ModifyItem.Doubles = sio.ModifyItem.Doubles
	io.ModifyItem.Longs = sio.ModifyItem.Longs
	io.ModifyItem.Strings = sio.ModifyItem.Strings
	io.Doubles = sio.Doubles
	io.Longs = sio.Longs
	io.Strings = sio.Strings
	return nil
}
