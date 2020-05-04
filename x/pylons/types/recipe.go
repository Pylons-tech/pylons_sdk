package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

const TypeRecipe = "recipe"

// Recipe is a game state machine step abstracted out as a cooking terminology
type Recipe struct {
	ID            string // the recipe guid
	CookbookID    string // the cookbook guid
	Name          string
	CoinInputs    CoinInputList
	ItemInputs    ItemInputList
	Entries       EntriesList
	Outputs       WeightedOutputsList
	Description   string
	BlockInterval int64
	Sender        sdk.AccAddress
	Disabled      bool
}

// RecipeList is a list of recipes
type RecipeList struct {
	Recipes []Recipe
}

// NewRecipe creates a new recipe
func NewRecipe(recipeName, cookbookID, description string,
	coinInputs CoinInputList, // coinOutputs CoinOutputList,
	itemInputs ItemInputList, // itemOutputs ItemOutputList,
	entries EntriesList,
	outputs WeightedOutputsList,
	execTime int64, sender sdk.AccAddress) Recipe {
	rcp := Recipe{
		Name:          recipeName,
		CookbookID:    cookbookID,
		CoinInputs:    coinInputs,
		ItemInputs:    itemInputs,
		Entries:       entries,
		Outputs:       outputs,
		BlockInterval: execTime,
		Description:   description,
		Sender:        sender,
	}

	rcp.ID = rcp.KeyGen()
	return rcp
}

// KeyGen generates key for the store
func (rcp Recipe) KeyGen() string {
	id := uuid.New()
	return rcp.Sender.String() + id.String()
}