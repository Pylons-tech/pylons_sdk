package fixturetest

import (
	"flag"
	"strings"
	"testing"

	inttest "github.com/Pylons-tech/pylons_sdk/cmd/test"
)

var runSerialMode = false
var useRest = false
var useKnownCookbook = false
var scenarios = ""
var accounts = ""

func init() {
	flag.BoolVar(&runSerialMode, "runserial", false, "true/false value to check if test will be running in parallel")
	flag.BoolVar(&useRest, "userest", false, "use rest endpoint for Tx send")
	flag.BoolVar(&useKnownCookbook, "use-known-cookbook", false, "use existing cookbook or not")
	flag.StringVar(&scenarios, "scenarios", "", "custom scenario file names")
	flag.StringVar(&accounts, "accounts", "", "custom account names")
}

func TestFixturesViaCLI(t *testing.T) {
	flag.Parse()
	FixtureTestOpts.IsParallel = !runSerialMode
	FixtureTestOpts.CreateNewCookbook = !useKnownCookbook
	if useRest {
		inttest.CLIOpts.RestEndpoint = "http://localhost:1317"
	}
	inttest.CLIOpts.MaxBroadcast = 50
	RegisterDefaultActionRunners()
	// Register custom action runners
	// RegisterActionRunner("custom_action", CustomActionRunner)
	scenarioFileNames := []string{}
	accountNames := []string{}
	if len(scenarios) > 0 {
		scenarioFileNames = strings.Split(scenarios, ",")
		accountNames = strings.Split(accounts, ",")
	}
	RunTestScenarios("scenarios", scenarioFileNames, accountNames, t)
}
