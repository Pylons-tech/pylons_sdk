package inttest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"
	log "github.com/sirupsen/logrus"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"

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
func GenTxWithMsg(messages []sdk.Msg) (auth.StdTx, error) {
	var err error
	cdc := GetAminoCdc()
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	viper.Set("keyring-backend", "test")
	viper.Set("chain-id", "pylonschain")

	txBldr := auth.NewTxBuilderFromCLI(&bytes.Buffer{}).WithTxEncoder(utils.GetTxEncoder(cdc)).WithChainID("pylonschain")
	if txBldr.SimulateAndExecute() {
		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, messages)
		if err != nil {
			return auth.StdTx{}, err
		}
	}

	stdSignMsg, err := txBldr.BuildSignMsg(messages)
	if err != nil {
		return auth.StdTx{}, err
	}
	stdSignMsg.Fee.Gas = 400000

	return auth.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo), nil
}

func broadcastTxFile(signedTxFile string, maxRetry int, t *testing.T) (string, error) {
	if len(CLIOpts.RestEndpoint) == 0 { // broadcast using cli
		// pylonscli tx broadcast signedCreateCookbookTx.json
		txBroadcastArgs := []string{"tx", "broadcast", signedTxFile}
		output, logstr, err := RunPylonsCli(txBroadcastArgs, "")
		// output2, logstr2, err := RunPylonsCli([]string{"query", "account", "cosmos10xgn8t2auxskrf2qjcht0hwq2h5chnrpx87dus"}, "")
		// t.WithFields(testing.Fields{
		// 	"broadcast_log": logstr,
		// 	"query_account": logstr2,
		// 	"output2":       string(output2),
		// }).Debug("debug log")

		t.MustNil(err)
		txResponse := sdk.TxResponse{}

		err = GetAminoCdc().UnmarshalJSON(output, &txResponse)
		// This can happen when "pylonscli config output json" is not set or when real issue is available
		if err != nil {
			t.WithFields(testing.Fields{
				"broadcast_output": string(output),
				"error":            err,
			}).Fatal("error in broadcasting signed transaction output")
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
		t.MustTrue(len(txResponse.TxHash) == 64)
		return txResponse.TxHash, nil
	}
	// broadcast using rest endpoint
	signedTx := ReadFile(signedTxFile, t)
	postBodyJSON := make(map[string]interface{})

	err := json.Unmarshal(signedTx, &postBodyJSON)
	t.MustNil(err)

	postBodyJSON["tx"] = postBodyJSON["value"]
	postBodyJSON["value"] = nil
	postBodyJSON["mode"] = "sync"
	postBody, err := json.Marshal(postBodyJSON)

	if err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("fatal log")
	}
	resp, err := http.Post(CLIOpts.RestEndpoint+"/txs", "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("fatal log")
	}

	var result map[string]string

	err = json.NewDecoder(resp.Body).Decode(&result)
	t.MustNil(err)
	defer resp.Body.Close()
	t.WithFields(testing.Fields{
		"get_pylons_api_response": result,
	}).Info("info log")
	t.MustTrue(len(result["txhash"]) == 64)
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
	t.MustNil(err)
	output, err := GetAminoCdc().MarshalJSON(txModel)
	t.MustNil(err)

	err = ioutil.WriteFile(rawTxFile, output, 0644)
	if err != nil {
		t.WithFields(testing.Fields{
			"tx_model_json": string(output),
			"error":         err,
		}).Fatal("error writing raw transaction")
	}

	// pylonscli tx sign raw_tx.json --from eugen --chain-id pylonschain > signed_tx.json
	txSignArgs := []string{"tx", "sign", rawTxFile,
		"--from", signer,
		"--chain-id", "pylonschain",
	}
	output, _, err = RunPylonsCli(txSignArgs, "")
	if err != nil {
		t.WithFields(testing.Fields{
			"signed_tx_json": string(output),
			"error":          err,
		}).Fatal("error signing transaction")
	}

	err = ioutil.WriteFile(signedTxFile, output, 0644)
	if err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("error writing signed transaction")
	}

	txhash, err := broadcastTxFile(signedTxFile, GetMaxBroadcastRetry(), t)
	if err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("broadcasting failure after maxRetry limitation")
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

	if fileExists(nonceFile) {
		nonceBytes := ReadFile(nonceFile, t)
		err := json.Unmarshal(nonceBytes, &nonceMap)
		if err != nil {
			return "error unmarshaling nonce map", err
		}
		nonce = nonceMap[signer]
	}
	t.Trace("tx_with_nonce.step.K")

	nonceMap[signer] = nonce + 1
	nonceOutput, err := json.Marshal(nonceMap)
	if err != nil {
		return "error marshaling nonceMap", err
	}
	t.Trace("tx_with_nonce.step.K")
	err = ioutil.WriteFile(nonceFile, nonceOutput, 0644)
	if err != nil {
		return "error writing nonce output file", err
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
			"error":         err,
		}).Fatal("error writing raw transaction")
	}

	t.Trace("tx_with_nonce.step.G")
	// pylonscli tx sign sample_transaction.json --account-number 2 --sequence 10 --offline --from eugen
	txSignArgs := []string{"tx", "sign", rawTxFile,
		"--from", signer,
		"--offline",
		"--chain-id", "pylonschain",
		"--sequence", strconv.FormatUint(nonce, 10),
		"--account-number", strconv.FormatUint(accInfo.GetAccountNumber(), 10),
	}
	output, _, err = RunPylonsCli(txSignArgs, "")
	// output, logstr, err := RunPylonsCli(txSignArgs, "")
	// t.WithFields(testing.Fields{
	// 	"error": err,
	// 	"log": logstr,
	// })("TX sign result")
	if err != nil {
		return "error signing transaction", err
	}
	t.Trace("tx_with_nonce.step.H")

	err = ioutil.WriteFile(signedTxFile, output, 0644)
	if err != nil {
		return "error writing signed transaction", err
	}

	t.Trace("tx_with_nonce.step.I")
	nonceMux.Unlock()

	txhash, err := broadcastTxFile(signedTxFile, GetMaxBroadcastRetry(), t)
	if err != nil {
		t.WithFields(testing.Fields{
			"error": err,
		}).Fatal("broadcasting failure after maxRetry limitation")
	}
	t.Trace("tx_with_nonce.step.J")

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
func TestTxWithMsgWithNonce(t *testing.T, msgValue sdk.Msg, signer string, isBech32Addr bool) string {
	txhash, err := SendMultiMsgTxWithNonce(t, []sdk.Msg{msgValue}, signer, isBech32Addr)
	if err != nil {
		t.WithFields(testing.Fields{
			"txhash": txhash,
			"error":  err,
			"func":   "TestTxWithMsgWithNonce",
		}).Fatal("fatal log")
	}
	return txhash
}
