package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Execution is a recipe execution used for tracking the execution - specifically a
// scheduled execution
type Execution struct {
	ID          string
	RecipeID    string // the recipe guid
	CookbookID  string
	CoinInputs  sdk.Coins
	ItemInputs  []Item
	BlockHeight int64
	Sender      sdk.AccAddress
	Completed   bool
}

// ExecutionList describes executions list
type ExecutionList struct {
	Executions []Execution
}

// NewExecution return a new Execution
func NewExecution(rcpID string, cbID string, ci sdk.Coins,
	itemInputs []Item,
	blockHeight int64, sender sdk.AccAddress,
	completed bool) Execution {

	exec := Execution{
		RecipeID:    rcpID,
		CookbookID:  cbID,
		CoinInputs:  ci,
		ItemInputs:  itemInputs,
		BlockHeight: blockHeight,
		Sender:      sender,
		Completed:   completed,
	}

	exec.ID = KeyGen(sender)
	return exec
}
