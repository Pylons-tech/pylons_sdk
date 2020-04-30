# pylons SDK

pylons SDK provides packages to build blockchain games on pylons eco system.

# Packages

"github.com/Pylons-tech/pylons/cmd/fixtures_test"
"github.com/Pylons-tech/pylons/cmd/test"

e.g. Utility functions within pylons/cmd tests

- GetAccountAddr
- GetAccountInfoFromName
- ListItemsViaCLI
- GetDaemonStatus
- CLIOpts.CustomNode
- CLIOpts.RestEndpoint
- WaitAndGetTxData
- ReadFile
- GetAminoCdc
- RunPylonsCli
- CleanFile
- GenTxWithMsg
- WaitForNextBlock
- WaitAndGetTxData
- GetHumanReadableErrorFromTxHash
- TestTxWithMsgWithNonce
- GetItemByGUID
- TestTxWithMsgWithNonce




"github.com/Pylons-tech/pylons/x/pylons/handlers"

Structs

- handlers.ExecuteRecipeResp
- handlers.ExecuteRecipeScheduleOutput{}
- handlers.CheckExecutionResp{}
- handlers.CreateCBResponse{}
- handlers.CreateRecipeResponse{}
- handlers.FulfillTradeResp{}
- handlers.PopularRecipeType
- handlers.GetParamsForPopularRecipe
- handlers.FiatItemResponse{}
- handlers.UpdateItemStringResp{}

"github.com/Pylons-tech/pylons/x/pylons/msgs"

All msg types
- MsgCheckExecution
- MsgCreateCookbook
- MsgCreateRecipe
- MsgCreateTrade
- MsgDisableRecipe
- MsgDisableTrade
- MsgEnableRecipe
- MsgEnableTrade
- MsgExecuteRecipe
- MsgFiatItem
- MsgFulfillTrade
- MsgGetPylons
- MsgSendPylons
- MsgUpdateItemString
- MsgUpdateCookbook
- MsgUpdateRecipe

Utility functions 

- msgs.NewMsgGetPylons
- msgs.NewMsgExecuteRecipe
- msgs.NewMsgCreateCookbook
- msgs.NewMsgGetPylons
- msgs.NewMsgUpdateItemString
- msgs.NewMsgCreateTrade
- msgs.NewMsgFulfillTrade
- msgs.NewMsgDisableTrade
- msgs.NewMsgCheckExecution 
- msgs.NewMsgFiatItem
- msgs.NewMsgCreateRecipe
- msgs.DefaultCostPerBlock


"github.com/Pylons-tech/pylons/x/pylons/types"


structs

- types.Item
- types.Cookbook
- types.Recipe
- types.Trade
- types.FloatString
- types.EntriesList
- types.TradeList
- types.Execution
- types.CoinOutput
- types.ItemModifyParams
- types.PremiumTier.Fee
- types.ItemList 
- types.ItemInputList
- types.ItemInput
- types.DoubleInputParamList
- types.DoubleInputParam
- types.LongInputParamList
- types.LongInputParam
- types.StringInputParamList
- types.StringInputParam
- types.CoinInputList,
- types.WeightedOutputsList,

Utility functions 

- types.NewPylon
- types.GenItemInputList
- types.GenEntries
- types.GenCoinInputList
- types.ItemInputList{},
- types.GenItemOnlyEntry
- types.GenCoinInputList
- types.GenEntriesFirstItemNameUpgrade(desItemName),
- types.GenOneOutput

"github.com/Pylons-tech/pylons/app"

- app.MakeCodec()

"github.com/Pylons-tech/pylons/x/pylons/queriers"

- queriers.ExecResp
- queriers.ItemResp
