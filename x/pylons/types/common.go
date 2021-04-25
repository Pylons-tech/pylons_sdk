package types

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

type StringKeyValueList []StringKeyValue
type DoubleKeyValueList []DoubleKeyValue
type LongKeyValueList []LongKeyValue

type LongParamList []LongParam
type DoubleParamList []DoubleParam
type StringParamList []StringParam

type DoubleInputParamList []DoubleInputParam
type LongInputParamList []LongInputParam
type StringInputParamList []StringInputParam
type WeightedOutputsList []WeightedOutputs

type CoinInputList []CoinInput
type ItemList []Item
type ItemInputList []ItemInput
type TradeItemInputList []TradeItemInput

type DoubleWeightTable []DoubleWeightRange
type IntWeightTable []IntWeightRange

// Tier defines the kind of cookbook this is
type Tier struct {
	Level int64
	Fee   sdk.Coins
}

// ItemHistory is a struct to store Item use history
type ItemHistory struct {
	ID       string
	Owner    sdk.AccAddress
	ItemID   string
	RecipeID string
	TradeID  string
}

// Reader struct is for entropy set on uuid
type Reader struct{}

// Tier defines the kind of cookbook this is
const (
	// Basic is the free int64 which does allow developers to use pylons ( paid currency ) in their
	// games
	Basic int64 = iota
	Premium
)

const (
	Pylon = "pylon"

	TypeCookbook    = "cookbook"
	TypeRecipe      = "recipe"
	TypeTrade       = "trade"
	TypeItem        = "item"
	TypeItemHistory = "item_history"
	TypeExecution   = "execution"
)

// tier fee types
var (
	BasicFee   = NewPylon(10000) // fee charged to create a basic cookbook
	PremiumFee = NewPylon(50000) // fee charged to create a premium cookbook
)

// BasicTier is the cookbook tier which doesn't allow paid recipes which means
// the developers cannot have recipes where they can actually carge a fee in pylons
var BasicTier = Tier{
	Level: Basic,
	Fee:   BasicFee,
}

// PremiumTier the cookbook tier which does allow paid recipes
var PremiumTier = Tier{
	Level: Premium,
	Fee:   PremiumFee,
}

// NewPylon Returns pylon currency
func NewPylon(amount int64) sdk.Coins {
	return sdk.Coins{sdk.NewInt64Coin(Pylon, amount)}
}

// NewCookbook return a new Cookbook
func NewCookbook(sEmail string, sender sdk.AccAddress, version, name, description, developer string, cpb int64) Cookbook {
	cb := Cookbook{
		NodeVersion:  "0.0.1",
		Name:         name,
		Description:  description,
		Version:      version,
		Developer:    developer,
		SupportEmail: sEmail,
		Sender:       sender.String(),
		CostPerBlock: cpb,
	}

	cb.ID = KeyGen(sender)
	return cb
}

// NewRecipe creates a new recipe
func NewRecipe(recipeName, cookbookID, description string,
	coinInputs CoinInputList, // coins to put on the recipe
	itemInputs ItemInputList, // items to put on the recipe
	entries EntriesList, // items that can be created from recipe
	outputs WeightedOutputsList, // item outputs listing by weight value
	blockInterval int64, // The amount of time to wait to finish running the recipe
	sender sdk.AccAddress,
	extraInfo string) Recipe {
	rcp := Recipe{
		NodeVersion:   "0.0.1",
		Name:          recipeName,
		CookbookID:    cookbookID,
		CoinInputs:    coinInputs,
		ItemInputs:    itemInputs,
		Entries:       entries,
		Outputs:       outputs,
		BlockInterval: blockInterval,
		Description:   description,
		Sender:        sender.String(),
		ExtraInfo:     extraInfo,
	}

	rcp.ID = KeyGen(sender)
	return rcp
}

// NewItem create a new item with an auto generated ID
func NewItem(cookbookID string, doubles DoubleKeyValueList, longs LongKeyValueList, strings StringKeyValueList, sender sdk.AccAddress, blockHeight int64, transferFee int64) Item {
	item := Item{
		NodeVersion: "0.0.1",
		CookbookID:  cookbookID,
		Doubles:     doubles,
		Longs:       longs,
		Strings:     strings,
		Sender:      sender.String(),
		// By default all items are tradable
		Tradable:    true,
		LastUpdate:  blockHeight,
		TransferFee: transferFee,
	}
	item.ID = KeyGen(sender)

	return item
}

// Equals compares two items
func (it Item) Equals(other Item) bool {
	return it.ID == other.ID &&
		reflect.DeepEqual(it.Doubles, other.Doubles) &&
		reflect.DeepEqual(it.Strings, other.Strings) &&
		reflect.DeepEqual(it.Longs, other.Longs) &&
		reflect.DeepEqual(it.CookbookID, other.CookbookID)
}

// MatchItemInput checks if the ItemInput matches the item
func (it Item) MatchItemInput(other ItemInput) bool {
	return reflect.DeepEqual(it.Doubles, other.Doubles) &&
		reflect.DeepEqual(it.Strings, other.Strings) &&
		reflect.DeepEqual(it.Longs, other.Longs)
}

// NewTradeError check if an item can be sent to someone else
func (it Item) NewTradeError() error {
	if !it.Tradable {
		return errors.New("Item Tradable flag is not set")
	}
	if it.OwnerRecipeID != "" {
		return errors.New("Item is owned by a recipe")
	}
	if it.OwnerTradeID != "" {
		return errors.New("Item is owned by a trade")
	}
	return nil
}

// FulfillTradeError check if an item can be sent to someone else
func (it Item) FulfillTradeError(tradeID string) error {
	if !it.Tradable {
		return errors.New("Item Tradable flag is not set")
	}
	if it.OwnerRecipeID != "" {
		return errors.New("Item is owned by a recipe")
	}
	if it.OwnerTradeID != tradeID {
		return errors.New("Item is not owned by the trade")
	}
	return nil
}

// NewExecution return a new Execution
func NewExecution(recipeID string, cookbookID string, ci sdk.Coins,
	itemInputs []Item,
	blockHeight int64, sender sdk.AccAddress,
	completed bool) Execution {

	exec := Execution{
		NodeVersion: "0.0.1",
		RecipeID:    recipeID,
		CookbookID:  cookbookID,
		CoinInputs:  ci,
		ItemInputs:  itemInputs,
		BlockHeight: blockHeight,
		Sender:      sender.String(),
		Completed:   completed,
	}

	exec.ID = KeyGen(sender)
	return exec
}

// NewRecipeExecutionError is a utility that shows if Recipe is compatible with recipe execution
func (it Item) NewRecipeExecutionError() error {
	if it.OwnerRecipeID != "" {
		return errors.New("Item is owned by a recipe")
	}
	if it.OwnerTradeID != "" {
		return errors.New("Item is owned by a trade")
	}
	return nil
}

// SetTransferFee set item's TransferFee
func (it *Item) SetTransferFee(transferFee int64) {
	it.TransferFee = transferFee
}

// FindDouble is a function to get a double attribute from an item
func (it Item) FindDouble(key string) (sdk.Dec, bool) {
	for _, v := range it.Doubles {
		if v.Key == key {
			return v.Value, true
		}
	}
	return sdk.NewDec(0), false
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
			return int(v.Value), true
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

// IDValidationError check if ID can be used as a variable name if available
func (ii ItemInput) IDValidationError() error {
	if len(ii.ID) == 0 {
		return nil
	}

	regex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z_0-9]*$`)
	if regex.MatchString(ii.ID) {
		return nil
	}

	return fmt.Errorf("ID is not empty nor fit the regular expression ^[a-zA-Z_][a-zA-Z_0-9]*$: id=%s", ii.ID)
}

// NewItemModifyOutput returns ItemOutput that is modified from item input
func NewItemModifyOutput(ID string, ItemInputRef string, ModifyParams ItemModifyParams) ItemModifyOutput {
	return ItemModifyOutput{
		ID:           ID,
		ItemInputRef: ItemInputRef,
		Doubles:      ModifyParams.Doubles,
		Longs:        ModifyParams.Longs,
		Strings:      ModifyParams.Strings,
		TransferFee:  ModifyParams.TransferFee,
	}
}

// SetTransferFee set generate item's transfer fee
func (mit *ItemModifyOutput) SetTransferFee(transferFee int64) {
	mit.TransferFee = transferFee
}

// NewItemOutput returns new ItemOutput generated from recipe
func NewItemOutput(ID string, Doubles DoubleParamList, Longs LongParamList, Strings StringParamList, TransferFee int64) ItemOutput {
	return ItemOutput{
		ID:          ID,
		Doubles:     Doubles,
		Longs:       Longs,
		Strings:     Strings,
		TransferFee: TransferFee,
	}
}

// SetTransferFee set generate item's transfer fee
func (io *ItemOutput) SetTransferFee(transferFee int64) {
	io.TransferFee = transferFee
}

// NewTrade creates a new trade
func NewTrade(extraInfo string,
	coinInputs CoinInputList, // coinOutputs CoinOutputList,
	itemInputs TradeItemInputList, // itemOutputs ItemOutputList,
	coinOutputs sdk.Coins, // newly created param instead of coinOutputs and itemOutputs
	itemOutputs ItemList,
	sender sdk.AccAddress) Trade {
	trd := Trade{
		NodeVersion: "0.0.1",
		CoinInputs:  coinInputs,
		ItemInputs:  itemInputs,
		CoinOutputs: coinOutputs,
		ItemOutputs: itemOutputs,
		ExtraInfo:   extraInfo,
		Sender:      sender.String(),
	}

	trd.ID = KeyGen(sender)
	return trd
}

// GetItemInputRefIndex get item input index from ref string
func (rcp Recipe) GetItemInputRefIndex(inputRef string) int {
	for idx, input := range rcp.ItemInputs {
		if input.ID == inputRef {
			return idx
		}
	}
	return -1
}

// ToCoins converts to coins
func (cil CoinInputList) ToCoins() sdk.Coins {
	var coins sdk.Coins
	for _, ci := range cil {
		coins = append(coins, sdk.NewInt64Coin(ci.Coin, ci.Count))
	}
	coins = coins.Sort()
	return coins
}

// Equal compares two inputlists
func (cil CoinInputList) Equal(other CoinInputList) bool {
	for _, inp := range cil {
		found := false
		for _, oinp := range other {
			if oinp.Coin == inp.Coin && oinp.Count == inp.Count {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// Validate is a function to check ItemInputList is valid
func (iil ItemInputList) Validate() error {
	return nil
}

// Validate is a function to check ItemInputList is valid
func (tiil TradeItemInputList) Validate() error {
	for _, ii := range tiil {
		if ii.CookbookID == "" {
			return errors.New("There should be no empty cookbook ID inputs for trades")
		}
	}
	return nil
}

// NewEntropyReader create an entropy reader
func NewEntropyReader() *Reader {
	return &Reader{}
}

func (r Reader) Read(b []byte) (n int, err error) {
	entropy := []byte{}
	for i := 0; i < len(b); i++ {
		entropy = append(entropy, byte(rand.Intn(256)))
	}

	n = copy(b, entropy)
	return n, nil
}

// ValidateLevel validates the level
func ValidateLevel(level int64) error {
	if level == Basic || level == Premium {
		return nil
	}

	return errors.New("Invalid cookbook plan")
}

// KeyGen generates key for the store
func KeyGen(sender sdk.AccAddress) string {
	id := uuid.New()
	return sender.String() + id.String()
}

// ValidateEmail validates the email string provided
func ValidateEmail(email string) error {
	exp := regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z0-9]{2,})$`)
	if exp.MatchString(email) {
		return nil
	}

	return errors.New("invalid email address")
}

// ValidateVersion validates the SemVer
func ValidateVersion(s string) error {
	regex := regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?$`)
	if regex.MatchString(s) {
		return nil
	}

	return errors.New("invalid semVer")
}

// Reference: BlockHeader struct is available at github.com/tendermint/tendermint@v0.33.4/types/proto3/block.pb.go:
// type Header struct {
// 		AppHash            []byte `protobuf:"bytes,11,opt,name=app_hash,json=appHash,proto3" json:"app_hash,omitempty"`
//   	...
// }

// RandomSeed calculate random seed from context and entity count
func RandomSeed(ctx sdk.Context, entityCount int) int64 {
	header := ctx.BlockHeader()
	appHash := header.AppHash
	seedValue := 0
	for i, bytv := range appHash { // len(appHash) = 11
		intv := int(bytv)
		seedValue += (i*i + 1) * intv
	}
	fmt.Println("RandomSeed entityCount:", entityCount, "BlockHeight:", header.Height)
	return int64(seedValue + entityCount)
}

// Max returns the larger of x or y.
func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

// Min returns the larger of x or y.
func Min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}
