package pylons

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// AppModuleBasic is app module basics object
type AppModuleBasic struct{}

// Name returns AppModuleBasic name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec implements RegisterCodec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis return GenesisState in JSON
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis do validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return nil
}

// GetQueryCmd get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	pylonsQueryCmd := &cobra.Command{
		Use:   RouterKey,
		Short: "Querying commands for the pylons module",
	}

	return pylonsQueryCmd
}

// GetTxCmd get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	pylonsTxCmd := &cobra.Command{
		Use:   RouterKey,
		Short: "Pylons transactions subcommands",
	}

	return pylonsTxCmd
}

// RegisterRESTRoutes rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
}
