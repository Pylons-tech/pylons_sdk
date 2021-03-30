package pylons

import (
	"encoding/json"
)

// AppModuleBasic is app module basics object
type AppModuleBasic struct{}

// Name returns AppModuleBasic name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// DefaultGenesis return GenesisState in JSON
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis do validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return nil
}
