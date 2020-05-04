package intTest

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"

	"github.com/Pylons-tech/pylons_sdk/x/pylons/queriers"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ListTradeViaCLI(account string) ([]types.Trade, error) {
	queryParams := []string{"query", "pylons", "list_trade"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, err := RunPylonsCli(queryParams, "")
	if err != nil {
		return []types.Trade{}, err
	}
	listTradesResp := types.TradeList{}
	err = GetAminoCdc().UnmarshalJSON(output, &listTradesResp)
	if err != nil {
		return []types.Trade{}, err
	}
	return listTradesResp.Trades, nil
}

func GetTradeIDFromExtraInfo(tradeExtraInfo string) (string, bool, error) {
	trdList, err := ListTradeViaCLI("")
	if err != nil {
		return "", false, err
	}
	trade, exist := FindTradeFromArrayByExtraInfo(trdList, tradeExtraInfo)
	return trade.ID, exist, nil
}

func ListCookbookViaCLI(account string) ([]types.Cookbook, error) {
	queryParams := []string{"query", "pylons", "list_cookbook"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, err := RunPylonsCli(queryParams, "")
	if err != nil {
		return []types.Cookbook{}, err
	}
	listCBResp := types.CookbookList{}
	err = GetAminoCdc().UnmarshalJSON(output, &listCBResp)
	if err != nil {
		return []types.Cookbook{}, err
	}
	return listCBResp.Cookbooks, err
}

func ListRecipesViaCLI(account string) ([]types.Recipe, error) {
	queryParams := []string{"query", "pylons", "list_recipe"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, err := RunPylonsCli(queryParams, "")
	if err != nil {
		return []types.Recipe{types.Recipe{}}, err
	}
	listRCPResp := types.RecipeList{}
	err = GetAminoCdc().UnmarshalJSON(output, &listRCPResp)
	if err != nil {
		return []types.Recipe{types.Recipe{}}, err
	}
	return listRCPResp.Recipes, err
}

func ListExecutionsViaCLI(account string, t *testing.T) ([]types.Execution, error) {
	queryParams := []string{"query", "pylons", "list_executions"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, err := RunPylonsCli(queryParams, "")
	if err != nil {
		t.Fatalf("error running list_executions cli command ::: %+v", err)
		return []types.Execution{}, err
	}
	var listExecutionsResp queriers.ExecResp
	err = GetAminoCdc().UnmarshalJSON(output, &listExecutionsResp)
	if err != nil {
		t.Fatalf("error unmarshaling list executions ::: %+v", err)
		return []types.Execution{}, err
	}
	return listExecutionsResp.Executions, err
}

func ListItemsViaCLI(account string) ([]types.Item, error) {
	queryParams := []string{"query", "pylons", "items_by_sender"}
	if len(account) != 0 {
		queryParams = append(queryParams, "--account", account)
	}
	output, err := RunPylonsCli(queryParams, "")
	if err != nil {
		return types.ItemList{}, err
	}
	var itemResp queriers.ItemResp
	err = GetAminoCdc().UnmarshalJSON(output, &itemResp)
	if err != nil {
		return types.ItemList{}, err
	}
	return itemResp.Items, err
}

func GetRawTxResponse(txhash string, t *testing.T) (sdk.TxResponse, error) {
	var tx sdk.TxResponse
	output, err := RunPylonsCli([]string{"query", "tx", txhash}, "")
	if err != nil {
		return tx, err
	}
	err = GetAminoCdc().UnmarshalJSON([]byte(output), &tx)
	return tx, err
}

func GetTxErrorFromRawTxResponse(tx sdk.TxResponse, t *testing.T) []byte {
	if len(tx.Logs) > 0 {
		return []byte(tx.Logs[0].Log)
	}
	return []byte{}
}

func GetTxErrorFromTxHash(txhash string, t *testing.T) ([]byte, error) {
	tx, err := GetRawTxResponse(txhash, t)
	if err != nil {
		return []byte{}, err
	}
	return GetTxErrorFromRawTxResponse(tx, t), nil
}

func GetHumanReadableErrorFromTxResponse(tx sdk.TxResponse, t *testing.T) (string, error) {
	txErrorBytes := GetTxErrorFromRawTxResponse(tx, t)
	if len(txErrorBytes) == 0 {
		return "", nil
	}
	hmrErr := struct {
		Codespace string `json:"codespace"`
		Code      int    `json:"code"`
		Message   string `json:"message"`
	}{}
	err := json.Unmarshal(txErrorBytes, &hmrErr)
	return hmrErr.Message, err
}

func GetHumanReadableErrorFromTxHash(txhash string, t *testing.T) (string, error) {
	tx, err := GetRawTxResponse(txhash, t)
	if err != nil {
		return "", err
	}
	return GetHumanReadableErrorFromTxResponse(tx, t)
}

func GetTxDataFromRawTxResponse(tx sdk.TxResponse, t *testing.T) ([]byte, error) {
	bs, err := hex.DecodeString(tx.Data)
	if err != nil {
		return []byte{}, err
	}
	return bs, nil
}

func GetTxDataFromTxHash(txhash string, t *testing.T) ([]byte, error) {
	tx, err := GetRawTxResponse(txhash, t)
	if err != nil {
		return []byte{}, err
	}
	return GetTxDataFromRawTxResponse(tx, t)
}

func WaitAndGetTxData(txhash string, maximum_wait_block int64, t *testing.T) ([]byte, error) {
	var tx sdk.TxResponse
	var txHandleResBytes []byte
	var err error

	tx, err = GetRawTxResponse(txhash, t)
	if err == nil {
		txHandleResBytes, err = GetTxDataFromRawTxResponse(tx, t)
	}
	if err != nil { // maybe transaction is not contained in block
		if maximum_wait_block == 0 {
			hmrError, err := GetHumanReadableErrorFromTxResponse(tx, t)
			var ltl TxLog
			if len(TxLogs) > 0 {
				ltl = TxLogs[len(TxLogs)-1]
			}
			return txHandleResBytes, errors.New(
				fmt.Sprintf("didn't get result waiting for maximum_wait_block.\nhmrError=%s, err=%+v\nLast TxLog=%+v",
					hmrError, err, ltl.String()))
		} else {
			WaitForNextBlock()
			return WaitAndGetTxData(txhash, maximum_wait_block-1, t)
		}
	}
	return txHandleResBytes, nil
}

func FindTradeFromArrayByExtraInfo(trades []types.Trade, extraInfo string) (types.Trade, bool) {
	for _, trade := range trades {
		if trade.ExtraInfo == extraInfo {
			return trade, true
		}
	}
	return types.Trade{}, false
}

func FindCookbookFromArrayByName(cbList []types.Cookbook, name string) (types.Cookbook, bool) {
	for _, cb := range cbList {
		if cb.Name == name {
			return cb, true
		}
	}
	return types.Cookbook{}, false
}

func FindRecipeFromArrayByName(recipes []types.Recipe, name string) (types.Recipe, bool) {
	for _, rcp := range recipes {
		if rcp.Name == name {
			return rcp, true
		}
	}
	return types.Recipe{}, false
}

func FindExecutionByRecipeID(execs []types.Execution, rcpID string) (types.Execution, bool) {
	for _, exec := range execs {
		if exec.RecipeID == rcpID {
			return exec, true
		}
	}
	return types.Execution{}, false
}

func FindItemFromArrayByName(items []types.Item, name string, includeLockedByRcp bool) (types.Item, bool) {
	for _, item := range items {
		itemName, _ := item.FindString("Name")
		if !includeLockedByRcp && len(item.OwnerRecipeID) != 0 {
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
	output, err := RunPylonsCli([]string{"query", "pylons", "get_cookbook", guid}, "")
	if err != nil {
		return types.Cookbook{}, err
	}
	var cookbook types.Cookbook
	err = GetAminoCdc().UnmarshalJSON(output, &cookbook)
	if err != nil {
		return types.Cookbook{}, err
	}
	return cookbook, err
}

func GetCookbookIDFromName(cbName string, account string) (string, bool, error) {
	cbList, err := ListCookbookViaCLI(account)
	if err != nil {
		return "", false, err
	}

	cb, exist := FindCookbookFromArrayByName(cbList, cbName)
	return cb.ID, exist, nil
}

func GetRecipeIDFromName(rcpName string) (string, bool, error) {
	rcpList, err := ListRecipesViaCLI("")
	if err != nil {
		return "", false, err
	}
	rcp, exist := FindRecipeFromArrayByName(rcpList, rcpName)
	return rcp.ID, exist, nil
}

func GetItemIDFromName(itemName string, includeLockedByRcp bool) (string, bool, error) {
	itemList, err := ListItemsViaCLI("")
	if err != nil {
		return "", false, err
	}
	rcp, exist := FindItemFromArrayByName(itemList, itemName, includeLockedByRcp)
	return rcp.ID, exist, nil
}

// GetRecipeByGUID is to get Recipe from ID
func GetRecipeByGUID(guid string) (types.Recipe, error) {
	output, err := RunPylonsCli([]string{"query", "pylons", "get_recipe", guid}, "")
	if err != nil {
		return types.Recipe{}, err
	}
	var rcp types.Recipe
	err = GetAminoCdc().UnmarshalJSON(output, &rcp)
	if err != nil {
		return types.Recipe{}, err
	}
	return rcp, err
}

// GetExecutionByGUID is to get Execution from ID
func GetExecutionByGUID(guid string) (types.Execution, error) {
	output, err := RunPylonsCli([]string{"query", "pylons", "get_execution", guid}, "")
	if err != nil {
		return types.Execution{}, err
	}
	var exec types.Execution
	err = GetAminoCdc().UnmarshalJSON(output, &exec)
	if err != nil {
		return types.Execution{}, err
	}
	return exec, err
}

// GetItemByGUID is to get Item from ID
func GetItemByGUID(guid string) (types.Item, error) {
	output, err := RunPylonsCli([]string{"query", "pylons", "get_item", guid}, "")
	if err != nil {
		return types.Item{}, err
	}
	var item types.Item
	err = GetAminoCdc().UnmarshalJSON(output, &item)
	if err != nil {
		return types.Item{}, err
	}
	return item, err
}
