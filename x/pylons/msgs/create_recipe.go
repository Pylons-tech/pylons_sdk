package msgs

import (
	"encoding/json"

	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgCreateRecipe defines a CreateRecipe message
type MsgCreateRecipe struct {
	// optional RecipeID if someone
	RecipeID      string `json:",omitempty"`
	Name          string
	CookbookID    string // the cookbook guid
	CoinInputs    types.CoinInputList
	ItemInputs    types.ItemInputList
	Entries       types.EntriesList
	Outputs       types.WeightedOutputsList
	BlockInterval int64
	Sender        sdk.AccAddress
	Description   string
}

// NewMsgCreateRecipe a constructor for CreateRecipe msg
func NewMsgCreateRecipe(recipeName, cookbookID, recipeID, description string,
	coinInputs types.CoinInputList,
	itemInputs types.ItemInputList,
	entries types.EntriesList,
	outputs types.WeightedOutputsList,
	blockInterval int64,
	sender sdk.AccAddress) MsgCreateRecipe {
	return MsgCreateRecipe{
		Name:          recipeName,
		CookbookID:    cookbookID,
		RecipeID:      recipeID,
		Description:   description,
		CoinInputs:    coinInputs,
		ItemInputs:    itemInputs,
		Entries:       entries,
		Outputs:       outputs,
		BlockInterval: int64(blockInterval),
		Sender:        sender,
	}
}

// Route should return the name of the module
func (msg MsgCreateRecipe) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateRecipe) Type() string { return "create_recipe" }

// ValidateBasic validates the Msg
func (msg MsgCreateRecipe) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}

	if len(msg.Description) < 20 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "the description should have more than 20 characters")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateRecipe) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners gets the signer who should have signed the message
func (msg MsgCreateRecipe) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
