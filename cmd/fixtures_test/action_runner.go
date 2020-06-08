package fixturetest

import (
	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"
)

// ActFunc describes the type of function used for action running test
type ActFunc func(FixtureStep, *testing.T)

var actFuncs = make(map[string]ActFunc)

// RegisterActionRunner registers action runner function
func RegisterActionRunner(action string, fn ActFunc) {
	actFuncs[action] = fn
}

// GetActionRunner get registered action runner function
func GetActionRunner(action string) ActFunc {
	return actFuncs[action]
}

// RunActionRunner execute registered action runner function
func RunActionRunner(action string, step FixtureStep, t *testing.T) {
	fn := GetActionRunner(action)
	t.MustTrue(fn != nil)
	fn(step, t)
}
