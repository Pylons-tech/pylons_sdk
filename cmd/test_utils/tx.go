package inttest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Pylons-tech/pylons_sdk/app"
	testing "github.com/Pylons-tech/pylons_sdk/cmd/evtesting"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var nonceMux sync.Mutex

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GenTxWithMsg is a function to generate transaction from msg
func GenTxWithMsg(messages []sdk.Msg) (authsigning.Tx, error) {
	var err error
	for i, msg := range messages {
		if err = msg.ValidateBasic(); err != nil {
			return nil, fmt.Errorf("%dth msg does not pass basic validation for %s", i, err.Error())
		}
	}

	viper.Set("keyring-backend", "test")
	viper.Set("chain-id", "pylonschain")

	txBldr := app.MakeEncodingConfig().TxConfig.NewTxBuilder()
	err = txBldr.SetMsgs(messages...)
	if err != nil {
		return nil, err
	}

	txBldr.SetGasLimit(10000000)
	return txBldr.GetTx(), nil
}

func broadcastTxFile(signedTxFile string, maxRetry int, t *testing.T) (string, error) {
	if len(CLIOpts.RestEndpoint) == 0 { // broadcast using cli
		// pylonsd tx broadcast signedCreateCookbookTx.json
		txBroadcastArgs := []string{"tx", "broadcast", signedTxFile}
		output, logstr, err := RunPylonsd(txBroadcastArgs, "")
		// output2, logstr2, err := RunPylonsd([]string{"query", "account", "cosmos10xgn8t2auxskrf2qjcht0hwq2h5chnrpx87dus"}, "")
		// t.WithFields(testing.Fields{
		// 	"query_account": logstr2,
		// 	"output2":       string(output2),
		// }).Debug("debug log")

		t.WithFields(testing.Fields{
			"broadcast_args":   txBroadcastArgs,
			"broadcast_output": string(output),
			"broadcast_log":    logstr,
		}).MustNil(err, "error running pylonsd broadcast command")
		txResponse := sdk.TxResponse{}

		err = GetAminoCdc().UnmarshalJSON(output, &txResponse)
		// This can happen when "pylonsd config output json" is not set or when real issue is available
		t.WithFields(testing.Fields{
			"broadcast_output":  string(output),
			"possible_solution": "You can set cli config output as json to solve this issue",
		}).MustNil(err, "error decoding transaction broadcast result")

		if err != nil {
			return txResponse.TxHash, err
		}

		if txResponse.Code == sdkerrors.ErrUnauthorized.ABCICode() &&
			strings.Contains(txResponse.RawLog, "signature verification failed") && maxRetry > 0 {
			t.WithFields(testing.Fields{
				"log":       logstr,
				"output":    string(output),
				"max_retry": maxRetry,
			}).Info("rebroadcasting after 1s...")
			time.Sleep(1 * time.Second)
			return broadcastTxFile(signedTxFile, maxRetry-1, t)
		}
		if txResponse.Code != 0 {
			return txResponse.TxHash, errors.New(txResponse.RawLog)
		}
		t.WithFields(testing.Fields{
			"txhash": txResponse.TxHash,
		}).MustTrue(len(txResponse.TxHash) == 64, "txhash length should have length of 64")
		return txResponse.TxHash, nil
	}
	// broadcast using rest endpoint
	signedTx := ReadFile(signedTxFile, t)
	postBodyJSON := make(map[string]interface{})

	err := json.Unmarshal(signedTx, &postBodyJSON)
	t.WithFields(testing.Fields{
		"signed_tx": string(signedTx),
	}).MustNil(err, "error decoding raw json")

	postBodyJSON["tx"] = postBodyJSON["value"]
	postBodyJSON["value"] = nil
	postBodyJSON["mode"] = "sync"
	postBody, err := json.Marshal(postBodyJSON)

	if err != nil {
		t.MustNil(err, "fatal log")
		return "", err
	}
	resp, err := http.Post(CLIOpts.RestEndpoint+"/txs", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.MustNil(err, "fatal log")
		return "", err
	}

	var result map[string]string

	err = json.NewDecoder(resp.Body).Decode(&result)
	t.MustNil(err, "error decoding raw json")
	defer resp.Body.Close()
	t.WithFields(testing.Fields{
		"get_pylons_api_response": result,
	}).Info("info log")
	t.WithFields(testing.Fields{
		"txhash": result["txhash"],
	}).MustTrue(len(result["txhash"]) == 64, "txhash length should have length of 64")
	return result["txhash"], nil
}

// TestTxWithMsg is a function to send transaction with message
func TestTxWithMsg(t *testing.T, msgValue sdk.Msg, signer string) string {
	tmpDir, err := ioutil.TempDir("", "pylons")
	if err != nil {
		panic(err.Error())
	}
	rawTxFile := filepath.Join(tmpDir, "raw_tx.json")
	signedTxFile := filepath.Join(tmpDir, "signed_tx.json")

	txModel, err := GenTxWithMsg([]sdk.Msg{msgValue})
	t.MustNil(err, "error while building transaction model from messages")
	output, err := GetAminoCdc().MarshalJSON(txModel)
	t.MustNil(err, "error encoding transaction model")

	err = ioutil.WriteFile(rawTxFile, output, 0644)
	if err != nil {
		t.WithFields(testing.Fields{
			"tx_model_json": string(output),
		}).MustNil(err, "error writing raw transaction")
		return ""
	}

	// pylonsd tx sign raw_tx.json --from eugen --chain-id pylonschain > signed_tx.json
	txSignArgs := []string{"tx", "sign", rawTxFile,
		"--from", signer,
		"--chain-id", "pylonschain",
	}
	output, _, err = RunPylonsd(txSignArgs, "")
	if err != nil {
		t.WithFields(testing.Fields{
			"signed_tx_json": string(output),
		}).MustNil(err, "error signing transaction")
		return ""
	}

	err = ioutil.WriteFile(signedTxFile, output, 0644)
	if err != nil {
		t.MustNil(err, "error writing signed transaction")
		return ""
	}

	txhash, err := broadcastTxFile(signedTxFile, GetMaxBroadcastRetry(), t)
	if err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Error("transaction broadcast failure")
		return ""
	}

	CleanFile(rawTxFile, t)
	CleanFile(signedTxFile, t)

	return txhash
}

// SendMultiMsgTxWithNonce is a function to send multiple messages in one transaction
func SendMultiMsgTxWithNonce(t *testing.T, msgs []sdk.Msg, signer string, isBech32Addr bool) (string, error) {
	t.WithFields(testing.Fields{
		"action":    "func_start",
		"signer":    signer,
		"is_bech32": isBech32Addr,
	}).
		AddFields(GetLogFieldsFromMsgs(msgs)).
		AddFields(log.Fields{
			"tx_msgs": AminoCodecFormatter(msgs),
		}).
		SetFieldsOrder(testing.SortCustomKey, []string{"action", "signer", "is_bech32"}).
		Debug("debug log")

	if len(msgs) == 0 {
		return "msgs validation error", errors.New("length of msgs shouldn't be zero")
	}
	t.Trace("tx_with_nonce.step.A")
	tmpDir, err := ioutil.TempDir("", "pylons")
	if err != nil {
		return "error creating pylons directory on temp folder", err
	}
	t.Trace("tx_with_nonce.step.B")
	nonceRootDir := "./"
	nonceFile := filepath.Join(nonceRootDir, "nonce.json")
	if !isBech32Addr {
		signer = GetAccountAddr(signer, t)
	}

	t.Trace("tx_with_nonce.step.C")
	accInfo := GetAccountInfoFromAddr(signer, t)
	nonce := accInfo.GetSequence()

	nonceMap := make(map[string]uint64)

	nonceMux.Lock()
	defer nonceMux.Unlock()

	if fileExists(nonceFile) {
		nonceBytes := ReadFile(nonceFile, t)
		err := json.Unmarshal(nonceBytes, &nonceMap)
		if err != nil {
			return "error unmarshaling nonce map", err
		}
		if existNonce, ok := nonceMap[signer]; ok {
			nonce = existNonce
		}
	}
	t.Trace("tx_with_nonce.step.D")

	txModel, err := GenTxWithMsg(msgs)
	if err != nil {
		return "error generating transaction with messages", err
	}
	t.Trace("tx_with_nonce.step.E")
	output, err := GetAminoCdc().MarshalJSON(txModel)
	if err != nil {
		return "error marshaling transaction into json", err
	}

	t.Trace("tx_with_nonce.step.F")
	rawTxFile := filepath.Join(tmpDir, "raw_tx_"+strconv.FormatUint(nonce, 10)+".json")
	signedTxFile := filepath.Join(tmpDir, "signed_tx_"+strconv.FormatUint(nonce, 10)+".json")
	err = ioutil.WriteFile(rawTxFile, output, 0644)
	if err != nil {
		t.WithFields(testing.Fields{
			"tx_model_json": string(output),
		}).MustNil(err, "error writing raw transaction")
		return "error writing raw transaction", err
	}

	t.Trace("tx_with_nonce.step.G")
	// pylonsd tx sign sample_transaction.json --account-number 2 --sequence 10 --offline --from eugen
	txSignArgs := []string{"tx", "sign", rawTxFile,
		"--from", signer,
		"--offline",
		"--chain-id", "pylonschain",
		"--sequence", strconv.FormatUint(nonce, 10),
		"--account-number", strconv.FormatUint(accInfo.GetAccountNumber(), 10),
	}
	output, _, err = RunPylonsd(txSignArgs, "")
	// output, logstr, err := RunPylonsd(txSignArgs, "")
	// t.WithFields(testing.Fields{
	// 	"error": err,
	// 	"log": logstr,
	// })("TX sign result")
	if err != nil {
		return "error signing transaction", fmt.Errorf("%s; %s", err.Error(), string(output))
	}
	t.Trace("tx_with_nonce.step.H")

	err = ioutil.WriteFile(signedTxFile, output, 0644)
	if err != nil {
		return "error writing signed transaction", err
	}

	t.Trace("tx_with_nonce.step.I")

	txhash, err := broadcastTxFile(signedTxFile, GetMaxBroadcastRetry(), t)
	t.WithFields(testing.Fields{
		"sequence":       strconv.FormatUint(nonce, 10),
		"account-number": strconv.FormatUint(accInfo.GetAccountNumber(), 10),
		"max-retry":      GetMaxBroadcastRetry(),
		"tx_msgs":        AminoCodecFormatter(msgs),
		"error":          err,
	}).Debug("transaction broadcast debug")
	if err != nil {
		return "error broadcasting tx file", err
	}
	// increase nonce file
	t.Trace("tx_with_nonce.step.J")
	nonceMap[signer] = nonce + 1
	nonceOutput, err := json.Marshal(nonceMap)
	if err != nil {
		return "error marshaling nonceMap", err
	}
	t.Trace("tx_with_nonce.step.K")
	err = ioutil.WriteFile(nonceFile, nonceOutput, 0644)
	if err != nil {
		exPath := ""
		ex, err := os.Executable()
		if err == nil {
			exPath = filepath.Dir(ex)
		} else {
			exPath = "UNDEFINED_EXPATH"
		}
		return fmt.Sprintf("error writing nonce output file at %s", exPath), err
	}
	t.Trace("tx_with_nonce.step.L")

	CleanFile(rawTxFile, t)
	CleanFile(signedTxFile, t)

	t.WithFields(testing.Fields{
		"action":    "func_end",
		"txhash":    txhash,
		"signer":    signer,
		"is_bech32": isBech32Addr,
	}).
		AddFields(GetLogFieldsFromMsgs(msgs)).
		AddFields(log.Fields{
			"tx_msgs": AminoCodecFormatter(msgs),
		}).
		SetFieldsOrder(testing.SortCustomKey, []string{"action", "txhash", "signer", "is_bech32"}).
		Debug("debug log")
	return txhash, nil
}

// TestTxWithMsgWithNonce is a function to send transaction with message and nonce
func TestTxWithMsgWithNonce(t *testing.T, msgValue sdk.Msg, signer string, isBech32Addr bool) (string, error) {
	output, err := SendMultiMsgTxWithNonce(t, []sdk.Msg{msgValue}, signer, isBech32Addr)
	if err != nil {
		// output is txhash if it's a success transaction, if fail, it's output log
		t.WithFields(testing.Fields{
			"output": output,
			"error":  err,
			"func":   "TestTxWithMsgWithNonce",
		}).Error("error log")
	}
	return output, err
}
