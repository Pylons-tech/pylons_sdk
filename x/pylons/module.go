package pylons

import (
	"encoding/json"

	"github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
)

// AppModuleBasic is app module basics object
type AppModuleBasic struct{}

// Name returns AppModuleBasic name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// DefaultGenesis return GenesisState in JSON
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis do validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, cl client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}

// GetQueryCmd get the root query command of this module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

// GetTxCmd get the root tx command of this module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
}

func (AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	RegisterCodec(amino)
}

func (AppModuleBasic) RegisterInterfaces(registry sdktypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&msgs.MsgCreateAccount{},
		&msgs.MsgGetPylons{},
		&msgs.MsgGoogleIAPGetPylons{},
		&msgs.MsgSendCoins{},
		&msgs.MsgSendItems{},
		&msgs.MsgCreateCookbook{},
		&msgs.MsgUpdateCookbook{},
		&msgs.MsgCreateRecipe{},
		&msgs.MsgUpdateRecipe{},
		&msgs.MsgExecuteRecipe{},
		&msgs.MsgDisableRecipe{},
		&msgs.MsgEnableRecipe{},
		&msgs.MsgCheckExecution{},
		&msgs.MsgFiatItem{},
		&msgs.MsgUpdateItemString{},
		&msgs.MsgCreateTrade{},
		&msgs.MsgFulfillTrade{},
		&msgs.MsgDisableTrade{},
		&msgs.MsgEnableTrade{},
	)

	msgs.RegisterMsgServiceDesc(registry)
}

// RegisterRESTRoutes rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
}
