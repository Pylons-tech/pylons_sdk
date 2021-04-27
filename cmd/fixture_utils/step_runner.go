package fixturetest

import (
	"encoding/base64"
	"encoding/json"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/evtesting"
	inttest "github.com/Pylons-tech/pylons_sdk/cmd/test_utils"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/handlers"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
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

		caTxHash, err := inttest.GetTxHashFromJson(result)
		t.MustNil(err, "error code detected parsing result json")
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
func GetPylonsMsgFromRef(ref string, t *testing.T) types.MsgGetPylons {
	gpAddr := GetAccountAddressFromTempName(ref, t)
	sdkAddr, err := sdk.AccAddressFromBech32(gpAddr)
	t.WithFields(testing.Fields{
		"temp_name": ref,
	}).MustNil(err, "error converting key to address")
	return types.NewMsgGetPylons(
		types.NewPylon(55000),
		sdkAddr.String(),
	)
}

// RunGetPylons is a function to run GetPylos message
func RunGetPylons(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		gpMsg := GetPylonsMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &gpMsg, gpMsg.Requester, true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		txHandleResBytes := GetTxHandleResult(txhash, t)
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgGetPylons{}).Type(), "MsgType should be accurate")
		resp := types.MsgGetPylonsResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// GoogleIAPGetPylonsMsgFromRef is a function to get GoogleIAPGetPylons message from reference
func GoogleIAPGetPylonsMsgFromRef(ref string, t *testing.T) types.MsgGoogleIAPGetPylons {
	byteValue := ReadFile(ref, t)
	// translate requester from account name to account address
	newByteValue := UpdateRequesterKeyToAddress(byteValue, t)

	var gigpType struct {
		ProductID     string
		PurchaseToken string
		ReceiptData   string
		Signature     string
		Requester     string
	}

	err := json.Unmarshal(newByteValue, &gigpType)
	t.WithFields(testing.Fields{
		"gigpType":  inttest.AminoCodecFormatter(gigpType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")

	receiptDataBase64 := base64.StdEncoding.EncodeToString([]byte(gigpType.ReceiptData))

	return types.NewMsgGoogleIAPGetPylons(
		gigpType.ProductID,
		gigpType.PurchaseToken,
		receiptDataBase64,
		gigpType.Signature,
		gigpType.Requester,
	)
}

// RunGoogleIAPGetPylons is a function to run GoogleIAPGetPylons message
func RunGoogleIAPGetPylons(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		gigpMsg := GoogleIAPGetPylonsMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &gigpMsg, gigpMsg.Requester, true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)

		txHandleResBytes := GetTxHandleResult(txhash, t)
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgGoogleIAPGetPylons{}).Type(), "MsgType should be accurate")
		resp := types.MsgGoogleIAPGetPylonsResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// SendCoinsMsgFromRef is a function to SendCoins message from reference
func SendCoinsMsgFromRef(ref string, t *testing.T) types.MsgSendCoins {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	newByteValue = UpdateReceiverKeyToAddress(newByteValue, t)

	var siType struct {
		Sender   string
		Receiver string
		Amount   string
	}

	err := json.Unmarshal(newByteValue, &siType)
	t.WithFields(testing.Fields{
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using json Unmarshaler")

	amount, err := sdk.ParseCoinsNormalized(siType.Amount)
	t.WithFields(testing.Fields{
		"amount": siType.Amount,
	}).MustNil(err, "error parsing amount")

	return types.NewMsgSendCoins(amount, siType.Sender, siType.Receiver)
}

// RunSendCoins is a function to send coins from one address to another
func RunSendCoins(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		scMsg := SendCoinsMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &scMsg, scMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgSendCoins{}).Type(), "MsgType should be accurate")
		resp := types.MsgGetPylonsResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
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
		var sender string
		for _, ref := range step.MsgRefs {
			var newMsg sdk.Msg
			switch ref.Action {
			case "fiat_item":
				msg := FiatItemMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "update_item_string":
				msg := UpdateItemStringMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "create_cookbook":
				msg := CreateCookbookMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "update_cookbook":
				msg := UpdateCookbookMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "create_recipe":
				msg := CreateRecipeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "update_recipe":
				msg := UpdateRecipeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "enable_recipe":
				msg := EnableRecipeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "disable_recipe":
				msg := DisableRecipeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "execute_recipe":
				msg := ExecuteRecipeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "check_execution":
				msg := CheckExecutionMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "create_trade":
				msg := CreateTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "fulfill_trade":
				msg := FulfillTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "disable_trade":
				msg := DisableTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			case "enable_trade":
				msg := EnableTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = &msg, msg.Sender
			}
			msgs = append(msgs, newMsg)
		}
		t.WithFields(testing.Fields{
			"sender":   sender,
			"len_msgs": len(msgs),
			"tx_msgs":  inttest.AminoCodecFormatter(msgs),
			"msg_refs": step.MsgRefs,
		}).AddFields(inttest.GetLogFieldsFromMsgs(msgs)).Debug("debug log")
		txhash, err := inttest.SendMultiMsgTxWithNonce(t, msgs, sender, true)
		if err != nil {
			TxBroadcastErrorCheck(err, txhash, step, t)
			return
		}

		WaitForNextBlockWithErrorCheck(t)
		GetTxHandleResult(txhash, t)
	}
}

// CheckExecutionMsgFromRef collect check execution message from reference string
func CheckExecutionMsgFromRef(ref string, t *testing.T) types.MsgCheckExecution {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate execRef to execID
	newByteValue = UpdateExecID(newByteValue, t)

	var execType struct {
		ExecID        string
		PayToComplete bool
		Sender        string
	}
	err := json.Unmarshal(newByteValue, &execType)
	t.WithFields(testing.Fields{
		"bytes": string(newByteValue),
	}).MustNil(err, "error reading using json Unmarshaler")

	return types.NewMsgCheckExecution(
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
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &chkExecMsg, chkExecMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgCheckExecution{}).Type(), "MsgType should be accurate")
		resp := types.MsgCheckExecutionResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// FiatItemMsgFromRef collect check execution message from reference string
func FiatItemMsgFromRef(ref string, t *testing.T) types.MsgFiatItem {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate cookbook name to cookbook ID
	newByteValue = UpdateCBNameToID(newByteValue, t)

	var itemType types.Item
	err := inttest.GetJSONMarshaler().UnmarshalJSON(newByteValue, &itemType)
	t.WithFields(testing.Fields{
		"itemType": inttest.AminoCodecFormatter(itemType),
	}).MustNil(err, "error reading using JSONMarshaler")

	return types.NewMsgFiatItem(
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
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &itmMsg, itmMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgFiatItem{}).Type(), "MsgType should be accurate")
		resp := types.MsgFiatItemResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.ItemID != "", "item id shouldn't be empty")
	}
}

// SendItemsMsgFromRef is a function to collect SendItems from reference string
func SendItemsMsgFromRef(ref string, t *testing.T) types.MsgSendItems {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	newByteValue = UpdateReceiverKeyToAddress(newByteValue, t)

	var siType struct {
		Sender   string
		Receiver string
		ItemIDs  []string `json:"ItemIDs"`
	}

	err := json.Unmarshal(newByteValue, &siType)
	t.WithFields(testing.Fields{
		"siType":    inttest.AminoCodecFormatter(siType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")

	// translate itemNames to itemIDs
	ItemIDs := GetItemIDsFromNames(newByteValue, siType.Sender, false, false, t)

	return types.NewMsgSendItems(ItemIDs, siType.Sender, siType.Receiver)
}

// RunSendItems is a function to send items to another user
func RunSendItems(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		siMsg := SendItemsMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &siMsg, siMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgSendItems{}).Type(), "MsgType should be accurate")
		resp := types.MsgSendItemsResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// UpdateItemStringMsgFromRef is a function to collect UpdateItemStringMsg from reference string
func UpdateItemStringMsgFromRef(ref string, t *testing.T) types.MsgUpdateItemString {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate item name to item ID
	newByteValue = UpdateItemIDFromName(newByteValue, false, false, t)

	var sTypeMsg types.MsgUpdateItemString
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
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &sTypeMsg, sTypeMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgUpdateItemString{}).Type(), "MsgType should be accurate")
		resp := types.MsgUpdateItemStringResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
	}
}

// CreateCookbookMsgFromRef is a function to get create cookbook message from reference
func CreateCookbookMsgFromRef(ref string, t *testing.T) types.MsgCreateCookbook {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)

	var cbType types.Cookbook
	err := inttest.GetJSONMarshaler().UnmarshalJSON(newByteValue, &cbType)
	t.WithFields(testing.Fields{
		"cbType":    inttest.AminoCodecFormatter(cbType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")

	return types.NewMsgCreateCookbook(
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

		txhash, err := inttest.TestTxWithMsgWithNonce(t, &cbMsg, cbMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgCreateCookbook{}).Type(), "MsgType should be accurate")
		resp := types.MsgCreateCookbookResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.CookbookID != "", "coookbook id shouldn't be empty")
	}
}

// UpdateCookbookMsgFromRef is a function to get update cookbook message from reference
func UpdateCookbookMsgFromRef(ref string, t *testing.T) types.MsgUpdateCookbook {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)

	var cbType types.Cookbook
	err := inttest.GetJSONMarshaler().UnmarshalJSON(newByteValue, &cbType)
	t.WithFields(testing.Fields{
		"cbType":    inttest.AminoCodecFormatter(cbType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")

	return types.NewMsgUpdateCookbook(
		cbType.ID,
		cbType.Description,
		cbType.Developer,
		cbType.Version,
		cbType.SupportEmail,
		cbType.Sender,
	)
}

// RunUpdateCookbook is a function to update cookbook
func RunUpdateCookbook(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		cbMsg := UpdateCookbookMsgFromRef(step.ParamsRef, t)

		txhash, err := inttest.TestTxWithMsgWithNonce(t, &cbMsg, cbMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgUpdateCookbook{}).Type(), "MsgType should be accurate")
		resp := types.MsgUpdateCookbookResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
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
func CreateRecipeMsgFromRef(ref string, t *testing.T) types.MsgCreateRecipe {
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
	err := json.Unmarshal(newByteValue, &rcpTempl)
	t.WithFields(testing.Fields{
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")

	return types.NewMsgCreateRecipe(
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
		t.WithFields(testing.Fields{
			"parsed_recipe": string(inttest.GetAminoCdc().MustMarshalJSON(rcpMsg)),
		}).Info("recipe info")

		txhash, err := inttest.TestTxWithMsgWithNonce(t, &rcpMsg, rcpMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgCreateRecipe{}).Type(), "MsgType should be accurate")
		resp := types.MsgCreateRecipeResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.RecipeID != "", "recipe id shouldn't be empty")
	}
}

// UpdateRecipeMsgFromRef is a function to get update recipe message from reference
func UpdateRecipeMsgFromRef(ref string, t *testing.T) types.MsgUpdateRecipe {
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
	err := json.Unmarshal(newByteValue, &rcpTempl)
	t.WithFields(testing.Fields{
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using json.Unmarshal")

	return types.NewMsgUpdateRecipe(
		rcpTempl.ID,
		rcpTempl.Name,
		rcpTempl.CookbookID,
		rcpTempl.Description,
		rcpTempl.CoinInputs,
		itemInputs,
		entries,
		rcpTempl.Outputs,
		rcpTempl.BlockInterval,
		rcpTempl.Sender,
	)
}

// RunUpdateRecipe is a function to update recipe
func RunUpdateRecipe(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		rcpMsg := UpdateRecipeMsgFromRef(step.ParamsRef, t)

		txhash, err := inttest.TestTxWithMsgWithNonce(t, &rcpMsg, rcpMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgUpdateRecipe{}).Type(), "MsgType should be accurate")
		resp := types.MsgUpdateRecipeResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.RecipeID != "", "recipe id shouldn't be empty")
	}
}

// EnableRecipeMsgFromRef is a function to get enable recipe message from reference
func EnableRecipeMsgFromRef(ref string, t *testing.T) types.MsgEnableRecipe {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate recipe name to recipe id
	newByteValue = UpdateRecipeName(newByteValue, t)

	var recipeType struct {
		RecipeID string
		Sender   string
	}

	err := json.Unmarshal(newByteValue, &recipeType)
	t.WithFields(testing.Fields{
		"rcpTempl":  inttest.AminoCodecFormatter(recipeType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")

	return types.NewMsgEnableRecipe(recipeType.RecipeID, recipeType.Sender)
}

// RunEnableRecipe is a function to enable recipe
func RunEnableRecipe(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		rcpMsg := EnableRecipeMsgFromRef(step.ParamsRef, t)

		txhash, err := inttest.TestTxWithMsgWithNonce(t, &rcpMsg, rcpMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgEnableRecipe{}).Type(), "MsgType should be accurate")
		resp := types.MsgEnableRecipeResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// DisableRecipeMsgFromRef is a function to get disable recipe message from reference
func DisableRecipeMsgFromRef(ref string, t *testing.T) types.MsgDisableRecipe {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate recipe name to recipe id
	newByteValue = UpdateRecipeName(newByteValue, t)

	var recipeType struct {
		RecipeID string
		Sender   string
	}

	err := json.Unmarshal(newByteValue, &recipeType)
	t.WithFields(testing.Fields{
		"rcpTempl":  inttest.AminoCodecFormatter(recipeType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetAminoCdc")

	return types.NewMsgDisableRecipe(recipeType.RecipeID, recipeType.Sender)
}

// RunDisableRecipe is a function to disable recipe
func RunDisableRecipe(step FixtureStep, t *testing.T) {
	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		rcpMsg := DisableRecipeMsgFromRef(step.ParamsRef, t)

		txhash, err := inttest.TestTxWithMsgWithNonce(t, &rcpMsg, rcpMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgDisableRecipe{}).Type(), "MsgType should be accurate")
		resp := types.MsgDisableRecipeResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// ExecuteRecipeMsgFromRef collect execute recipe msg from reference string
func ExecuteRecipeMsgFromRef(ref string, t *testing.T) types.MsgExecuteRecipe {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate recipe name to recipe id
	newByteValue = UpdateRecipeName(newByteValue, t)

	var execType struct {
		RecipeID string
		Sender   string
		ItemIDs  []string `json:"ItemIDs"`
	}

	err := json.Unmarshal(newByteValue, &execType)
	t.WithFields(testing.Fields{
		"execType":  inttest.AminoCodecFormatter(execType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")
	// translate itemNames to itemIDs
	ItemIDs := GetItemIDsFromNames(newByteValue, execType.Sender, false, false, t)

	return types.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
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
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &execMsg, execMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgExecuteRecipe{}).Type(), "MsgType should be accurate")
		resp := types.MsgExecuteRecipeResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)

		if resp.Message == "scheduled the recipe" { // delayed execution
			var scheduleRes handlers.ExecuteRecipeScheduleOutput

			err := json.Unmarshal(resp.Output, &scheduleRes)
			t.WithFields(testing.Fields{
				"response_output": string(resp.Output),
			}).MustNil(err, "error decoding raw json")
			execIDRWMutex.Lock()
			execIDs[step.ID] = scheduleRes.ExecID
			execIDRWMutex.Unlock()
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
func CreateTradeMsgFromRef(ref string, t *testing.T) types.MsgCreateTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// get item inputs from fileNames
	tradeItemInputs := GetTradeItemInputsFromBytes(newByteValue, t)
	var trdType types.Trade
	err := json.Unmarshal(newByteValue, &trdType)
	t.WithFields(testing.Fields{
		"trdType":   inttest.AminoCodecFormatter(trdType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")

	// get ItemOutputs from ItemOutputNames
	itemOutputs := GetItemOutputsFromBytes(newByteValue, trdType.Sender, t)

	return types.NewMsgCreateTrade(
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
		}).AddFields(inttest.GetLogFieldsFromMsgs([]sdk.Msg{&createTrd})).Debug("createTrd")
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &createTrd, createTrd.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgCreateTrade{}).Type(), "MsgType should be accurate")
		resp := types.MsgCreateTradeResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(resp.TradeID != "", "trade id shouldn't be empty")
	}
}

// FulfillTradeMsgFromRef collect fulfill trade message from reference string
func FulfillTradeMsgFromRef(ref string, t *testing.T) types.MsgFulfillTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate extra info to trade id
	newByteValue = UpdateTradeExtraInfoToID(newByteValue, t)

	var trdType struct {
		TradeID string
		Sender  sdk.AccAddress
		ItemIDs []string `json:"ItemIDs"`
	}

	err := json.Unmarshal(newByteValue, &trdType)
	t.WithFields(testing.Fields{
		"trdType":   inttest.AminoCodecFormatter(trdType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")
	// translate itemNames to itemIDs
	ItemIDs := GetItemIDsFromNames(newByteValue, trdType.Sender.String(), false, false, t)

	return types.NewMsgFulfillTrade(trdType.TradeID, trdType.Sender.String(), ItemIDs)
}

// RunFulfillTrade is a function to fulfill trade
func RunFulfillTrade(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		ffTrdMsg := FulfillTradeMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &ffTrdMsg, ffTrdMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgFulfillTrade{}).Type(), "MsgType should be accurate")
		resp := types.MsgFulfillTradeResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// DisableTradeMsgFromRef collect disable trade msg from reference string
func DisableTradeMsgFromRef(ref string, t *testing.T) types.MsgDisableTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate extra info to trade id
	newByteValue = UpdateTradeExtraInfoToID(newByteValue, t)

	var trdType struct {
		TradeID string
		Sender  string
	}

	err := json.Unmarshal(newByteValue, &trdType)
	t.WithFields(testing.Fields{
		"trdType":   inttest.AminoCodecFormatter(trdType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")

	return types.NewMsgDisableTrade(trdType.TradeID, trdType.Sender)
}

// RunDisableTrade is a function to disable trade
func RunDisableTrade(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		dsTrdMsg := DisableTradeMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &dsTrdMsg, dsTrdMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgDisableTrade{}).Type(), "MsgType should be accurate")
		resp := types.MsgDisableTradeResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}

// EnableTradeMsgFromRef collect enable trade msg from reference string
func EnableTradeMsgFromRef(ref string, t *testing.T) types.MsgEnableTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate extra info to trade id
	newByteValue = UpdateTradeExtraInfoToID(newByteValue, t)

	var trdType struct {
		TradeID string
		Sender  string
	}

	err := json.Unmarshal(newByteValue, &trdType)
	t.WithFields(testing.Fields{
		"trdType":   inttest.AminoCodecFormatter(trdType),
		"new_bytes": string(newByteValue),
	}).MustNil(err, "error reading using GetJSONMarshaler")

	return types.NewMsgEnableTrade(trdType.TradeID, trdType.Sender)
}

// RunEnableTrade is a function to enable trade
func RunEnableTrade(step FixtureStep, t *testing.T) {

	if FixtureTestOpts.VerifyOnly {
		return
	}
	if step.ParamsRef != "" {
		dsTrdMsg := EnableTradeMsgFromRef(step.ParamsRef, t)
		txhash, err := inttest.TestTxWithMsgWithNonce(t, &dsTrdMsg, dsTrdMsg.Sender, true)
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
		txMsgData := &sdk.TxMsgData{
			Data: make([]*sdk.MsgData, 0, 1),
		}
		err = proto.Unmarshal(txHandleResBytes, txMsgData)
		TxResultDecodingErrorCheck(err, txhash, t)
		t.MustTrue(len(txMsgData.Data) == 1, "number of msgs should be 1")
		t.MustTrue(txMsgData.Data[0].MsgType == (types.MsgEnableTrade{}).Type(), "MsgType should be accurate")
		resp := types.MsgEnableTradeResponse{}
		err = proto.Unmarshal(txMsgData.Data[0].Data, &resp)
		TxResultDecodingErrorCheck(err, txhash, t)
		TxResultStatusMessageCheck(resp.Status, resp.Message, txhash, step, t)
	}
}
