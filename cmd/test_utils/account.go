package inttest

import (
	"encoding/json"
	"errors"
)

// AddNewLocalKey is a function to add key cli
func AddNewLocalKey(key string) (map[string]string, error) {
	result := make(map[string]string)
	if len(key) == 0 {
		return result, errors.New("key is empty")
	}
	params := []string{"keys", "add", key}
	output, logstr, err := RunPylonsCli(params, "")
	if err != nil {
		result["logstr"] = logstr
		result["output"] = string(output)
		return result, err
	}
	err = json.Unmarshal(output, &result)
	return result, err
}

// CreateChainAccount is a function to create account on chain
func CreateChainAccount(key string) (string, string, error) {
	if len(key) == 0 {
		return "", "", errors.New("key is empty")
	}
	params := []string{"tx", "pylons", "create-account", "--from", key}
	output, logstr, err := RunPylonsCli(params, "y\n")
	return string(output), logstr, err
}
