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

// NewGenesisState returns new genesis state
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// ValidateGenesis do validate genesis
func ValidateGenesis(data GenesisState) error {
	return nil
}

// DefaultGenesisState returns default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Cookbooks: []types.Cookbook{},
		Recipies:  []types.Recipe{},
		Items:     []types.Item{},
	}
}
