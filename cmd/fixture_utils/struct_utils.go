package fixturetest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/evtesting"
	inttest "github.com/Pylons-tech/pylons_sdk/cmd/test_utils"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
)

var execIDRWMutex sync.Mutex
var execIDs = make(map[string]string)

// ReadFile is a function to read file
func ReadFile(fileURL string, t *testing.T) []byte {
	jsonFile, err := os.Open(path.Join(FixtureTestOpts.BaseDirectory, fileURL))
	t.MustNil(err, "fatal log reading file")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

// UnmarshalIntoEmptyInterface is a function to convert bytes into map[string]interface{}
func UnmarshalIntoEmptyInterface(bytes []byte, t *testing.T) map[string]interface{} {
	var raw map[string]interface{}
	err := json.Unmarshal(bytes, &raw)
	t.WithFields(testing.Fields{
		"bytes": string(bytes),
	}).MustNil(err, "error unmarshaling")
	return raw
}

// RegisterDefaultAccountKeys register the accounts configured on FixtureTestOpts.AccountNames
func RegisterDefaultAccountKeys() {
	runtimeKeyGenMux.Lock()
	defer runtimeKeyGenMux.Unlock()
	for idx, key := range FixtureTestOpts.AccountNames {
		runtimeAccountKeys[fmt.Sprintf("account%d", idx+1)] = key
	}
}

// GetAccountKeyFromTempName is a function to get account key from temp name
func GetAccountKeyFromTempName(tempName string, t *testing.T) string {
	t.MustTrue(len(tempName) > 0, "account key should not be an empty string")
	runtimeKeyGenMux.Lock()
	defer runtimeKeyGenMux.Unlock()
	key, ok := runtimeAccountKeys[tempName]
	if !ok {
		key = fmt.Sprintf("FixtureRuntime_%s_%d", tempName, time.Now().Unix())
		runtimeAccountKeys[tempName] = key
	}
	return key
}

// GetAccountAddressFromTempName is a function to get account address from temp name
func GetAccountAddressFromTempName(tempName string, t *testing.T) string {
	accountKey := GetAccountKeyFromTempName(tempName, t)
	return inttest.GetAccountAddr(accountKey, t)
}

// GetSenderKeyFromRef is a function to get create cookbook message from reference
func GetSenderKeyFromRef(ref string, t *testing.T) string {
	bytes := ReadFile(ref, t)
	raw := UnmarshalIntoEmptyInterface(bytes, t)
	sender, ok := raw["Sender"].(string)
	t.MustTrue(ok, "sender field is empty")
	return sender
}

// UpdateSenderKeyToAddress is a function to update sender key to sender's address
func UpdateSenderKeyToAddress(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	senderTempName, ok := raw["Sender"].(string)
	t.MustTrue(ok, "sender field is empty")

	raw["Sender"] = GetAccountAddressFromTempName(senderTempName, t)
	newBytes, err := json.Marshal(raw)
	t.WithFields(testing.Fields{
		"updated_sender_interface": raw,
	}).MustNil(err, "error encoding raw json")
	return newBytes
}

// UpdateReceiverKeyToAddress is a function to update receiver key to receiver's address
func UpdateReceiverKeyToAddress(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	receiverTempName, ok := raw["Receiver"].(string)
	t.MustTrue(ok, "receiver field is empty")
	raw["Receiver"] = GetAccountAddressFromTempName(receiverTempName, t)
	newBytes, err := json.Marshal(raw)
	t.WithFields(testing.Fields{
		"updated_receiver_interface": raw,
	}).MustNil(err, "error encoding raw json")
	return newBytes
}

// UpdateRequesterKeyToAddress is a function to update requester key to requester's address
func UpdateRequesterKeyToAddress(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	requesterTempName, ok := raw["Requester"].(string)
	t.MustTrue(ok, "requester field is empty")
	raw["Requester"] = GetAccountAddressFromTempName(requesterTempName, t)
	newBytes, err := json.Marshal(raw)
	t.WithFields(testing.Fields{
		"updated_requester_interface": raw,
	}).MustNil(err, "error encoding raw json")
	return newBytes
}

// UpdateCBNameToID is a function to update cookbook name to cookbook id if it has cookbook name field
func UpdateCBNameToID(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	cbName, ok := raw["CookbookName"].(string)
	if !ok {
		return bytes
	}
	cbID, exist, err := inttest.GetCookbookIDFromName(cbName, "")
	if exist && err != nil {
		raw["CookbookID"] = cbID
		newBytes, err := json.Marshal(raw)
		t.WithFields(testing.Fields{
			"updated_cookbook_id_interface": raw,
		}).MustNil(err, "error encoding raw json")
		return newBytes
	}
	return bytes
}

// UpdateRecipeName is a function to update recipe name into recipe id
func UpdateRecipeName(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	rcpName, ok := raw["RecipeName"].(string)
	t.MustTrue(ok, "recipe name field is empty")
	rcpID, exist, err := inttest.GetRecipeIDFromName(rcpName)
	t.WithFields(testing.Fields{
		"recipe_name": rcpName,
	}).MustTrue(exist, "there's no recipe id with specific recipe name")
	t.MustNil(err, "error getting recipe id from name")
	raw["RecipeID"] = rcpID
	newBytes, err := json.Marshal(raw)
	t.WithFields(testing.Fields{
		"updated_recipe_id_interface": raw,
	}).MustNil(err, "error encoding raw json")
	return newBytes
}

// UpdateTradeExtraInfoToID is a function to update trade extra info into trade id
func UpdateTradeExtraInfoToID(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	trdInfo, ok := raw["TradeInfo"].(string)
	t.MustTrue(ok, "trade info does not exist in json")
	trdID, exist, err := inttest.GetTradeIDFromExtraInfo(trdInfo)
	t.WithFields(testing.Fields{
		"trade_info": trdInfo,
	}).MustTrue(exist, "there's not trade id with specific info")
	t.MustNil(err, "error getting trade id from info")
	raw["TradeID"] = trdID
	newBytes, err := json.Marshal(raw)
	t.WithFields(testing.Fields{
		"updated_trade_id_interface": raw,
	}).MustNil(err, "error encoding raw json")
	return newBytes
}

// UpdateExecID is a function to set execute id from execID reference
func UpdateExecID(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	var execRefReader struct {
		ExecRef string
	}
	err := json.Unmarshal(bytes, &execRefReader)
	t.WithFields(testing.Fields{
		"bytes": string(bytes),
	}).MustNil(err, "error unmarshaling into exec ref")

	var ok bool
	execIDRWMutex.Lock()
	raw["ExecID"], ok = execIDs[execRefReader.ExecRef]
	execIDRWMutex.Unlock()
	t.WithFields(testing.Fields{
		"execRef": execRefReader.ExecRef,
	}).MustTrue(ok, "execID not available")
	newBytes, err := json.Marshal(raw)
	t.WithFields(testing.Fields{
		"updated_exec_id_interface": raw,
	}).MustNil(err, "error encoding raw json")
	return newBytes
}

// UpdateItemIDFromName is a function to set item id from item name
func UpdateItemIDFromName(bytes []byte, includeLockedByRecipe, includeLockedByTrade bool, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	itemName, ok := raw["ItemName"].(string)
	t.MustTrue(ok, "item name does not exist in json")
	sender, ok := raw["Sender"].(string)
	t.MustTrue(ok, "sender address does not exist in json")

	itemID, exist, err := inttest.GetItemIDFromName(sender, itemName, includeLockedByRecipe, includeLockedByTrade)
	if !exist {
		t.WithFields(testing.Fields{
			"item_name":                itemName,
			"include_locked_by_recipe": includeLockedByRecipe,
			"include_locked_by_trade":  includeLockedByTrade,
		}).Debug("no item fit params")
	}
	t.WithFields(testing.Fields{
		"item_name":                itemName,
		"include_locked_by_recipe": includeLockedByRecipe,
		"include_locked_by_trade":  includeLockedByTrade,
	}).MustNil(err, "error getting item id from name")
	raw["ItemID"] = itemID
	newBytes, err := json.Marshal(raw)
	t.WithFields(testing.Fields{
		"updated_item_id_interface": raw,
	}).MustNil(err, "error encoding raw json")
	return newBytes
}

// GetItemIDsFromNames is a function to set item ids from names for recipe execution
func GetItemIDsFromNames(bytes []byte, sender string, includeLockedByRecipe, includeLockedByTrade bool, t *testing.T) []string {
	var itemNamesResp struct {
		ItemNames []string
	}
	err := json.Unmarshal(bytes, &itemNamesResp)
	t.WithFields(testing.Fields{
		"bytes": string(bytes),
	}).MustNil(err, "error unmarshaling into item names")
	ItemIDs := []string{}

	for _, itemName := range itemNamesResp.ItemNames {
		itemID, exist, err := inttest.GetItemIDFromName(sender, itemName, includeLockedByRecipe, includeLockedByTrade)
		if !exist {
			t.WithFields(testing.Fields{
				"item_name":                itemName,
				"include_locked_by_recipe": includeLockedByRecipe,
				"include_locked_by_trade":  includeLockedByTrade,
			}).Debug("no item fit params")
		}
		t.WithFields(testing.Fields{
			"item_name":                itemName,
			"include_locked_by_recipe": includeLockedByRecipe,
			"include_locked_by_trade":  includeLockedByTrade,
		}).MustNil(err, "error getting item id from name")
		ItemIDs = append(ItemIDs, itemID)
	}
	return ItemIDs
}

// GetItemInputsFromBytes is a function to get item input list from bytes
func GetItemInputsFromBytes(bytes []byte, t *testing.T) types.ItemInputList {
	var itemInputRefsReader struct {
		ItemInputs []struct {
			ID  string
			Ref string
		}
	}
	var itemInputDirectReader struct {
		ItemInputs []types.ItemInput
	}

	err := json.Unmarshal(bytes, &itemInputRefsReader)
	t.WithFields(testing.Fields{
		"bytes": string(bytes),
	}).MustNil(err, "error unmarshaling into item input refs")

	err = json.Unmarshal(bytes, &itemInputDirectReader)
	t.WithFields(testing.Fields{
		"bytes": string(bytes),
	}).MustNil(err, "error unmarshaling into item input direct")
	t.MustTrue(len(itemInputDirectReader.ItemInputs) == len(itemInputRefsReader.ItemInputs), "item input refs count and item input count does not match")

	var itemInputs types.ItemInputList

	for iii, iia := range itemInputRefsReader.ItemInputs {
		if len(iia.Ref) > 0 {
			var ii types.ItemInput
			iiBytes := ReadFile(iia.Ref, t)
			err := json.Unmarshal(iiBytes, &ii)
			if err != nil {
				t.WithFields(testing.Fields{
					"item_input_bytes": string(iiBytes),
				}).MustNil(err, "error unmarshaling item inputs")
			}
			ii.ID = iia.ID
			itemInputs = append(itemInputs, ii)
		} else {
			itemInputs = append(itemInputs, itemInputDirectReader.ItemInputs[iii])
		}
	}
	return itemInputs
}

// GetTradeItemInputsFromBytes is a function to get item input list from bytes
func GetTradeItemInputsFromBytes(bytes []byte, t *testing.T) types.TradeItemInputList {
	var tradeItemInputRefsReader struct {
		ItemInputRefs []string
	}
	err := json.Unmarshal(bytes, &tradeItemInputRefsReader)
	t.WithFields(testing.Fields{
		"bytes": string(bytes),
	}).MustNil(err, "error unmarshaling into trade item input refs")

	var itemInputs types.TradeItemInputList

	for _, tiiRef := range tradeItemInputRefsReader.ItemInputRefs {
		var tii types.TradeItemInput
		tiiBytes := ReadFile(tiiRef, t)
		err := inttest.GetJSONMarshaler().UnmarshalJSON(tiiBytes, &tii)
		t.WithFields(testing.Fields{
			"trade_item_input_bytes": string(tiiBytes),
		}).MustNil(err, "error unmarshaling trading item inputs")
		itemInputs = append(itemInputs, tii)
	}
	return itemInputs
}

// GetItemOutputsFromBytes is a function to get item outputs from bytes
func GetItemOutputsFromBytes(bytes []byte, sender string, t *testing.T) types.ItemList {
	var itemOutputNamesReader struct {
		ItemOutputNames []string
	}
	err := json.Unmarshal(bytes, &itemOutputNamesReader)
	t.WithFields(testing.Fields{
		"bytes": string(bytes),
	}).MustNil(err, "error unmarshaling into item output names")

	var itemOutputs types.ItemList

	for _, iN := range itemOutputNamesReader.ItemOutputNames {
		var io types.Item
		iID, ok, err := inttest.GetItemIDFromName(sender, iN, false, false)
		t.MustTrue(ok, "item id with specific name does not exist")
		t.WithFields(testing.Fields{
			"item_name": iN,
		}).MustNil(err, "error getting item id from name")
		io, err = inttest.GetItemByGUID(iID)
		t.WithFields(testing.Fields{
			"item_id": iID,
		}).MustNil(err, "error getting item from id")
		itemOutputs = append(itemOutputs, io)
	}
	return itemOutputs
}

// GetEntriesFromBytes is a function to get entries from bytes
func GetEntriesFromBytes(bytes []byte, t *testing.T) types.EntriesList {
	var entriesReader struct {
		Entries struct {
			CoinOutputs       []types.CoinOutput
			ItemModifyOutputs []struct {
				ID              string
				ItemInputRef    string
				ModifyParamsRef string
			}
			ItemOutputs []struct {
				ID  string
				Ref string
			}
		}
	}

	var entriesDirectReader struct {
		Entries struct {
			CoinOutputs       []types.CoinOutput
			ItemModifyOutputs []types.ItemModifyOutput
			ItemOutputs       []types.ItemOutput
		}
	}

	err := json.Unmarshal(bytes, &entriesReader)
	t.WithFields(testing.Fields{
		"bytes": string(bytes),
	}).MustNil(err, "error unmarshaling into entries reader")

	err = json.Unmarshal(bytes, &entriesDirectReader)
	t.WithFields(testing.Fields{
		"bytes": string(bytes),
	}).MustNil(err, "error unmarshaling into entries direct")
	t.MustTrue(len(entriesReader.Entries.ItemModifyOutputs) == len(entriesDirectReader.Entries.ItemModifyOutputs), "entry parsing modify outputs array length different for direct reader and ref reader")
	t.MustTrue(len(entriesReader.Entries.ItemOutputs) == len(entriesDirectReader.Entries.ItemOutputs), "entry parsing outputs array length different for direct reader and ref reader")

	var wpl types.EntriesList
	for _, co := range entriesReader.Entries.CoinOutputs {
		wpl.CoinOutputs = append(wpl.CoinOutputs, co)
	}

	for ioidx, io := range entriesReader.Entries.ItemModifyOutputs {
		if len(io.ModifyParamsRef) > 0 {
			ModifyParams := GetModifyParamsFromRef(io.ModifyParamsRef, t)
			pio := types.NewItemModifyOutput(io.ID, io.ItemInputRef, ModifyParams)

			// This is hot fix for signature verification failed issue of item output Doubles: [] instead of Doubles: nil
			if pio.Doubles != nil && len(pio.Doubles) == 0 {
				pio.Doubles = nil
			}
			if pio.Longs != nil && len(pio.Longs) == 0 {
				pio.Longs = nil
			}
			if pio.Strings != nil && len(pio.Strings) == 0 {
				pio.Strings = nil
			}
			wpl.ItemModifyOutputs = append(wpl.ItemModifyOutputs, pio)
		} else {
			wpl.ItemModifyOutputs = append(wpl.ItemModifyOutputs, entriesDirectReader.Entries.ItemModifyOutputs[ioidx])
		}
	}

	for ioidx, io := range entriesReader.Entries.ItemOutputs {
		var pio types.ItemOutput
		if len(io.Ref) > 0 {
			ioBytes := ReadFile(io.Ref, t)
			err := json.Unmarshal(ioBytes, &pio)
			t.WithFields(testing.Fields{
				"item_output_bytes": string(ioBytes),
			}).MustNil(err, "error unmarshaling into item outputs")
		} else {
			pio = entriesDirectReader.Entries.ItemOutputs[ioidx]
		}
		pio.ID = io.ID
		// This is hot fix for signature verification failed issue of item output Doubles: [] instead of Doubles: nil
		if pio.Doubles != nil && len(pio.Doubles) == 0 {
			pio.Doubles = nil
		}
		if pio.Longs != nil && len(pio.Longs) == 0 {
			pio.Longs = nil
		}
		if pio.Strings != nil && len(pio.Strings) == 0 {
			pio.Strings = nil
		}
		wpl.ItemOutputs = append(wpl.ItemOutputs, pio)
	}

	return wpl
}

// GetModifyParamsFromRef is a function to get modifying fields from reference file
func GetModifyParamsFromRef(ref string, t *testing.T) types.ItemModifyParams {
	var iup types.ItemModifyParams
	if len(ref) == 0 {
		return iup
	}
	modBytes := ReadFile(ref, t)
	err := json.Unmarshal(modBytes, &iup)
	t.WithFields(testing.Fields{
		"modify_param_bytes": string(modBytes),
	}).MustNil(err, "error unmarshaling")

	return iup
}
