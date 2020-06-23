package fixturetest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	originT "testing"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"
	inttest "github.com/Pylons-tech/pylons_sdk/cmd/test"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FixtureStep struct describes what should be done in one fixture testcase
type FixtureStep struct {
	ID       string `json:"ID"`
	RunAfter struct {
		PreCondition []string `json:"precondition"`
		BlockWait    int64    `json:"blockWait"`
	} `json:"runAfter"`
	Action    string `json:"action"`
	ParamsRef string `json:"paramsRef"`
	MsgRefs   []struct {
		Action    string `json:"action"`
		ParamsRef string `json:"paramsRef"`
	} `json:"msgRefs"`
	Output struct {
		TxResult struct {
			Status         string `json:"status"`
			Message        string `json:"message"`
			ErrorLog       string `json:"errLog"`
			BroadcastError string `json:"broadcastError"`
		} `json:"txResult"`
		Property []struct {
			Owner          string   `json:"owner"`
			ShouldNotExist bool     `json:"shouldNotExist"`
			Cookbooks      []string `json:"cookbooks"`
			Recipes        []string `json:"recipes"`
			Items          []struct {
				StringKeys   []string                     `json:"stringKeys"`
				StringValues map[string]string            `json:"stringValues"`
				DblKeys      []string                     `json:"dblKeys"`
				DblValues    map[string]types.FloatString `json:"dblValues"`
				LongKeys     []string                     `json:"longKeys"`
				LongValues   map[string]int               `json:"longValues"`
			} `json:"items"`
			Coins []struct {
				Coin   string `json:"denom"`
				Amount int64  `json:"amount"`
			} `json:"coins"`
		} `json:"property"`
	} `json:"output"`
}

// TestOptions is options struct to manage test options
type TestOptions struct {
	IsParallel        bool
	CreateNewCookbook bool
}

// FixtureTestOpts is a variable to have fixture test options
var FixtureTestOpts = TestOptions{
	IsParallel: true,
}

// CheckItemWithStringKeys checks if string keys are all available
func CheckItemWithStringKeys(item types.Item, stringKeys []string) bool {
	for _, sK := range stringKeys {
		keyExist := false
		for _, sKV := range item.Strings {
			if sK == sKV.Key {
				keyExist = true
			}
		}
		if !keyExist {
			return false
		}
	}
	return true
}

// CheckItemWithStringValues checks if string value/key set are all available
func CheckItemWithStringValues(item types.Item, stringValues map[string]string) bool {
	for sK, sV := range stringValues {
		keyExist := false
		for _, sKV := range item.Strings {
			if sK == sKV.Key && sV == sKV.Value {
				keyExist = true
			}
		}
		if !keyExist {
			return false
		}
	}
	return true
}

// CheckItemWithDblKeys checks if double keys are all available
func CheckItemWithDblKeys(item types.Item, dblKeys []string) bool {
	for _, sK := range dblKeys {
		keyExist := false
		for _, sKV := range item.Doubles {
			if sK == sKV.Key {
				keyExist = true
			}
		}
		if !keyExist {
			return false
		}
	}
	return true
}

// CheckItemWithDblValues checks if double key/values are all available
func CheckItemWithDblValues(item types.Item, dblValues map[string]types.FloatString) bool {
	for sK, sV := range dblValues {
		keyExist := false
		for _, sKV := range item.Doubles {
			if sK == sKV.Key && sV.Float() == sKV.Value.Float() {
				keyExist = true
			}
		}
		if !keyExist {
			return false
		}
	}
	return true
}

// CheckItemWithLongKeys checks if long keys are all available
func CheckItemWithLongKeys(item types.Item, longKeys []string) bool {
	for _, sK := range longKeys {
		keyExist := false
		for _, sKV := range item.Longs {
			if sK == sKV.Key {
				keyExist = true
			}
		}
		if !keyExist {
			return false
		}
	}
	return true
}

// CheckItemWithLongValues checks if long key/values are all available
func CheckItemWithLongValues(item types.Item, longValues map[string]int) bool {
	for sK, sV := range longValues {
		keyExist := false
		for _, sKV := range item.Longs {
			if sK == sKV.Key && sV == sKV.Value {
				keyExist = true
			}
		}
		if !keyExist {
			return false
		}
	}
	return true
}

// CheckErrorOnTxFromTxHash validate if there's an error on transaction
func CheckErrorOnTxFromTxHash(txhash string, t *testing.T) {
	hmrErrMsg := inttest.GetHumanReadableErrorFromTxHash(txhash, t)
	if len(hmrErrMsg) > 0 {
		t.WithFields(testing.Fields{
			"txhash":   txhash,
			"tx_error": hmrErrMsg,
		}).Fatal("tx_error exist")
	}
}

// PropertyExistCheck function check if an account has required property that needs to be available
func PropertyExistCheck(step FixtureStep, t *testing.T) {
	for _, pCheck := range step.Output.Property {
		shouldNotExist := pCheck.ShouldNotExist
		var pOwnerAddr string
		if len(pCheck.Owner) == 0 {
			pOwnerAddr = ""
		} else {
			ownerAddrIndex, err := strconv.Atoi(strings.TrimLeft(pCheck.Owner, "account"))
			t.MustNil(err, "temp account name doesn't match to the account args")
			ownerAddrIndex--

			t.MustTrue(ownerAddrIndex < len(accountNames), "temp account name doesn't match to the account args")
			// pOwnerAddr := accountNames[ownerAddrIndex]

			pOwnerAddr = inttest.GetAccountAddr(accountNames[ownerAddrIndex], t)
		}
		if len(pCheck.Cookbooks) > 0 {
			for _, cbName := range pCheck.Cookbooks {
				_, exist, err := inttest.GetCookbookIDFromName(cbName, pOwnerAddr)
				if err != nil {
					t.WithFields(testing.Fields{
						"error": err,
					}).Fatal("error checking cookbook exist")
				}
				if !shouldNotExist {
					if exist {
						t.WithFields(testing.Fields{
							"cookbook_name": cbName,
						}).Info("cookbook exist, ok")
					} else {
						t.WithFields(testing.Fields{
							"cookbook_name": cbName,
						}).Fatal("cookbook does not exist, but should exist")
					}
				} else {
					if exist {
						t.WithFields(testing.Fields{
							"cookbook_name": cbName,
						}).Fatal("cookbook exist, but shouldn't exist")
					} else {
						t.WithFields(testing.Fields{
							"cookbook_name": cbName,
						}).Info("cookbook does not exist as expected, ok")
					}
				}
			}
		}
		if len(pCheck.Recipes) > 0 {
			for _, rcpName := range pCheck.Recipes {
				guid, err := inttest.GetRecipeGUIDFromName(rcpName, pOwnerAddr)
				if err != nil {
					t.WithFields(testing.Fields{
						"error": err,
					}).Fatal("error checking if recipe already exist")
				}

				if !shouldNotExist {
					if len(guid) > 0 {
						t.WithFields(testing.Fields{
							"recipe_name": rcpName,
						}).Info("recipe exist, ok")
					} else {
						t.WithFields(testing.Fields{
							"recipe_name": rcpName,
						}).Fatal("recipe with does not exist, but should exist")
					}
				} else {
					if len(guid) > 0 {
						t.WithFields(testing.Fields{
							"recipe_name": rcpName,
						}).Fatal("recipe exist but shouldn't exist")
					} else {
						t.WithFields(testing.Fields{
							"recipe_name": rcpName,
						}).Info("recipe does not exist as expected, ok")
					}
				}
			}
		}
		if len(pCheck.Items) > 0 {
			for _, itemCheck := range pCheck.Items {
				fitItemExist := false
				// t.WithFields(testing.Fields{
				// 	// "id": idx,
				// 	"item_spec": itemCheck,
				// }).Info("checking item")
				items, err := inttest.ListItemsViaCLI(pOwnerAddr)
				if err != nil {
					t.WithFields(testing.Fields{
						"error": err,
					}).Fatal("error listing items")
				}
				for _, item := range items {
					if CheckItemWithStringKeys(item, itemCheck.StringKeys) &&
						CheckItemWithStringValues(item, itemCheck.StringValues) &&
						CheckItemWithDblKeys(item, itemCheck.DblKeys) &&
						CheckItemWithDblValues(item, itemCheck.DblValues) &&
						CheckItemWithLongKeys(item, itemCheck.LongKeys) &&
						CheckItemWithLongValues(item, itemCheck.LongValues) {
						fitItemExist = true
					}
				}

				if !shouldNotExist {
					if fitItemExist {
						t.WithFields(testing.Fields{
							"owner_address": pOwnerAddr,
							"item_spec":     inttest.JSONFormatter(itemCheck),
						}).Info("checked item existence")
					} else {
						t.WithFields(testing.Fields{
							"owner_address": pOwnerAddr,
							"item_spec":     inttest.JSONFormatter(itemCheck),
						}).Fatal("no item exist which fit item spec")
					}
				} else {
					if fitItemExist {
						t.WithFields(testing.Fields{
							"owner_address": pOwnerAddr,
							"item_spec":     inttest.JSONFormatter(itemCheck),
						}).Fatal("item exist but shouldn't exist")
					} else {
						t.WithFields(testing.Fields{
							"owner_address": pOwnerAddr,
							"item_spec":     inttest.JSONFormatter(itemCheck),
						}).Info("item does not exist as expected, ok")
					}
				}
			}
		}
		if len(pCheck.Coins) > 0 {
			for _, coinCheck := range pCheck.Coins {
				accInfo := inttest.GetAccountInfoFromName(pCheck.Owner, t)
				// TODO should we have the case of using GTE, LTE, GT or LT ?
				t.MustTrue(accInfo.Coins.AmountOf(coinCheck.Coin).Equal(sdk.NewInt(coinCheck.Amount)), "account balance is incorrect")
			}
		}
	}
}

// ProcessSingleFixtureQueueItem executes a fixture queue item
func ProcessSingleFixtureQueueItem(file string, idx int, fixtureSteps []FixtureStep, lv1t *testing.T) {
	step := fixtureSteps[idx]
	lv1t.Run(strconv.Itoa(idx)+"_"+step.ID, func(t *testing.T) {
		if FixtureTestOpts.IsParallel {
			t.Parallel()
		}
		if step.RunAfter.BlockWait > 0 {
			err := inttest.WaitForBlockInterval(step.RunAfter.BlockWait)
			if err != nil {
				t.Fatal(err)
			}
		}
		RunActionRunner(step.Action, step, t)
		PropertyExistCheck(step, t)
		UpdateWorkQueueStatus(file, idx, fixtureSteps, Done, t)
	})
}

// RunSingleFixtureTest add a work queue into fixture test runner and execute work queues
func RunSingleFixtureTest(file string, t *testing.T) {
	t.Run(file, func(t *testing.T) {
		if FixtureTestOpts.IsParallel {
			t.Parallel()
		}
		var fixtureSteps []FixtureStep
		byteValue := ReadFile(file, t)

		err := json.Unmarshal([]byte(byteValue), &fixtureSteps)
		t.WithFields(testing.Fields{
			"raw_json": string(byteValue),
		}).MustNil(err, "something went wrong decoding fixture steps")

		CheckSteps(fixtureSteps, t)

		for idx, step := range fixtureSteps {
			workQueues = append(workQueues, QueueItem{
				fixtureFileName: file,
				idx:             idx,
				stepID:          step.ID,
				status:          NotStarted,
			})
		}

		for idx := range fixtureSteps {
			UpdateWorkQueueStatus(file, idx, fixtureSteps, InProgress, t)
		}
	})
}

// RunTestScenarios execute all scenarios
func RunTestScenarios(scenarioDir string, scenarioFileNames []string, t *originT.T) {
	newT := testing.NewT(t)

	var files []string

	scenarioDirectory := "scenarios"
	err := filepath.Walk(scenarioDirectory, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".json" {
			return nil
		}
		scenarioName := strings.TrimSuffix(info.Name(), ".json")
		t.Log(fmt.Sprintf("checking %s from %+v", scenarioName, scenarioFileNames))
		if len(scenarioFileNames) != 0 && !inttest.Exists(scenarioFileNames, scenarioName) {
			return nil
		}
		t.Log("added", scenarioName)
		files = append(files, path)
		return nil
	})
	if err != nil {
		t.Fatal("error walking through scenario directory", err)
	}
	for _, file := range files {
		t.Log("Running scenario path=", file)
		RunSingleFixtureTest(file, &newT)
	}
}
