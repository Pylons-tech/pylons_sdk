package handlers

import (
	"github.com/Pylons-tech/pylons_sdk/x/pylons/keep"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExecProcess struct {
	ctx          sdk.Context
	keeper       keep.Keeper
	recipe       types.Recipe
	matchedItems []types.Item
	ec           types.CelEnvCollection
}
