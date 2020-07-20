package inttest

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/evtesting"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/queriers"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ListTradeViaCLI is a function to get list of trades from cli
func ListTradeViaCLI(account string) ([]types.Trade, error) {
	queryParams := []string{"query", "pylons", "list_trade"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, logstr, err := RunPylonsCli(queryParams, "")
	if err != nil {
		return []types.Trade{}, fmt.Errorf("%s: %s", logstr, err.Error())
	}
	listTradesResp := types.TradeList{}
	err = GetAminoCdc().UnmarshalJSON(output, &listTradesResp)
	return listTradesResp.Trades, err
}

// GetTradeIDFromExtraInfo is a function to get trade id from trade extra info
func GetTradeIDFromExtraInfo(tradeExtraInfo string) (string, bool, error) {
	trdList, err := ListTradeViaCLI("")
	if err != nil {
		return "", false, err
	}
	trade, exist := FindTradeFromArrayByExtraInfo(trdList, tradeExtraInfo)
	return trade.ID, exist, nil
}

// ListCookbookViaCLI is a function to list cookbooks via cli
func ListCookbookViaCLI(account string) ([]types.Cookbook, error) {
	listCBResp := types.CookbookList{}
	queryParams := []string{"query", "pylons", "list_cookbook"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, logstr, err := RunPylonsCli(queryParams, "")
	if err != nil {
		return listCBResp.Cookbooks, fmt.Errorf("%s: %s", logstr, err.Error())
	}
	err = GetAminoCdc().UnmarshalJSON(output, &listCBResp)
	return listCBResp.Cookbooks, err
}

// ListLockedCoinsViaCLI is a function to list locked coins via cli
func ListLockedCoinsViaCLI(account string) (types.LockedCoin, error) {
	lcResp := types.LockedCoin{}
	queryParams := []string{"query", "pylons", "list_locked_coins"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, logstr, err := RunPylonsCli(queryParams, "")
	if err != nil {
		return lcResp, fmt.Errorf("%s: %s", logstr, err.Error())
	}
	err = GetAminoCdc().UnmarshalJSON(output, &lcResp)
	return lcResp, err
}

// ListRecipesViaCLI is a function to list recipes via cli
func ListRecipesViaCLI(account string) ([]types.Recipe, error) {
	queryParams := []string{"query", "pylons", "list_recipe"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, _, err := RunPylonsCli(queryParams, "")
	if err != nil {
		return []types.Recipe{}, err
	}
	listRCPResp := types.RecipeList{}
	err = GetAminoCdc().UnmarshalJSON(output, &listRCPResp)
	return listRCPResp.Recipes, err
}

// ListExecutionsViaCLI is a function to list executions via cli
func ListExecutionsViaCLI(account string, t *testing.T) ([]types.Execution, error) {
	queryParams := []string{"query", "pylons", "list_executions"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, _, err := RunPylonsCli(queryParams, "")
	if err != nil {
		t.MustNil(err, "error running list_executions cli command")
		return []types.Execution{}, err
	}
	var listExecutionsResp queriers.ExecResponse
	err = GetAminoCdc().UnmarshalJSON(output, &listExecutionsResp)
	t.WithFields(testing.Fields{
		"list_executions_output": string(output),
	}).MustNil(err, "error unmarshaling list executions")
	return listExecutionsResp.Executions, err
}

// ListItemsViaCLI is a function to list items via cli
func ListItemsViaCLI(account string) ([]types.Item, error) {
	queryParams := []string{"query", "pylons", "items_by_sender"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, logstr, err := RunPylonsCli(queryParams, "")
	if err != nil {
		return types.ItemList{}, fmt.Errorf("%s: %s", logstr, err.Error())
	}
	var ItemResponse queriers.ItemResponse
	err = GetAminoCdc().UnmarshalJSON(output, &ItemResponse)
	return ItemResponse.Items, err
}

// WaitAndGetTxError is a function to wait and get transaction error from hash
func WaitAndGetTxError(txhash string, maxWaitBlock int64, t *testing.T) ([]byte, error) {
	txErrorResBytes, err := GetTxError(txhash, t)
	if err != nil { // maybe transaction is not contained in block
		if maxWaitBlock == 0 {
			t.WithFields(testing.Fields{
				"output": string(txErrorResBytes),
				"error":  err,
			}).Error("didn't get result waiting for maximum wait block")
			return txErrorResBytes, errors.New("didn't get result waiting for maximum wait block")
		}
		if err = WaitForNextBlock(); err != nil {
			return txErrorResBytes, err
		}
		return WaitAndGetTxError(txhash, maxWaitBlock-1, t)
	}
	return txErrorResBytes, nil
}

// IsJSON checks if bytes is in json
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// GetTxError is a function to get transaction error from txhash
func GetTxError(txhash string, t *testing.T) ([]byte, error) {
	output, logstr, err := RunPylonsCli([]string{"query", "tx", txhash}, "")
	if err != nil {
		return []byte{}, fmt.Errorf("%s: %s", logstr, err.Error())
	}
	var tx sdk.TxResponse
	err = GetAminoCdc().UnmarshalJSON([]byte(output), &tx)
	if err != nil {
		return []byte{}, err
	}

	if !IsJSON(tx.RawLog) {
		return []byte(tx.RawLog), nil
	}
	return []byte{}, nil
}

// GetHumanReadableErrorFromTxHash is a function to get human readable error from txhash
func GetHumanReadableErrorFromTxHash(txhash string, t *testing.T) string {
	txErrorBytes, err := WaitAndGetTxError(txhash, GetMaxWaitBlock(), t)
	t.WithFields(testing.Fields{
		"tx_error_bytes": string(txErrorBytes),
	}).MustNil(err, "error while waiting for a transaction")
	return string(txErrorBytes)
}

// GetTxData is a function to get transaction result data by txhash
func GetTxData(txhash string, t *testing.T) ([]byte, error) {
	output, _, err := RunPylonsCli([]string{"query", "tx", txhash}, "")
	if err != nil {
		t.WithFields(testing.Fields{
			"output": string(output),
			"error":  err,
		}).Debug("query for tx") // do debug as in this step, transaction could be in mempool
		return output, err
	}
	var tx sdk.TxResponse
	err = GetAminoCdc().UnmarshalJSON([]byte(output), &tx)
	if err != nil {
		return []byte{}, err
	}
	bs, err := hex.DecodeString(tx.Data)
	return bs, err
}

// WaitAndGetTxData is a function to get transaction data after transaction is processed
func WaitAndGetTxData(txhash string, maxWaitBlock int64, t *testing.T) ([]byte, error) {
	txHandleResBytes, err := GetTxData(txhash, t)
	if err != nil { // maybe transaction is not contained in block
		if maxWaitBlock == 0 {
			t.WithFields(testing.Fields{
				"action": "func_end",
			}).Error("didn't get result waiting for maximum wait block")
			return txHandleResBytes, errors.New("didn't get result waiting for maximum wait block")
		}
		if err = WaitForNextBlock(); err != nil {
			return txHandleResBytes, err
		}
		return WaitAndGetTxData(txhash, maxWaitBlock-1, t)
	}
	return txHandleResBytes, nil
}

// FindTradeFromArrayByExtraInfo is a function to find trade from extra info
func FindTradeFromArrayByExtraInfo(trades []types.Trade, extraInfo string) (types.Trade, bool) {
	for _, trade := range trades {
		if trade.ExtraInfo == extraInfo {
			return trade, true
		}
	}
	return types.Trade{}, false
}

// FindCookbookFromArrayByName is a function to cookbook from array by name
func FindCookbookFromArrayByName(cbList []types.Cookbook, name string) (types.Cookbook, bool) {
	for _, cb := range cbList {
		if cb.Name == name {
			return cb, true
		}
	}
	return types.Cookbook{}, false
}

// FindRecipeFromArrayByName is a function to get recipes from array by name
func FindRecipeFromArrayByName(recipes []types.Recipe, name string) (types.Recipe, bool) {
	for _, rcp := range recipes {
		if rcp.Name == name {
			return rcp, true
		}
	}
	return types.Recipe{}, false
}

// FindExecutionByRecipeID is a function to get execution by recipe id
func FindExecutionByRecipeID(execs []types.Execution, rcpID string) (types.Execution, bool) {
	for _, exec := range execs {
		if exec.RecipeID == rcpID {
			return exec, true
		}
	}
	return types.Execution{}, false
}

// FindItemFromArrayByName is a function to get item array by name
func FindItemFromArrayByName(
	items []types.Item, name string,
	includeLockedByRecipe bool,
	includeLockedByTrade bool,
) (types.Item, bool) {
	for _, item := range items {
		itemName, _ := item.FindString("Name")
		if !includeLockedByRecipe && len(item.OwnerRecipeID) != 0 {
			continue
		}
		if !includeLockedByTrade && len(item.OwnerTradeID) != 0 {
			continue
		}
		if itemName == name {
			return item, true
		}
	}
	return types.Item{}, false
}

// GetCookbookByGUID is to get Cookbook from ID
func GetCookbookByGUID(guid string) (types.Cookbook, error) {
	output, _, err := RunPylonsCli([]string{"query", "pylons", "get_cookbook", guid}, "")
	if err != nil {
		return types.Cookbook{}, err
	}
	var cookbook types.Cookbook
	err = GetAminoCdc().UnmarshalJSON(output, &cookbook)
	return cookbook, err
}

// GetCookbookIDFromName is a function to get cookbook id from name
func GetCookbookIDFromName(cbName string, account string) (string, bool, error) {
	cbList, err := ListCookbookViaCLI(account)
	if err != nil {
		return "", false, err
	}

	cb, exist := FindCookbookFromArrayByName(cbList, cbName)
	return cb.ID, exist, nil
}

// GetRecipeIDFromName is a function to get recipe id from name
func GetRecipeIDFromName(rcpName string) (string, bool, error) {
	rcpList, err := ListRecipesViaCLI("")
	if err != nil {
		return "", false, err
	}
	rcp, exist := FindRecipeFromArrayByName(rcpList, rcpName)
	return rcp.ID, exist, nil
}

// GetItemIDFromName is a function to get item id from name
func GetItemIDFromName(itemName string, includeLockedByRecipe bool, includeLockedByTrade bool) (string, bool, error) {
	itemList, err := ListItemsViaCLI("")
	if err != nil {
		return "", false, err
	}
	rcp, exist := FindItemFromArrayByName(itemList, itemName, includeLockedByRecipe, includeLockedByTrade)
	return rcp.ID, exist, nil
}

// GetRecipeByGUID is to get Recipe from ID
func GetRecipeByGUID(guid string) (types.Recipe, error) {
	output, _, err := RunPylonsCli([]string{"query", "pylons", "get_recipe", guid}, "")
	if err != nil {
		return types.Recipe{}, err
	}
	var rcp types.Recipe
	err = GetAminoCdc().UnmarshalJSON(output, &rcp)
	return rcp, err
}

// GetExecutionByGUID is to get Execution from ID
func GetExecutionByGUID(guid string) (types.Execution, error) {
	output, _, err := RunPylonsCli([]string{"query", "pylons", "get_execution", guid}, "")
	if err != nil {
		return types.Execution{}, err
	}
	var exec types.Execution
	err = GetAminoCdc().UnmarshalJSON(output, &exec)
	return exec, err
}

// GetItemByGUID is to get Item from ID
func GetItemByGUID(guid string) (types.Item, error) {
	output, _, err := RunPylonsCli([]string{"query", "pylons", "get_item", guid}, "")
	if err != nil {
		return types.Item{}, err
	}
	var item types.Item
	err = GetAminoCdc().UnmarshalJSON(output, &item)
	return item, err
}

// GetRecipeGUIDFromName is a function to get recipe id from name
func GetRecipeGUIDFromName(name string, account string) (string, error) {
	rcpList, err := ListRecipesViaCLI(account)
	if err != nil {
		return "", err
	}
	rcp, _ := FindRecipeFromArrayByName(rcpList, name)
	return rcp.ID, nil
}
