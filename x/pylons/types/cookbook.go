package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TypeCookbook is a store key for cookbook
const TypeCookbook = "cookbook"

// Cookbook is a struct that contains all the metadata of a cookbook
type Cookbook struct {
	NodeVersion  string
	ID           string // the cookbook guid
	Name         string
	Description  string
	Version      string
	Developer    string
	Level        int64
	SupportEmail string
	CostPerBlock int `json:",omitempty"`
	Sender       sdk.AccAddress
}

// NewCookbook return a new Cookbook
func NewCookbook(sEmail string, sender sdk.AccAddress, version, name, description, developer string, cpb int) Cookbook {
	cb := Cookbook{
		NodeVersion:  "0.0.1",
		Name:         name,
		Description:  description,
		Version:      version,
		Developer:    developer,
		SupportEmail: sEmail,
		Sender:       sender,
		CostPerBlock: cpb,
	}

	cb.ID = KeyGen(sender)
	return cb
}

func (cb Cookbook) String() string {
	return fmt.Sprintf(`
	Cookbook{ 
		NodeVersion: %s,
		Name: %s,
		Description: %s,
		Version: %s,
		Developer: %s,
		Level: %d,
		SupportEmail: %s,
		CostPerBlock: %d,
		Sender: %s,
	}`, cb.NodeVersion, cb.Name, cb.Description, cb.Version, cb.Developer, cb.Level, cb.SupportEmail, cb.CostPerBlock, cb.Sender)
}
