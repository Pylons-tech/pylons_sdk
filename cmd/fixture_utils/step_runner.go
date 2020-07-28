package fixturetest

import (
	"encoding/json"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/evtesting"

	inttest "github.com/Pylons-tech/pylons_sdk/cmd/test_utils"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"

	"github.com/Pylons-tech/pylons_sdk/x/pylons/handlers"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TxBroadcastErrorCheck check error is same as expected when it exist
func TxBroadcastErrorCheck(err error, txhash string, step FixtureStep, t *testing.T) {
	if step.Output.TxResult.BroadcastError != "" {
		t.WithFields(testing.Fields{
			"txhash": txhash,
		}).MustContain(err.Error(), step.Output.TxResult.BroadcastError, "broadcast error is different from expected one")
	} else {
		t.WithFields(testing.Fields{
			"txhash": txhash,
		}).MustNil(err, "unexpected transaction broadcast error")
	}
}

// TxErrorLogCheck check expected error log is produced correctly
func TxErrorLogCheck(txhash string, ErrorLog string, t *testing.T) {
	if len(ErrorLog) > 0 {
		hmrErrMsg := inttest.GetHumanReadableErrorFromTxHash(txhash, t)
		t.WithFields(testing.Fields{
			"txhash": txhash,
		}).MustContain(hmrErrMsg, ErrorLog, "transaction error log is different from expected one.")
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
	txErrorBytes, getTxLogErr := inttest.GetTxError(txhash, t)
	t.WithFields(testing.Fields{
		"txhash":         txhash,
		"tx_error_bytes": string(txErrorBytes),
		"get_tx_log_err": getTxLogErr,
	}).MustNil(err, "error unmarshaling tx response")
}

// GetTxHandleResult check error on tx by hash and return handle result
func GetTxHandleResult(txhash string, t *testing.T) []byte {
	txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
	t.WithFields(testing.Fields{
		"tx_result_bytes": string(txHandleResBytes),
		"error":           err,
	}).MustNil(err, "error getting tx result bytes")
	CheckErrorOnTxFromTxHash(txhash, t)
	return txHandleResBytes
}

// WaitForNextBlockWithErrorCheck wait 1 block and check the error result
func WaitForNextBlockWithErrorCheck(t *testing.T) {
	err := inttest.WaitForNextBlock()
	t.MustNil(err, "error waiting for next block")
}

// RunCreateAccount is a function to create account
func RunCreateAccount(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		caKey := GetAccountKeyFromTempName(step.ParamsRef, t)
		localKeyResult, err := inttest.AddNewLocalKey(caKey)
		t.WithFields(testing.Fields{
			"key":              caKey,
			"local_key_result": localKeyResult,
		}).MustNil(err, "error creating local Key")
		result, logstr, err := inttest.CreateChainAccount(caKey)
		t.WithFields(testing.Fields{
			"result": result,
			"logstr": logstr,
		}).MustNil(err, "error creating account on chain")

		caTxHash := inttest.GetTxHashFromLog(result)
		t.MustTrue(caTxHash != "", "error fetching txhash from result")
		t.WithFields(testing.Fields{
			"txhash": caTxHash,
		}).Info("waiting for create account transaction")
		txResponseBytes, err := inttest.WaitAndGetTxData(caTxHash, inttest.GetMaxWaitBlock(), t)
		t.WithFields(testing.Fields{
			"result": string(txResponseBytes),
		}).MustNil(err, "error waiting for create account transaction")
		inttest.GetAccountInfoFromAddr(localKeyResult["address"], t)
	}
}

// GetPylonsMsgFromRef is a function to get GetPylons message from reference
func GetPylonsMsgFromRef(ref string, t *testing.T) msgs.MsgGetPylons {
	gpAddr := GetAccountAddressFromTempName(ref, t)
	sdkAddr, err := sdk.AccAddressFromBech32(gpAddr)
	t.WithFields(testing.Fields{
		"temp_name": ref,
	}).MustNil(err, "error converting key to address")
	return msgs.NewMsgGetPylons(
		types.NewPylon(55000),
		sdkAddr,
	)
}

// SendCoinsMsgFromRef is a function to SendCoins message from reference
func SendCoinsMsgFromRef(ref string, t *testing.T) msgs.MsgSendCoins {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	newByteValue = UpdateReceiverKeyToAddress(newByteValue, t)

	var siType struct {
		Sender   sdk.AccAddress
		Receiver sdk.AccAddress
		Amount   sdk.Coins
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &siType)
	t.WithFields(testing.Fields{
		"siType":    inttest.AminoCodecFormatter(siType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")

	return msgs.NewMsgSendCoins(siType.Amount, siType.Sender, siType.Receiver)
}

// RunGetPylons is a function to run GetPylos message
func RunGetPylons(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		gpMsg := GetPylonsMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, gpMsg, gpMsg.Requester.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.GetPylonsResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// RunSendCoins is a function to send coins from one address to another
func RunSendCoins(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		scMsg := SendCoinsMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, scMsg, scMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.GetPylonsResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// RunMockAccount = RunCreateAccount + RunGetPylons
func RunMockAccount(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		RunCreateAccount(step, t)
		RunGetPylons(step, t)
	}
}

// RunMultiMsgTx is a function to send multiple messages in a transaction
// This support only 1 sender multi transaction for now
// TODO we need to support multi-message multi sender transaction
func RunMultiMsgTx(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
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
	t.WithFields(testing.Fields{
		"execType": inttest.AminoCodecFormatter(execType),
	}).MustNil(err, "error reading using GetAminoCdc")

	return msgs.NewMsgCheckExecution(
		execType.ExecID,
		execType.PayToComplete,
		execType.Sender,
	)
}

// RunCheckExecution is a function to execute check execution
func RunCheckExecution(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		chkExecMsg := CheckExecutionMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, chkExecMsg, chkExecMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

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
	t.WithFields(testing.Fields{
		"itemType": inttest.AminoCodecFormatter(itemType),
	}).MustNil(err, "error reading using GetAminoCdc")

	return msgs.NewMsgFiatItem(
		itemType.CookbookID,
		itemType.Doubles,
		itemType.Longs,
		itemType.Strings,
		itemType.Sender,
		itemType.TransferFee,
	)
}

// RunFiatItem is a function to execute fiat item
func RunFiatItem(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		itmMsg := FiatItemMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, itmMsg, itmMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

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

	ItemIDs := GetItemIDsFromNames(newByteValue, false, false, t)

	var siType struct {
		Sender   sdk.AccAddress
		Receiver sdk.AccAddress
		ItemIDs  []string `json:"ItemIDs"`
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &siType)
	t.WithFields(testing.Fields{
		"siType":    inttest.AminoCodecFormatter(siType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")

	return msgs.NewMsgSendItems(ItemIDs, siType.Sender, siType.Receiver)
}

// RunSendItems is a function to send items to another user
func RunSendItems(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		siMsg := SendItemsMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, siMsg, siMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.SendItemsResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// UpdateItemStringMsgFromRef is a function to collect UpdateItemStringMsg from reference string
func UpdateItemStringMsgFromRef(ref string, t *testing.T) msgs.MsgUpdateItemString {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate item name to item ID
	newByteValue = UpdateItemIDFromName(newByteValue, false, false, t)

	var sTypeMsg msgs.MsgUpdateItemString
	err := json.Unmarshal(newByteValue, &sTypeMsg)
	t.WithFields(testing.Fields{
		"sTypeMsg":  inttest.AminoCodecFormatter(sTypeMsg),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")
	return sTypeMsg
}

// RunUpdateItemString is a function to update item's string value
func RunUpdateItemString(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		sTypeMsg := UpdateItemStringMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, sTypeMsg, sTypeMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}
		WaitForNextBlockWithErrorCheck(t)

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

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
	t.WithFields(testing.Fields{
		"cbType":    inttest.AminoCodecFormatter(cbType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")

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
	if !FixtureTestOpts.CreateNewCookbook || FixtureTestOpts.VerifyOnly {
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

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.CreateCookbookResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.CookbookID != "", "coookbook id shouldn't be empty")
	}
}

// RunMockCookbook = RunMockAccount + RunCreateCookbook
func RunMockCookbook(step FixtureStep, t *testing.T) {
	if !FixtureTestOpts.CreateNewCookbook || FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		sender := GetSenderKeyFromRef(step.ParamsRef, t)
		RunMockAccount(FixtureStep{ParamsRef: sender}, t)
		RunCreateCookbook(step, t)
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
	t.WithFields(testing.Fields{
		"rcpTempl":  inttest.AminoCodecFormatter(rcpTempl),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")

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
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		rcpMsg := CreateRecipeMsgFromRef(step.ParamsRef, t)

		txhash, err := inttest.TestTxWithMsgWithNonce(t, rcpMsg, rcpMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

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
	ItemIDs := GetItemIDsFromNames(newByteValue, false, false, t)

	var execType struct {
		RecipeID string
		Sender   sdk.AccAddress
		ItemIDs  []string `json:"ItemIDs"`
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &execType)
	t.WithFields(testing.Fields{
		"execType":  inttest.AminoCodecFormatter(execType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")

	return msgs.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
}

// RunExecuteRecipe is executed when an action "execute_recipe" is called
func RunExecuteRecipe(step FixtureStep, t *testing.T) {
	// TODO should check item ID is returned
	// TODO when items are generated, rather than returning whole should return only ID [if multiple, array of item IDs]

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		execMsg := ExecuteRecipeMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, execMsg, execMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.ExecuteRecipeResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)

		if resp.Message == "scheduled the recipe" { // delayed execution
			var scheduleRes handlers.ExecuteRecipeScheduleOutput

			err := json.Unmarshal(resp.Output, &scheduleRes)
			t.WithFields(testing.Fields{
				"response_output": string(resp.Output),
			}).MustNil(err, "error decoding raw json")
			execIDs[step.ID] = scheduleRes.ExecID
			for _, itemID := range execMsg.ItemIDs {
				item, err := inttest.GetItemByGUID(itemID)
				t.WithFields(testing.Fields{
					"item_id": itemID,
				}).MustNil(err, "error getting item from id")
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

// CreateTradeMsgFromRef collect create trade msg from reference
func CreateTradeMsgFromRef(ref string, t *testing.T) msgs.MsgCreateTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// get item inputs from fileNames
	tradeItemInputs := GetTradeItemInputsFromBytes(newByteValue, t)
	var trdType types.Trade
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &trdType)
	t.WithFields(testing.Fields{
		"trdType":   inttest.AminoCodecFormatter(trdType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")

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
	if FixtureTestOpts.VerifyOnly {
		return
	}
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

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

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
	ItemIDs := GetItemIDsFromNames(newByteValue, false, false, t)

	var trdType struct {
		TradeID string
		Sender  sdk.AccAddress
		ItemIDs []string `json:"ItemIDs"`
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &trdType)
	t.WithFields(testing.Fields{
		"trdType":   inttest.AminoCodecFormatter(trdType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")

	return msgs.NewMsgFulfillTrade(trdType.TradeID, trdType.Sender, ItemIDs)
}

// RunFulfillTrade is a function to fulfill trade
func RunFulfillTrade(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		ffTrdMsg := FulfillTradeMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, ffTrdMsg, ffTrdMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.FulfillTradeResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
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
	t.WithFields(testing.Fields{
		"trdType":   inttest.AminoCodecFormatter(trdType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")

	return msgs.NewMsgDisableTrade(trdType.TradeID, trdType.Sender)
}

// RunDisableTrade is a function to disable trade
func RunDisableTrade(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		dsTrdMsg := DisableTradeMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, dsTrdMsg, dsTrdMsg.Sender.String(), true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		TxErrorLogCheck(txhash, step.Output.TxResult.ErrorLog, t)
		if len(step.Output.TxResult.ErrorLog) > 0 {
			return
		}

		txHandleResBytes := GetTxHandleResult(txhash, t)
		resp := handlers.DisableTradeResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}
