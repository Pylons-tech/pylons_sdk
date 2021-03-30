package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewRecipe creates a new recipe
func NewRecipe(recipeName, cookbookID, description string,
	coinInputs CoinInputList, // coins to put on the recipe
	itemInputs ItemInputList, // items to put on the recipe
	entries EntriesList, // items that can be created from recipe
	outputs WeightedOutputsList, // item outputs listing by weight value
	blockInterval int64, // The amount of time to wait to finish running the recipe
	sender sdk.AccAddress) Recipe {
	rcp := Recipe{
		NodeVersion:   SemVer("0.0.1"),
		Name:          recipeName,
		CookbookID:    cookbookID,
		CoinInputs:    coinInputs,
		ItemInputs:    itemInputs,
		Entries:       entries,
		Outputs:       outputs,
		BlockInterval: blockInterval,
		Description:   description,
		Sender:        sender,
	}

	rcp.ID = KeyGen(sender)
	return rcp
}

// GetItemInputRefIndex get item input index from ref string
func (rcp Recipe) GetItemInputRefIndex(inputRef string) int {
	for idx, input := range rcp.ItemInputs {
		if input.ID == inputRef {
			return idx
		}
	}
	return -1
}
