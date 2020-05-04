package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

const TypeExecution = "execution"

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

	exec.ID = exec.KeyGen()
	return exec
}

// KeyGen generates key for the execution
func (exec Execution) KeyGen() string {
	id := uuid.New()
	return exec.Sender.String() + id.String()
}