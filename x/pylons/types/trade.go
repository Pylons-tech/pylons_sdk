package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Trade is a construct to perform exchange of items and coins between users. Initiated by the sender and completed by
// the FulFiller.
type Trade struct {
	ID          string             // the trade guid
	CoinInputs  CoinInputList      // coins that the fulfiller should send to creator
	ItemInputs  TradeItemInputList // items that the fulfiller should send to creator
	CoinOutputs sdk.Coins          // coins that the creator should send to fulfiller
	ItemOutputs ItemList           // items that the creator should send to fulfiller
	ExtraInfo   string             // custom trade info text
	Sender      sdk.AccAddress     // trade creator address
	FulFiller   sdk.AccAddress     // trade fulfiller address (acceptor)
	Disabled    bool               // disabled flag
	Completed   bool               // completed flag
}

// TradeList is a list of trades
type TradeList struct {
	Trades []Trade
}

// NewTrade creates a new trade
func NewTrade(extraInfo string,
	coinInputs CoinInputList, // coinOutputs CoinOutputList,
	itemInputs TradeItemInputList, // itemOutputs ItemOutputList,
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
