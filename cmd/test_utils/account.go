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
	output, _, err := RunPylonsCli(params, "")
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(output, &result)
	return result, err
}

// CreateChainAccount is a function to create account on chain
func CreateChainAccount(key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key is empty")
	}
	params := []string{"tx", "pylons", "create-account", "add", "--from", key}
	output, _, err := RunPylonsCli(params, "")
	return string(output), err
}
