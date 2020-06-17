package fixturetest

import (
	"encoding/json"
	"io/ioutil"
	"os"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"

	inttest "github.com/Pylons-tech/pylons_sdk/cmd/test"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
)

var execIDs map[string]string = make(map[string]string)

// ReadFile is a function to read file
func ReadFile(fileURL string, t *testing.T) []byte {
	jsonFile, err := os.Open(fileURL)
	if err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("fatal log")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

// UnmarshalIntoEmptyInterface is a function to convert bytes into map[string]interface{}
func UnmarshalIntoEmptyInterface(bytes []byte, t *testing.T) map[string]interface{} {
	var raw map[string]interface{}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("error unmarshaling")
	}
	return raw
}

// UpdateSenderKeyToAddress is a function to update sender key to sender's address
func UpdateSenderKeyToAddress(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	senderName, ok := raw["Sender"].(string)
	t.MustTrue(ok)
	raw["Sender"] = inttest.GetAccountAddr(senderName, t)
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
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
		t.MustNil(err)
		return newBytes
	}
	return bytes
}

// UpdateRecipeName is a function to update recipe name into recipe id
func UpdateRecipeName(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	rcpName, ok := raw["RecipeName"].(string)
	t.MustTrue(ok)
	rcpID, exist, err := inttest.GetRecipeIDFromName(rcpName)
	t.MustTrue(exist)
	t.MustNil(err)
	raw["RecipeID"] = rcpID
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
	return newBytes
}

// UpdateTradeExtraInfoToID is a function to update trade extra info into trade id
func UpdateTradeExtraInfoToID(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	trdInfo, ok := raw["TradeInfo"].(string)
	t.MustTrue(ok)
	trdID, exist, err := inttest.GetTradeIDFromExtraInfo(trdInfo)
	t.MustTrue(exist)
	t.MustNil(err)
	raw["TradeID"] = trdID
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
	return newBytes
}

// UpdateExecID is a function to set execute id from execID reference
func UpdateExecID(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	var execRefReader struct {
		ExecRef string
	}
	if err := json.Unmarshal(bytes, &execRefReader); err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("error unmarshaling")
	}

	var ok bool
	raw["ExecID"], ok = execIDs[execRefReader.ExecRef]
	if !ok {
		t.WithFields(testing.Fields{
			"execRef": execRefReader.ExecRef,
		}).Fatal("execID not available")
	}
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
	return newBytes
}

// UpdateItemIDFromName is a function to set item id from item name
func UpdateItemIDFromName(bytes []byte, includeLockedByRcp bool, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	itemName, ok := raw["ItemName"].(string)

	t.MustTrue(ok)
	itemID, exist, err := inttest.GetItemIDFromName(itemName, includeLockedByRcp)
	if !exist {
		t.WithFields(testing.Fields{
			"name":           itemName,
			"include_locked": includeLockedByRcp,
		}).Debug("no item fit params")
	}
	t.MustTrue(exist)
	t.MustNil(err)
	raw["ItemID"] = itemID
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
	return newBytes
}

// GetItemIDsFromNames is a function to set item ids from names for recipe execution
func GetItemIDsFromNames(bytes []byte, includeLockedByRcp bool, t *testing.T) []string {
	var itemNamesResp struct {
		ItemNames []string
	}
	if err := json.Unmarshal(bytes, &itemNamesResp); err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("error unmarshaling")
	}
	ItemIDs := []string{}

	for _, itemName := range itemNamesResp.ItemNames {
		itemID, exist, err := inttest.GetItemIDFromName(itemName, includeLockedByRcp)
		if !exist {
			t.WithFields(testing.Fields{
				"name":           itemName,
				"include_locked": includeLockedByRcp,
			}).Debug("no item fit params")
		}
		t.MustTrue(exist)
		t.MustNil(err)
		ItemIDs = append(ItemIDs, itemID)
	}
	return ItemIDs
}

// GetItemInputsFromBytes is a function to get item input list from bytes
func GetItemInputsFromBytes(bytes []byte, t *testing.T) types.ItemInputList {
	var itemInputRefsReader struct {
		ItemInputRefs []string
	}
	if err := json.Unmarshal(bytes, &itemInputRefsReader); err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("error unmarshaling")
	}

	var itemInputs types.ItemInputList

	for _, iiRef := range itemInputRefsReader.ItemInputRefs {
		var ii types.ItemInput
		iiBytes := ReadFile(iiRef, t)
		err := inttest.GetAminoCdc().UnmarshalJSON(iiBytes, &ii)
		if err != nil {
			t.WithFields(testing.Fields{
				"error": err,
			}).Fatal("error unmarshaling")
		}
		itemInputs = append(itemInputs, ii)
	}
	return itemInputs
}

// GetTradeItemInputsFromBytes is a function to get item input list from bytes
func GetTradeItemInputsFromBytes(bytes []byte, t *testing.T) types.TradeItemInputList {
	var itemInputRefsReader struct {
		ItemInputRefs []string
	}
	if err := json.Unmarshal(bytes, &itemInputRefsReader); err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("error unmarshaling")
	}

	var itemInputs types.TradeItemInputList

	for _, tiiRef := range itemInputRefsReader.ItemInputRefs {
		var tii types.TradeItemInput
		tiiBytes := ReadFile(tiiRef, t)
		err := inttest.GetAminoCdc().UnmarshalJSON(tiiBytes, &tii)
		if err != nil {
			t.WithFields(testing.Fields{
				"error": err,
			}).Fatal("error unmarshaling")
		}
		itemInputs = append(itemInputs, tii)
	}
	return itemInputs
}

// GetItemOutputsFromBytes is a function to get item outputs from bytes
func GetItemOutputsFromBytes(bytes []byte, sender string, t *testing.T) types.ItemList {
	var itemOutputNamesReader struct {
		ItemOutputNames []string
	}
	if err := json.Unmarshal(bytes, &itemOutputNamesReader); err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("error unmarshaling")
	}

	var itemOutputs types.ItemList

	for _, iN := range itemOutputNamesReader.ItemOutputNames {
		var io types.Item
		iID, ok, err := inttest.GetItemIDFromName(iN, false)
		t.MustNil(err)
		t.MustTrue(ok)
		io, err = inttest.GetItemByGUID(iID)
		t.MustNil(err)
		itemOutputs = append(itemOutputs, io)
	}
	return itemOutputs
}

// GetEntriesFromBytes is a function to get entries from bytes
func GetEntriesFromBytes(bytes []byte, t *testing.T) types.EntriesList {
	var entriesReader struct {
		Entries struct {
			CoinOutputs []types.CoinOutput
			ItemOutputs []struct {
				Ref        string
				ModifyItem struct {
					ItemInputRef    *int
					ModifyParamsRef string
				}
			}
		}
	}

	if err := json.Unmarshal(bytes, &entriesReader); err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("error unmarshaling")
	}

	var wpl types.EntriesList
	for _, co := range entriesReader.Entries.CoinOutputs {
		wpl = append(wpl, co)
	}

	for _, io := range entriesReader.Entries.ItemOutputs {
		var pio types.ItemOutput
		if len(io.Ref) > 0 {
			ioBytes := ReadFile(io.Ref, t)
			err := json.Unmarshal(ioBytes, &pio)
			if err != nil {
				t.WithFields(testing.Fields{
					"item_output_bytes": string(ioBytes),
					"error":             err,
				}).Fatal("error unmarshaling")
			}
		}
		if io.ModifyItem.ItemInputRef == nil {
			pio.ModifyItem.ItemInputRef = -1
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
		} else {
			pio.ModifyItem.ItemInputRef = *io.ModifyItem.ItemInputRef
		}
		ModifyParams := GetModifyParamsFromRef(io.ModifyItem.ModifyParamsRef, t)
		pio.ModifyItem.Doubles = ModifyParams.Doubles
		pio.ModifyItem.Longs = ModifyParams.Longs
		pio.ModifyItem.Strings = ModifyParams.Strings
		// This is hot fix for signature verification failed issue of item output Doubles: [] instead of Doubles: nil
		if pio.ModifyItem.Doubles != nil && len(pio.ModifyItem.Doubles) == 0 {
			pio.ModifyItem.Doubles = nil
		}
		if pio.ModifyItem.Longs != nil && len(pio.ModifyItem.Longs) == 0 {
			pio.ModifyItem.Longs = nil
		}
		if pio.ModifyItem.Strings != nil && len(pio.ModifyItem.Strings) == 0 {
			pio.ModifyItem.Strings = nil
		}
		wpl = append(wpl, pio)
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
	if err != nil {
		t.WithFields(testing.Fields{
			"modify_param_bytes": string(modBytes),
			"error":              err,
		}).Fatal("error unmarshaling")
	}

	return iup
}
