package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DoubleKeyValue describes double key/value set
type DoubleKeyValue struct {
	Key   string
	Value FloatString
}

// LongKeyValue describes long key/value set
type LongKeyValue struct {
	Key   string
	Value int
}

// StringKeyValue describes string key/value set
type StringKeyValue struct {
	Key   string
	Value string
}

// Item is a tradable asset
type Item struct {
	ID      string
	Doubles []DoubleKeyValue
	Longs   []LongKeyValue
	Strings []StringKeyValue
	// as items are unique per cookbook
	CookbookID    string
	Sender        sdk.AccAddress
	OwnerRecipeID string
	Tradable      bool
	LastUpdate    int64
}

// ItemList is a list of items
type ItemList []Item

// FindDouble is a function to get a double attribute from an item
func (it Item) FindDouble(key string) (float64, bool) {
	for _, v := range it.Doubles {
		if v.Key == key {
			return v.Value.Float(), true
		}
	}
	return 0, false
}

// FindDoubleKey is a function get double key index
func (it Item) FindDoubleKey(key string) (int, bool) {
	for i, v := range it.Doubles {
		if v.Key == key {
			return i, true
		}
	}
	return 0, false
}

// FindLong is a function to get a long attribute from an item
func (it Item) FindLong(key string) (int, bool) {
	for _, v := range it.Longs {
		if v.Key == key {
			return v.Value, true
		}
	}
	return 0, false
}

// FindLongKey is a function to get long key index
func (it Item) FindLongKey(key string) (int, bool) {
	for i, v := range it.Longs {
		if v.Key == key {
			return i, true
		}
	}
	return 0, false
}

// FindString is a function to get a string attribute from an item
func (it Item) FindString(key string) (string, bool) {
	for _, v := range it.Strings {
		if v.Key == key {
			return v.Value, true
		}
	}
	return "", false
}

// FindStringKey is a function to get string key index
func (it Item) FindStringKey(key string) (int, bool) {
	for i, v := range it.Strings {
		if v.Key == key {
			return i, true
		}
	}
	return 0, false
}

// SetString set item's string attribute
func (it Item) SetString(key string, value string) bool {
	for i, v := range it.Strings {
		if v.Key == key {
			it.Strings[i].Value = value
			return true
		}
	}
	return false
}

// NewItem create a new item with an auto generated ID
func NewItem(cookbookID string, doubles []DoubleKeyValue, longs []LongKeyValue, strings []StringKeyValue, sender sdk.AccAddress, BlockHeight int64) *Item {
	item := &Item{
		CookbookID: cookbookID,
		Doubles:    doubles,
		Longs:      longs,
		Strings:    strings,
		Sender:     sender,
		// By default all items are tradable
		Tradable:   true,
		LastUpdate: BlockHeight,
	}
	item.ID = KeyGen(sender)
	return item
}
