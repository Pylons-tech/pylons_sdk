package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Trade is a construct to perform exchange of items and coins between users. Initiated by the sender and completed by
// the FulFiller.
type Trade struct {
	ID          string // the recipe guid
	CoinInputs  CoinInputList
	ItemInputs  ItemInputList
	CoinOutputs sdk.Coins
	ItemOutputs ItemList
	ExtraInfo   string
	Sender      sdk.AccAddress
	FulFiller   sdk.AccAddress
	Disabled    bool
	Completed   bool
}

// TradeList is a list of trades
type TradeList struct {
	Trades []Trade
}

// NewTrade creates a new trade
func NewTrade(extraInfo string,
	coinInputs CoinInputList, // coinOutputs CoinOutputList,
	itemInputs ItemInputList, // itemOutputs ItemOutputList,
	coinOutputs sdk.Coins, // newly created param instead of coinOutputs and itemOutputs
	itemOutputs ItemList,
	sender sdk.AccAddress) Trade {
	trd := Trade{
		CoinInputs:  coinInputs,
		ItemInputs:  itemInputs,
		CoinOutputs: coinOutputs,
		ItemOutputs: itemOutputs,
		ExtraInfo:   extraInfo,
		Sender:      sender,
	}

	trd.ID = KeyGen(sender)
	return trd
}
