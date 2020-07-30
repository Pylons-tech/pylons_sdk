package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Cookbook is a struct that contains all the metadata of a cookbook
type Cookbook struct {
	NodeVersion  SemVer
	ID           string // the cookbook guid
	Name         string
	Description  string
	Version      SemVer
	Developer    string
	Level        Level
	SupportEmail Email
	CostPerBlock int `json:",omitempty"`
	Sender       sdk.AccAddress
}

// CookbookList is a list of cookbook
type CookbookList struct {
	Cookbooks []Cookbook
}

// NewCookbook return a new Cookbook
func NewCookbook(sEmail Email, sender sdk.AccAddress, version SemVer, name, description, developer string, cpb int) Cookbook {
	cb := Cookbook{
		NodeVersion:  SemVer("0.0.1"),
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
