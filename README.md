# pylons SDK

pylons SDK provides packages to build blockchain games on pylons eco system.

## Setup development environment

```
git clone https://github.com/Pylons-tech/pylons_sdk
brew install pre-commit
brew install golangci/tap/golangci-lint
go get -u golang.org/x/lint/golint
pre-commit install
```

# SDK publish preparation

#### Check fixture test works fine
```
make fixture_tests
```

#### All the features added should have fixture test and it should be well documented.

# Packages

## Fixture Test Package
github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test

## Integration Test Utils Package
github.com/Pylons-tech/pylons_sdk/cmd/test

#### GetAccountAddr  
GetAccountAddr is a function to get account address from key
#### GetAccountInfoFromName  
GetAccountInfoFromName is a function to get account information from account key
#### ListItemsViaCLI  
ListItemsViaCLI is a function to list items via cli
#### GetDaemonStatus  
GetDaemonStatus is a function to get daemon status
#### CLIOpts  
CLIOpts is a variable to manage pylonscli options.  
`CustomNode` is for custom node tcp endpoint and `RestEndpoint` is for custom node http endpoint.
#### WaitAndGetTxData  
WaitAndGetTxData is a function to get transaction data after transaction is processed
#### ReadFile  
ReadFile is a utility function to read file
#### CleanFile  
CleanFile is a function to remove file
#### GetAminoCdc  
GetAminoCdc is a utility function to get amino codec
#### RunPylonsCli  
RunPylonsCli is a function to run pylonscli
#### GenTxWithMsg  
GenTxWithMsg is a function to generate transaction from msg
#### WaitForNextBlock  
WaitForNextBlock is a function to wait until next block
#### WaitAndGetTxData  
WaitAndGetTxData is a function to get transaction data after transaction is processed
#### GetHumanReadableErrorFromTxHash  
GetHumanReadableErrorFromTxHash is a function to get human readable error from txhash
#### TestTxWithMsgWithNonce  
TestTxWithMsgWithNonce is a function to send transaction with message and nonce
#### GetItemByGUID  
GetItemByGUID is to get Item from ID
#### SendMultiMsgTxWithNonce  
SendMultiMsgTxWithNonce is an integration test utility to send multiple message transaction from a single sender, single signed transaction.
#### RegisterDefaultActionRunners  
RegisterDefaultActionRunners register default test functions.
#### RegisterActionRunner  
RegisterActionRunner registers action runner function
#### GetActionRunner  
GetActionRunner get registered action runner function
#### RunActionRunner  
RunActionRunner execute registered action runner function

## Handlers struct package
github.com/Pylons-tech/pylons_sdk/x/pylons/handlers

### Structs

#### ExecuteRecipeResp  
ExecuteRecipeResp is the response for executeRecipe
#### ExecuteRecipeScheduleOutput  
ExecuteRecipeScheduleOutput is a struct that shows how execute recipe schedule output works
#### CheckExecutionResp  
CheckExecutionResp is the response for checkExecution
#### CreateCBResponse  
CreateCBResponse is a struct of create cookbook response
#### CreateRecipeResponse  
CreateRecipeResponse is struct of create recipe response
#### FulfillTradeResp  
FulfillTradeResp is a struct to control fulfill trade response
#### PopularRecipeType  
PopularRecipeType is a type for popular recipes
#### GetParamsForPopularRecipe  
GetParamsForPopularRecipe is a function to get popular recipe's attributes
#### FiatItemResponse  
FiatItemResponse is a struct to control fiat item response
#### UpdateItemStringResp  
UpdateItemStringResp is a struct to control update item string response
## Msgs package
github.com/Pylons-tech/pylons_sdk/x/pylons/msgs

### Msg structs

#### MsgCheckExecution  
MsgCheckExecution defines a CheckExecution message
#### MsgCreateCookbook  
MsgCreateCookbook defines a CreateCookbook message
#### MsgCreateRecipe  
MsgCreateRecipe defines a CreateRecipe message
#### MsgCreateTrade  
MsgCreateTrade defines a CreateTrade message
#### MsgDisableRecipe  
MsgDisableRecipe defines a DisableRecipe message
#### MsgDisableTrade  
MsgDisableTrade defines a DisableTrade message
#### MsgEnableRecipe  
MsgEnableRecipe defines a EnableRecipe message
#### MsgEnableTrade  
MsgEnableTrade defines a EnableTrade message
#### MsgExecuteRecipe  
MsgExecuteRecipe defines a SetName message
#### MsgFiatItem  
MsgFiatItem is a msg struct to be used to fiat item
#### MsgFulfillTrade  
MsgFulfillTrade defines a FulfillTrade message
#### MsgGetPylons  
MsgGetPylons defines a GetPylons message
#### MsgSendPylons  
MsgSendPylons defines a SendPylons message
#### MsgUpdateItemString  
MsgUpdateItemString defines a SendPylons message
#### MsgUpdateCookbook  
MsgUpdateCookbook defines a UpdateCookbook message
#### MsgUpdateRecipe  
MsgUpdateRecipe defines a UpdateRecipe message

### Msg Utility functions 

#### NewMsgGetPylons  
NewMsgGetPylons is a function to get MsgGetPylons msg from required params

#### NewMsgExecuteRecipe  
NewMsgExecuteRecipe a constructor for ExecuteCookbook msg

#### NewMsgCreateCookbook  
NewMsgCreateCookbook a constructor for CreateCookbook msg

#### NewMsgGetPylons  
NewMsgGetPylons is a function to get MsgGetPylons msg from required params

#### NewMsgUpdateItemString  
NewMsgUpdateItemString is a function to get MsgUpdateItemString msg from required params

#### NewMsgCreateTrade  
NewMsgCreateTrade a constructor for CreateTrade msg

#### NewMsgFulfillTrade  
NewMsgFulfillTrade a constructor for FulfillTrade msg

#### NewMsgDisableTrade  
NewMsgDisableTrade a constructor for DisableTrade msg

#### NewMsgCheckExecution   
NewMsgCheckExecution a constructor for CheckExecution msg

#### NewMsgFiatItem  
NewMsgFiatItem a constructor for MsgFiatItem msg

#### NewMsgCreateRecipe  
NewMsgCreateRecipe a constructor for CreateRecipe msg

#### DefaultCostPerBlock  
DefaultCostPerBlock the amount of pylons to be charged by default

## Types package
github.com/Pylons-tech/pylons_sdk/x/pylons/types


### Type structs

#### Item  
Item is a tradable asset

#### Cookbook  
Cookbook is a struct that contains all the metadata of a cookbook

#### Recipe  
Recipe is a game state machine step abstracted out as a cooking terminology

#### Trade  
Trade is a construct to perform exchange of items and coins between users. Initiated by the sender and completed by the FulFiller.

#### FloatString  
FloatString is a wrapper to resolve the amino issues

#### EntriesList  
EntriesList is a struct to keep list of items and coins

#### TradeList  
TradeList is a list of trades

#### Execution  
Execution is a recipe execution used for tracking the execution #### specifically a scheduled execution

#### CoinOutput  
CoinOutput is the game elements that are needs as output to a recipe

#### ItemModifyParams  
ItemModifyParams describes the fields that needs to be modified

#### BasicTier  
BasicTier is the cookbook tier which doesn't allow paid receipes which means the developers cannot have receipes where they can actually carge a fee in pylons.

#### PremiumTier  
PremiumTier the cookbook tier which does allow paid receipes

#### ItemList   
ItemList is a list of items

#### ItemInput  
ItemInput is a wrapper struct for Item for recipes

#### ItemInputList  
ItemInputList is a list of ItemInputs for convinience

#### DoubleInputParamList  
DoubleInputParamList is a list of DoubleInputParam

#### DoubleInputParam  
DoubleInputParamList is a list of DoubleInputParam

#### LongInputParam  
LongInputParam describes the bounds on an item input/output parameter of type int64

#### LongInputParamList  
LongInputParamList is a list of LongInputParam

#### StringInputParam  
StringInputParam describes the bounds on an item input/output parameter of type string

#### StringInputParamList  
StringInputParamList is a list of StringInputParam

#### CoinInputList  
CoinInputList is a list of Coin inputs

#### WeightedOutputsList  
WeightedOutputsList is a struct to keep items which can be generated by weight

### Type Utility functions 

#### NewPylon  
NewPylon Returns pylon currency

#### GenItemInputList  
GenItemInputList is a utility function to genearte item input list

#### GenEntries  
GenEntries is a function to generate entries from coin name and item name

#### GenCoinInputList  
GenCoinInputList is a utility function to genearte coin input list

#### GenItemOnlyEntry  
GenItemOnlyEntry is a utility function to generate item only entry

#### GenCoinInputList  
GenCoinInputList is a utility function to genearte coin input list

#### GenEntriesFirstItemNameUpgrade  
GenEntriesFirstItemNameUpgrade is a function to generate entries that update first item's name

#### GenOneOutput  
GenOneOutput is a function to generate output with one from entry list

## App package
github.com/Pylons-tech/pylons_sdk/app

#### MakeCodec  
MakeCodec make codec for message marshal/unmarshal

## Queriers package
github.com/Pylons-tech/pylons_sdk/x/pylons/queriers

#### ExecResp  
ExecResp is the response for ListExecutions

#### ItemResp  
ItemResp is the response for Items

