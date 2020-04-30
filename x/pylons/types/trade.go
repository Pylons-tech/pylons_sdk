package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

const TypeTrade = "trade"

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

func (cbl TradeList) String() string {
	output := "TradeList{"
	for _, cb := range cbl.Trades {
		output += cb.String()
		output += ",\n"
	}
	output += "}"
	return output
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

	trd.ID = trd.KeyGen()
	return trd
}

// KeyGen generates key for the store
func (trd Trade) KeyGen() string {
	id := uuid.New()
	return trd.Sender.String() + id.String()
}
