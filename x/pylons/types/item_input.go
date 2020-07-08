package types

import "errors"

// ItemInput is a wrapper struct for Item for recipes
type ItemInput struct {
	Doubles     DoubleInputParamList
	Longs       LongInputParamList
	Strings     StringInputParamList
	TransferFee FeeInputParam
}

// ItemInputList is a list of ItemInputs for convinience
type ItemInputList []ItemInput

// Validate is a function to check ItemInputList is valid
func (iil ItemInputList) Validate() error {
	return nil
}

// TradeItemInput is a wrapper struct for Item for trades
type TradeItemInput struct {
	ItemInput  ItemInput
	CookbookID string
}

// TradeItemInputList is a list of ItemInputs for convinience
type TradeItemInputList []TradeItemInput

// Validate is a function to check ItemInputList is valid
func (tiil TradeItemInputList) Validate() error {
	for _, ii := range tiil {
		if ii.CookbookID == "" {
			return errors.New("There should be no empty cookbook ID inputs for trades")
		}
	}
	return nil
}
