package pylons

import (
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
)

// GenesisState empty genesis for pylons
type GenesisState struct {
	Cookbooks []types.Cookbook
	Recipies  []types.Recipe
	Items     []types.Item
}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Cookbooks: []types.Cookbook{},
		Recipies:  []types.Recipe{},
		Items:     []types.Item{},
	}
}
