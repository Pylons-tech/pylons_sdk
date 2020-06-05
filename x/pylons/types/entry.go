package types

import (
	"encoding/json"
)

// Entry describes an output which can be produced from a recipe
type Entry interface{}

// EntriesList is a struct to keep list of items and coins
type EntriesList []Entry

type serializeEntriesList struct {
	CoinOutputs []CoinOutput
	ItemOutputs []ItemOutput
}

func (wpl EntriesList) MarshalJSON() ([]byte, error) {
	var sel serializeEntriesList
	for _, wp := range wpl {
		switch wp := wp.(type) {
		case CoinOutput:
			sel.CoinOutputs = append(sel.CoinOutputs, wp)
		case ItemOutput:
			sel.ItemOutputs = append(sel.ItemOutputs, wp)
		default:
		}
	}
	return json.Marshal(sel)
}

func (el *EntriesList) UnmarshalJSON(data []byte) error {
	var sel serializeEntriesList
	err := json.Unmarshal(data, &sel)
	if err != nil {
		return err
	}

	for _, co := range sel.CoinOutputs {
		*el = append(*el, co)
	}
	for _, io := range sel.ItemOutputs {
		*el = append(*el, io)
	}
	return nil
}
