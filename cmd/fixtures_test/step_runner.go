package fixturetest

import (
	"encoding/json"
	"strings"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"

	inttest "github.com/Pylons-tech/pylons_sdk/cmd/test"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"

	"github.com/Pylons-tech/pylons_sdk/x/pylons/handlers"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TxBroadcastErrorCheck check error is same as expected when it exist
func TxBroadcastErrorCheck(err error, txhash string, step FixtureStep, t *testing.T) {
	if step.Output.TxResult.BroadcastError != "" {
		t.WithFields(testing.Fields{
			"original_error": err.Error(),
			"target_error":   step.Output.TxResult.BroadcastError,
		}).MustTrue(strings.Contains(err.Error(), step.Output.TxResult.BroadcastError), "broadcast error is different from expected one")
	} else {
		t.WithFields(testing.Fields{
			"txhash": txhash,
			"error":  err,
		}).Fatal("unexpected transaction broadcast error")
	}
}

// TxResultStatusMessageCheck check result status and message
func TxResultStatusMessageCheck(status, message, txhash string, step FixtureStep, t *testing.T) {
	if len(step.Output.TxResult.Status) > 0 {
		t.WithFields(testing.Fields{
			"txhash":          txhash,
			"original_status": status,
			"target_status":   step.Output.TxResult.Status,
		}).MustTrue(status == step.Output.TxResult.Status, "transaction result status is different from expected")
	}
	if len(step.Output.TxResult.Message) > 0 {
		t.WithFields(testing.Fields{
			"txhash":           txhash,
			"original_message": message,
			"target_message":   step.Output.TxResult.Message,
		}).MustTrue(message == step.Output.TxResult.Message, "transaction result message is different from expected")
	}
}

// TxResultDecodingErrorCheck check error for tx response data unmarshal
func TxResultDecodingErrorCheck(err error, txhash string, t *testing.T) {
	if err != nil {
		t.WithFields(testing.Fields{
			"txhash": txhash,
			"error":  err,
		}).Fatal("error unmarshaling tx response")
	}
}

// GetTxHandleResult check error on tx by hash and return handle result
func GetTxHandleResult(txhash string, t *testing.T) []byte {
	txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
	if err != nil {
		t.WithFields(testing.Fields{
			"tx_result_bytes": string(txHandleResBytes),
			"error":           err,
		}).Fatal("error getting tx result bytes")
	}
	CheckErrorOnTxFromTxHash(txhash, t)
	return txHandleResBytes
}

// WaitForNextBlockWithErrorCheck wait 1 block and check the error result
func WaitForNextBlockWithErrorCheck(t *testing.T) {
	err := inttest.WaitForNextBlock()
	if err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("error waiting for check execution")
	}
}

// RunMultiMsgTx is a function to send multiple messages in a transaction
// This support only 1 sender multi transaction for now
// TODO we need to support multi-message multi sender transaction
func RunMultiMsgTx(step FixtureStep, t *testing.T) {
	if len(step.MsgRefs) != 0 {
		var msgs []sdk.Msg
		var sender sdk.AccAddress
		for _, ref := range step.MsgRefs {
			var newMsg sdk.Msg
			switch ref.Action {
			case "fiat_item":
				msg := FiatItemMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "update_item_string":
				msg := UpdateItemStringMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "create_cookbook":
				msg := CreateCookbookMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "create_recipe":
				msg := CreateRecipeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "execute_recipe":
				msg := ExecuteRecipeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "check_execution":
				msg := CheckExecutionMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "create_trade":
				msg := CreateTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "fulfill_trade":
				msg := FulfillTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "disable_trade":
				msg := DisableTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			}
			msgs = append(msgs, newMsg)
		}
		t.WithFields(testing.Fields{
			"sender":   sender.String(),
			"len_msgs": len(msgs),
			"tx_msgs":  inttest.AminoCodecFormatter(msgs),
			"msg_refs": step.MsgRefs,
		}).AddFields(inttest.GetLogFieldsFromMsgs(msgs)).Debug("debug log")
		txhash, err := inttest.SendMultiMsgTxWithNonce(t, msgs, sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)
		GetTxHandleResult(txhash, t)
		t.WithFields(testing.Fields{
			"txhash": txhash,
		}).Debug("debug log")
	}
}

// CheckExecutionMsgFromRef collect check execution message from reference string
func CheckExecutionMsgFromRef(ref string, t *testing.T) msgs.MsgCheckExecution {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate execRef to execID
	newByteValue = UpdateExecID(newByteValue, t)

	var execType struct {
		ExecID        string
		PayToComplete bool
		Sender        sdk.AccAddress
	}
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &execType)
	if err != nil {
		t.WithFields(testing.Fields{
			"execType": inttest.AminoCodecFormatter(execType),
			"error":    err,
		}).Fatal("error reading using GetAminoCdc")
	}

	return msgs.NewMsgCheckExecution(
		execType.ExecID,
		execType.PayToComplete,
		execType.Sender,
	)
}

// RunCheckExecution is a function to execute check execution
func RunCheckExecution(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		chkExecMsg := CheckExecutionMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, chkExecMsg, chkExecMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.CheckExecutionResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// FiatItemMsgFromRef collect check execution message from reference string
func FiatItemMsgFromRef(ref string, t *testing.T) msgs.MsgFiatItem {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate cookbook name to cookbook ID
	newByteValue = UpdateCBNameToID(newByteValue, t)

	var itemType types.Item
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &itemType)
	if err != nil {
		t.WithFields(testing.Fields{
			"itemType": inttest.AminoCodecFormatter(itemType),
			"error":    err,
		}).Fatal("error reading using GetAminoCdc itemType", itemType, err)
	}

	t.Fatal(itemType)

	return msgs.NewMsgFiatItem(
		itemType.CookbookID,
		itemType.Doubles,
		itemType.Longs,
		itemType.Strings,
		itemType.Sender,
		itemType.AdditionalTransferFee,
	)
}

// RunFiatItem is a function to execute fiat item
func RunFiatItem(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		itmMsg := FiatItemMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, itmMsg, itmMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.FiatItemResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.ItemID != "", "item id shouldn't be empty")
	}
}

// SendItemsMsgFromRef is a function to collect SendItems from reference string
func SendItemsMsgFromRef(ref string, t *testing.T) msgs.MsgSendItems {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	newByteValue = UpdateReceiverKeyToAddress(newByteValue, t)

	ItemIDs := GetItemIDsFromNames(newByteValue, false, t)

	var siType struct {
		Sender   sdk.AccAddress
		Receiver sdk.AccAddress
		ItemIDs  []string `json:"ItemIDs"`
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &siType)
	if err != nil {
		t.WithFields(testing.Fields{
			"siType":    inttest.AminoCodecFormatter(siType),
			"new_bytes": string(newByteValue),
			"error":     err,
		}).Fatal("error reading using GetAminoCdc")
	}

	return msgs.NewMsgSendItems(ItemIDs, siType.Sender, siType.Receiver)
}

// RunSendItems is a function to send items to another user
func RunSendItems(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		siMsg := SendItemsMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, siMsg, siMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		if len(step.Output.TxResult.ErrorLog) > 0 {
		} else {
			txHandleResBytes := GetTxHandleResult(txhash, t)
			resp := handlers.SendItemsResponse{}
			err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
			if err != nil {
				t.WithFields(testing.Fields{
					"txhash": txhash,
				}).Fatal("failed to parse transaction result")
			}
			TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
		}
	}
}

// UpdateItemStringMsgFromRef is a function to collect UpdateItemStringMsg from reference string
func UpdateItemStringMsgFromRef(ref string, t *testing.T) msgs.MsgUpdateItemString {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate item name to item ID
	newByteValue = UpdateItemIDFromName(newByteValue, false, t)

	var sTypeMsg msgs.MsgUpdateItemString
	err := json.Unmarshal(newByteValue, &sTypeMsg)
	if err != nil {
		t.WithFields(testing.Fields{
			"sTypeMsg":  inttest.AminoCodecFormatter(sTypeMsg),
			"new_bytes": string(newByteValue),
			"error":     err,
		}).Fatal("error reading using GetAminoCdc")
	}
	return sTypeMsg
}

// RunUpdateItemString is a function to update item's string value
func RunUpdateItemString(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		sTypeMsg := UpdateItemStringMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, sTypeMsg, sTypeMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}
		WaitForNextBlockWithErrorCheck(t)
		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.UpdateItemStringResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
	}
}

// CreateCookbookMsgFromRef is a function to get create cookbook message from reference
func CreateCookbookMsgFromRef(ref string, t *testing.T) msgs.MsgCreateCookbook {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)

	var cbType types.Cookbook
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &cbType)
	if err != nil {
		t.WithFields(testing.Fields{
			"cbType":    inttest.AminoCodecFormatter(cbType),
			"new_bytes": string(newByteValue),
			"error":     err,
		}).Fatal("error reading using GetAminoCdc")
	}

	return msgs.NewMsgCreateCookbook(
		cbType.Name,
		cbType.ID,
		cbType.Description,
		cbType.Developer,
		cbType.Version,
		cbType.SupportEmail,
		cbType.Level,
		cbType.CostPerBlock,
		cbType.Sender,
	)
}

// RunCreateCookbook is a function to create cookbook
func RunCreateCookbook(step FixtureStep, t *testing.T) {
	if !FixtureTestOpts.CreateNewCookbook {
		return
	}
	if step.ParamsRef != "" {
		cbMsg := CreateCookbookMsgFromRef(step.ParamsRef, t)

		txhash, err := inttest.TestTxWithMsgWithNonce(t, cbMsg, cbMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.CreateCookbookResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.CookbookID != "", "coookbook id shouldn't be empty")
	}
}

// CreateRecipeMsgFromRef is a function to get create cookbook message from reference
func CreateRecipeMsgFromRef(ref string, t *testing.T) msgs.MsgCreateRecipe {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate cookbook name to cookbook id
	newByteValue = UpdateCBNameToID(newByteValue, t)
	// get item inputs from fileNames
	itemInputs := GetItemInputsFromBytes(newByteValue, t)
	// get entries from fileNames
	entries := GetEntriesFromBytes(newByteValue, t)

	var rcpTempl types.Recipe
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &rcpTempl)
	if err != nil {
		t.WithFields(testing.Fields{
			"rcpTempl":  inttest.AminoCodecFormatter(rcpTempl),
			"new_bytes": string(newByteValue),
			"error":     err,
		}).Fatal("error reading using GetAminoCdc")
	}

	return msgs.NewMsgCreateRecipe(
		rcpTempl.Name,
		rcpTempl.CookbookID,
		rcpTempl.ID,
		rcpTempl.Description,
		rcpTempl.CoinInputs,
		itemInputs,
		entries,
		rcpTempl.Outputs,
		rcpTempl.BlockInterval,
		rcpTempl.Sender,
	)
}

// RunCreateRecipe is a function to create recipe
func RunCreateRecipe(step FixtureStep, t *testing.T) {
	if step.ParamsRef != "" {
		rcpMsg := CreateRecipeMsgFromRef(step.ParamsRef, t)

		txhash, err := inttest.TestTxWithMsgWithNonce(t, rcpMsg, rcpMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)
		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.CreateRecipeResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.RecipeID != "", "recipe id shouldn't be empty")
	}
}

// ExecuteRecipeMsgFromRef collect execute recipe msg from reference string
func ExecuteRecipeMsgFromRef(ref string, t *testing.T) msgs.MsgExecuteRecipe {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate recipe name to recipe id
	newByteValue = UpdateRecipeName(newByteValue, t)
	// translate itemNames to itemIDs
	ItemIDs := GetItemIDsFromNames(newByteValue, false, t)

	var execType struct {
		RecipeID string
		Sender   sdk.AccAddress
		ItemIDs  []string `json:"ItemIDs"`
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &execType)
	if err != nil {
		t.WithFields(testing.Fields{
			"execType":  inttest.AminoCodecFormatter(execType),
			"new_bytes": string(newByteValue),
			"error":     err,
		}).Fatal("error reading using GetAminoCdc")
	}

	return msgs.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
}

// RunExecuteRecipe is executed when an action "execute_recipe" is called
func RunExecuteRecipe(step FixtureStep, t *testing.T) {
	// TODO should check item ID is returned
	// TODO when items are generated, rather than returning whole should return only ID [if multiple, array of item IDs]

	if step.ParamsRef != "" {
		execMsg := ExecuteRecipeMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, execMsg, execMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			hmrErrMsg := inttest.GetHumanReadableErrorFromTxHash(txhash, t)

			t.WithFields(testing.Fields{
				"tx_error":       hmrErrMsg,
				"expected_error": step.Output.TxResult.ErrorLog,
			}).MustTrue(strings.Contains(hmrErrMsg, step.Output.TxResult.ErrorLog), "transaction error log is different from expected one.")
		} else {
			txHandleResBytes := GetTxHandleResult(txhash, t)
			resp := handlers.ExecuteRecipeResponse{}
			err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
			if err != nil {
				t.WithFields(testing.Fields{
					"txhash": txhash,
				}).Fatal("failed to parse transaction result")
			}
			TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)

			if resp.Message == "scheduled the recipe" { // delayed execution
				var scheduleRes handlers.ExecuteRecipeScheduleOutput

				err := json.Unmarshal(resp.Output, &scheduleRes)
				t.WithFields(testing.Fields{
					"response_output": string(resp.Output),
				}).MustNil(err, "something went wrong decoding raw json")
				execIDs[step.ID] = scheduleRes.ExecID
				for _, itemID := range execMsg.ItemIDs {
					item, err := inttest.GetItemByGUID(itemID)
					t.WithFields(testing.Fields{
						"item_id": itemID,
					}).MustNil(err, "there's an issue while getting item from id")
					t.MustTrue(len(item.OwnerRecipeID) != 0, "OwnerRecipeID shouldn't be set but it's set")
				}

				t.WithFields(testing.Fields{
					"exec_id": scheduleRes.ExecID,
				}).Debug("scheduled execution")
			} else { // straight execution
				t.WithFields(testing.Fields{
					"output": string(resp.Output),
				}).Debug("straight execution result")
			}
		}
	}
}

// CreateTradeMsgFromRef collect create trade msg from reference
func CreateTradeMsgFromRef(ref string, t *testing.T) msgs.MsgCreateTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// get item inputs from fileNames
	tradeItemInputs := GetTradeItemInputsFromBytes(newByteValue, t)
	var trdType types.Trade
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &trdType)
	if err != nil {
		t.WithFields(testing.Fields{
			"trdType":   inttest.AminoCodecFormatter(trdType),
			"new_bytes": string(newByteValue),
			"error":     err,
		}).Fatal("error reading using GetAminoCdc")
	}

	// get ItemOutputs from ItemOutputNames
	itemOutputs := GetItemOutputsFromBytes(newByteValue, trdType.Sender.String(), t)

	return msgs.NewMsgCreateTrade(
		trdType.CoinInputs,
		tradeItemInputs,
		trdType.CoinOutputs,
		itemOutputs,
		trdType.ExtraInfo,
		trdType.Sender,
	)
}

// RunCreateTrade is a function to create trade
func RunCreateTrade(step FixtureStep, t *testing.T) {
	if step.ParamsRef != "" {
		createTrd := CreateTradeMsgFromRef(step.ParamsRef, t)
		t.WithFields(testing.Fields{
			"tx_msgs": inttest.AminoCodecFormatter(createTrd),
		}).AddFields(inttest.GetLogFieldsFromMsgs([]sdk.Msg{createTrd})).Debug("createTrd")
		txhash, err := inttest.TestTxWithMsgWithNonce(t, createTrd, createTrd.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)
		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.CreateTradeResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.TradeID != "", "trade id shouldn't be empty")
	}
}

// FulfillTradeMsgFromRef collect fulfill trade message from reference string
func FulfillTradeMsgFromRef(ref string, t *testing.T) msgs.MsgFulfillTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate extra info to trade id
	newByteValue = UpdateTradeExtraInfoToID(newByteValue, t)
	// translate itemNames to itemIDs
	ItemIDs := GetItemIDsFromNames(newByteValue, false, t)

	var trdType struct {
		TradeID string
		Sender  sdk.AccAddress
		ItemIDs []string `json:"ItemIDs"`
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &trdType)
	if err != nil {
		t.WithFields(testing.Fields{
			"trdType":   inttest.AminoCodecFormatter(trdType),
			"new_bytes": string(newByteValue),
			"error":     err,
		}).Fatal("error reading using GetAminoCdc")
	}

	return msgs.NewMsgFulfillTrade(trdType.TradeID, trdType.Sender, ItemIDs)
}

// RunFulfillTrade is a function to fulfill trade
func RunFulfillTrade(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		ffTrdMsg := FulfillTradeMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, ffTrdMsg, ffTrdMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		if len(step.Output.TxResult.ErrorLog) > 0 {
		} else {
			txHandleResBytes := GetTxHandleResult(txhash, t)
			resp := handlers.FulfillTradeResponse{}
			err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
			if err != nil {
				t.WithFields(testing.Fields{
					"txhash": txhash,
				}).Fatal("failed to parse transaction result")
			}
			TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
		}
	}
}

// DisableTradeMsgFromRef collect disable trade msg from reference string
func DisableTradeMsgFromRef(ref string, t *testing.T) msgs.MsgDisableTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate extra info to trade id
	newByteValue = UpdateTradeExtraInfoToID(newByteValue, t)

	var trdType struct {
		TradeID string
		Sender  sdk.AccAddress
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &trdType)
	if err != nil {
		t.WithFields(testing.Fields{
			"trdType":   inttest.AminoCodecFormatter(trdType),
			"new_bytes": string(newByteValue),
			"error":     err,
		}).Fatal("error reading using GetAminoCdc")
	}

	return msgs.NewMsgDisableTrade(trdType.TradeID, trdType.Sender)
}

// RunDisableTrade is a function to disable trade
func RunDisableTrade(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		dsTrdMsg := DisableTradeMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, dsTrdMsg, dsTrdMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		if len(step.Output.TxResult.ErrorLog) > 0 {
		} else {
			txHandleResBytes := GetTxHandleResult(txhash, t)
			resp := handlers.DisableTradeResponse{}
			err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
			if err != nil {
				t.WithFields(testing.Fields{
					"txhash": txhash,
				}).Fatal("failed to parse transaction result")
			}
			TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
		}
	}
}
