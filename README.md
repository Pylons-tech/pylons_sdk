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

# How to implement fixture test using fixture test SDK

As in `cmd/fixture_test/fixture_test.go` you can add single `*_test.go` file and add flag parser, register default action runners, custom action runners and run test scenarios.

```go
package fixturetest

import (
	"flag"
	"testing"

	inttestSDK "github.com/Pylons-tech/pylons_sdk/cmd/test"
	fixturetestSDK "github.com/Pylons-tech/pylons_sdk/cmd/fixture_test"
)

var runSerialMode bool = false
var useRest bool = false
var useKnownCookbook bool = false

func init() {
	flag.BoolVar(&runSerialMode, "runserial", false, "true/false value to check if test will be running in parallel")
	flag.BoolVar(&useRest, "userest", false, "use rest endpoint for Tx send")
	flag.BoolVar(&useKnownCookbook, "use-known-cookbook", false, "use existing cookbook or not")
}

func TestFixturesViaCLI(t *testing.T) {
	flag.Parse()
	fixturetestSDK.FixtureTestOpts.IsParallel = !runSerialMode
	fixturetestSDK.FixtureTestOpts.CreateNewCookbook = !useKnownCookbook
	if useRest {
		inttestSDK.CLIOpts.RestEndpoint = "http://localhost:1317"
	}
	inttest.CLIOpts.MaxBroadcast = 50
	fixturetestSDK.RegisterDefaultActionRunners()
	// Register custom action runners
	// fixturetestSDK.RegisterActionRunner("custom_action", CustomActionRunner)
	fixturetestSDK.RunTestScenarios("scenarios", t)
}

```


# Packages

## Fixture Test Package
github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test

## Integration Test Utils Package
github.com/Pylons-tech/pylons_sdk/cmd/test
| No | Type | Name                            | Description                                                                                                                                         |
|----|------|---------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| 1  | Config   | CLIOpts                         | CLIOpts is a variable to manage pylonscli options.`CustomNode` is for custom node tcp endpoint and `RestEndpoint` is for custom node http endpoint. |
| 2  | Fn   | CleanFile                       | CleanFile is a function to remove file                                                                                                              |
| 3  | Fn   | GenTxWithMsg                    | GenTxWithMsg is a function to generate transaction from msg                                                                                         |
| 4  | Fn   | GetAccountAddr                  | GetAccountAddr is a function to get account address from key                                                                                        |
| 5  | Fn   | GetAccountInfoFromName          | GetAccountInfoFromName is a function to get account information from account key                                                                    |
| 6  | Fn   | GetActionRunner                 | GetActionRunner get registered action runner function                                                                                               |
| 7  | Fn   | GetAminoCdc                     | GetAminoCdc is a utility function to get amino codec                                                                                                |
| 8  | Fn   | GetDaemonStatus                 | GetDaemonStatus is a function to get daemon status                                                                                                  |
| 9  | Fn   | GetHumanReadableErrorFromTxHash | GetHumanReadableErrorFromTxHash is a function to get human readable error from txhash                                                               |
| 10 | Fn   | GetItemByGUID                   | GetItemByGUID is to get Item from ID                                                                                                                |
| 11 | Fn   | ListItemsViaCLI                 | ListItemsViaCLI is a function to list items via cli                                                                                                 |
| 12 | Fn   | ReadFile                        | ReadFile is a utility function to read file                                                                                                         |
| 13 | Fn   | RegisterActionRunner            | RegisterActionRunner registers action runner function                                                                                               |
| 14 | Fn   | RegisterDefaultActionRunners    | RegisterDefaultActionRunners register default test functions.                                                                                       |
| 15 | Fn   | RunActionRunner                 | RunActionRunner execute registered action runner function                                                                                           |
| 16 | Fn   | RunPylonsCli                    | RunPylonsCli is a function to run pylonscli                                                                                                         |
| 17 | Fn   | SendMultiMsgTxWithNonce         | SendMultiMsgTxWithNonce is an integration test utility to send multiple message transaction from a single sender, single signed transaction.        |
| 18 | Fn   | TestTxWithMsgWithNonce          | TestTxWithMsgWithNonce is a function to send transaction with message and nonce                                                                     |
| 19 | Fn   | WaitAndGetTxData                | WaitAndGetTxData is a function to get transaction data after transaction is processed                                                               |
| 20 | Fn   | WaitAndGetTxData                | WaitAndGetTxData is a function to get transaction data after transaction is processed                                                               |
| 21 | Fn   | WaitForNextBlock                | WaitForNextBlock is a function to wait until next block                                                                                             |

## Handlers struct package
github.com/Pylons-tech/pylons_sdk/x/pylons/handlers

### Structs

| No | Type   | Name                        | Description                                                                                 |
|----|--------|-----------------------------|---------------------------------------------------------------------------------------------|
| 1  | Struct | CheckExecutionResponse      | CheckExecutionResponse is the response for checkExecution                                   |
| 2  | Struct | CreateCookbookResponse      | CreateCookbookResponse is a struct of create cookbook response                              |
| 3  | Struct | CreateRecipeResponse        | CreateRecipeResponse is struct of create recipe response                                    |
| 4  | Struct | ExecuteRecipeResponse       | ExecuteRecipeResponse is the response for executeRecipe                                     |
| 5  | Struct | ExecuteRecipeScheduleOutput | ExecuteRecipeScheduleOutput is a struct that shows how execute recipe schedule output works |
| 6  | Struct | FiatItemResponse            | FiatItemResponse is a struct to control fiat item response                                  |
| 7  | Struct | FulfillTradeResponse        | FulfillTradeResponse is a struct to control fulfill trade response                          |
| 8  | Struct | GetParamsForPopularRecipe   | GetParamsForPopularRecipe is a function to get popular recipe's attributes                  |
| 9  | Struct | PopularRecipeType           | PopularRecipeType is a type for popular recipes                                             |
| 10 | Struct | UpdateItemStringResponse    | UpdateItemStringResponse is a struct to control update item string response                 |

## Msgs package
github.com/Pylons-tech/pylons_sdk/x/pylons/msgs

### Msg structs

| No | Type   | Name                | Description                                                                              |
|----|--------|---------------------|------------------------------------------------------------------------------------------|
| 1  | Struct | MsgCheckExecution   | MsgCheckExecution defines a CheckExecution message                                       |
| 2  | Struct | MsgCreateCookbook   | MsgCreateCookbook defines a CreateCookbook message                                       |
| 3  | Struct | MsgCreateRecipe     | NewMsgCreateRecipe a constructor for CreateRecipe msg                                    |
| 4  | Struct | MsgCreateTrade      | MsgCreateTrade defines a CreateTrade message                                             |
| 5  | Struct | MsgDisableRecipe    | MsgDisableRecipe defines a DisableRecipe message                                         |
| 6  | Struct | MsgDisableTrade     | MsgDisableTrade defines a DisableTrade message                                           |
| 7  | Struct | MsgEnableRecipe     | MsgEnableRecipe defines a EnableRecipe message                                           |
| 8  | Struct | MsgEnableTrade      | MsgEnableTrade defines a EnableTrade message                                             |
| 9  | Struct | MsgExecuteRecipe    | MsgExecuteRecipe defines a SetName message                                               |
| 10 | Struct | MsgFiatItem         | MsgFiatItem is a msg struct to be used to fiat item                                      |
| 11 | Struct | MsgFulfillTrade     | NewMsgFulfillTrade a constructor for FulfillTrade msg                                    |
| 12 | Struct | MsgGetPylons        | MsgGetPylons defines a GetPylons message                                                 |
| 13 | Struct | MsgSendPylons       | MsgSendPylons defines a SendPylons message                                               |
| 14 | Struct | MsgUpdateCookbook   | MsgUpdateCookbook defines a UpdateCookbook message                                       |
| 15 | Struct | MsgUpdateItemString | MsgUpdateItemString defines a UpdateItemString message                                   |
| 16 | Struct | MsgUpdateRecipe     | MsgUpdateRecipe defines a UpdateRecipe message                                           |

### Msg Utility functions 

| No | Type     | Name                   | Description                                                                              |
|----|----------|------------------------|------------------------------------------------------------------------------------------|
| 1  | Constant | DefaultCostPerBlock    | DefaultCostPerBlock the amount of pylons to be charged by default                        |
| 2  | Fn       | NewMsgCheckExecution   | NewMsgCheckExecution a constructor for CheckExecution msg                                |
| 3  | Fn       | NewMsgCreateCookbook   | NewMsgCreateCookbook a constructor for CreateCookbook msg                                |
| 4  | Fn       | NewMsgCreateRecipe     | NewMsgCreateRecipe a constructor for CreateRecipe msg                                    |
| 5  | Fn       | NewMsgCreateTrade      | NewMsgCreateTrade a constructor for CreateTrade msg                                      |
| 6  | Fn       | NewMsgDisableTrade     | NewMsgDisableTrade a constructor for DisableTrade msg                                    |
| 7  | Fn       | NewMsgExecuteRecipe    | NewMsgExecuteRecipe a constructor for ExecuteCookbook msg                                |
| 8  | Fn       | NewMsgFiatItem         | NewMsgFiatItem a constructor for MsgFiatItem msg                                         |
| 9  | Fn       | NewMsgFulfillTrade     | NewMsgFulfillTrade a constructor for FulfillTrade msg                                    |
| 10 | Fn       | NewMsgGetPylons        | NewMsgGetPylons is a function to get MsgGetPylons msg from required params               |
| 11 | Fn       | NewMsgUpdateItemString | NewMsgUpdateItemString is a function to get MsgUpdateItemString msg from required params |

### Type structs

| No | Type     | Name                 | Description                                                                                                                                                     |
|----|----------|----------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1  | Constant | BasicTier            | BasicTier is the cookbook tier which doesn't allow paid receipes which means the developers cannot have receipes where they can actually carge a fee in pylons. |
| 2  | Struct   | CoinInputList        | CoinInputList is a list of Coin inputs                                                                                                                          |
| 3  | Struct   | CoinOutput           | CoinOutput is the game elements that are needs as output to a recipe                                                                                            |
| 4  | Struct   | Cookbook             | Cookbook is a struct that contains all the metadata of a cookbook                                                                                               |
| 5  | Struct   | DoubleInputParam     | DoubleInputParam describes the bounds on an item input/output parameter of type float64                                                                         |
| 6  | Struct   | DoubleInputParamList | DoubleInputParamList is a list of DoubleInputParam                                                                                                              |
| 7  | Struct   | EntriesList          | EntriesList is a struct to keep list of items and coins                                                                                                         |
| 8  | Struct   | Execution            | Execution is a recipe execution used for tracking the execution #### specifically a scheduled execution                                                         |
| 9  | Struct   | FloatString          | FloatString is a wrapper to resolve the amino issues                                                                                                            |
| 10 | Struct   | Item                 | Item is a tradable asset                                                                                                                                        |
| 11 | Struct   | ItemInput            | ItemInput is a wrapper struct for Item for recipes                                                                                                              |
| 12 | Struct   | ItemInputList        | ItemInputList is a list of ItemInputs for convinience                                                                                                           |
| 13 | Struct   | ItemList             | ItemList is a list of items                                                                                                                                     |
| 14 | Struct   | ItemModifyParams     | ItemModifyParams describes the fields that needs to be modified                                                                                                 |
| 15 | Struct   | LongInputParam       | LongInputParam describes the bounds on an item input/output parameter of type int64                                                                             |
| 16 | Struct   | LongInputParamList   | LongInputParamList is a list of LongInputParam                                                                                                                  |
| 17 | Struct   | PremiumTier          | PremiumTier the cookbook tier which does allow paid receipes                                                                                                    |
| 18 | Struct   | Recipe               | Recipe is a game state machine step abstracted out as a cooking terminology                                                                                     |
| 19 | Struct   | StringInputParam     | StringInputParam describes the bounds on an item input/output parameter of type string                                                                          |
| 20 | Struct   | StringInputParamList | StringInputParamList is a list of StringInputParam                                                                                                              |
| 21 | Struct   | Trade                | Trade is a construct to perform exchange of items and coins between users. Initiated by the sender and completed by the FulFiller.                              |
| 22 | Struct   | TradeList            | TradeList is a list of trades                                                                                                                                   |
| 23 | Struct   | WeightedOutputsList  | WeightedOutputsList is a struct to keep items which can be generated by weight                                                                                  |

### Type Utility functions 

| No | Type | Name                           | Description                                                                                    |
|----|------|--------------------------------|------------------------------------------------------------------------------------------------|
| 1  | Fn   | GenCoinInputList               | GenCoinInputList is a utility function to genearte coin input list                             |
| 2  | Fn   | GenEntries                     | GenEntries is a function to generate entries from coin name and item name                      |
| 3  | Fn   | GenEntriesFirstItemNameUpgrade | GenEntriesFirstItemNameUpgrade is a function to generate entries that update first item's name |
| 4  | Fn   | GenItemInputList               | GenItemInputList is a utility function to genearte item input list                             |
| 5  | Fn   | GenItemOnlyEntry               | GenItemOnlyEntry is a utility function to generate item only entry                             |
| 6  | Fn   | GenOneOutput                   | GenOneOutput is a function to generate output with one from entry list                         |
| 7  | Fn   | NewPylon                       | NewPylon Returns pylon currency                                                                |

## App package
github.com/Pylons-tech/pylons_sdk/app

#### MakeCodec  
MakeCodec make codec for message marshal/unmarshal

## Queriers package
github.com/Pylons-tech/pylons_sdk/x/pylons/queriers

#### ExecResponse  
ExecResponse is the response for ListExecutions

#### ItemResponse  
ItemResponse is the response for Items

