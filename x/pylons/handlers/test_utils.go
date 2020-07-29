package handlers

import (
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	mnemonicEntropySize = 256
)

// PopularRecipeType is a type for popular recipes
type PopularRecipeType int

// describes popular recipes
const (
	// a default recipe
	RcpDefault PopularRecipeType = 0
	// a recipe to convert 5x woodcoin to chaircoin
	Rcp5xWoodcoinTo1xChaircoin PopularRecipeType = 1
	// a recipe to convert 5x woodcoin to chaircoin, which is 5 block delayed
	Rcp5BlockDelayed5xWoodcoinTo1xChaircoin PopularRecipeType = 2
	// a recipe to convert 5x woodcoin to 1x raichu item
	Rcp5xWoodcoinTo1xRaichuItemBuy PopularRecipeType = 3
	// a recipe to upgrade raichu's name
	RcpRaichuNameUpgrade PopularRecipeType = 4
	// a recipe to upgrade raichu's name with catalyst item
	RcpRaichuNameUpgradeWithCatalyst PopularRecipeType = 5
	// a recipe to upgrade knife, which is 2 block delayed
	Rcp2BlockDelayedKnifeUpgrade PopularRecipeType = 6
	// a recipe to merge two knives, which is 2 block delayed
	Rcp2BlockDelayedKnifeMerge PopularRecipeType = 7
	// a recipe to buy a knife, which is 2 block delayed
	Rcp2BlockDelayedKnifeBuyer PopularRecipeType = 8
)

// GetParamsForPopularRecipe is a function to get popular recipe's attributes
func GetParamsForPopularRecipe(hfrt PopularRecipeType) (types.CoinInputList, types.ItemInputList, types.EntriesList, types.WeightedOutputsList, int64) {
	switch hfrt {
	case Rcp5xWoodcoinTo1xChaircoin: // 5 x woodcoin -> 1 x chair coin recipe
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenCoinOnlyEntry("chair"),
			types.GenOneOutput("chair"),
			0
	case Rcp5BlockDelayed5xWoodcoinTo1xChaircoin: // 5 x woodcoin -> 1 x chair coin recipe, 5 block delayed
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenCoinOnlyEntry("chair"),
			types.GenOneOutput("chair"),
			5
	case Rcp5xWoodcoinTo1xRaichuItemBuy:
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenItemOnlyEntry("Raichu"),
			types.GenOneOutput("Raichu"),
			0
	case RcpRaichuNameUpgrade:
		return types.CoinInputList{},
			types.GenItemInputList("Raichu"),
			types.GenEntriesFirstItemNameUpgrade("RaichuV2"),
			types.GenOneOutput("RaichuV2"),
			0
	case RcpRaichuNameUpgradeWithCatalyst:
		return types.CoinInputList{},
			types.GenItemInputList("RaichuTC", "catalyst"),
			types.GenEntriesFirstItemNameUpgrade("RaichuTCV2"),
			types.GenOneOutput("RaichuTCV2"),
			0
	case Rcp2BlockDelayedKnifeUpgrade:
		return types.CoinInputList{},
			types.GenItemInputList("Knife"),
			types.GenEntriesFirstItemNameUpgrade("KnifeV2"),
			types.GenOneOutput("KnifeV2"),
			2
	case Rcp2BlockDelayedKnifeMerge:
		return types.CoinInputList{},
			types.GenItemInputList("Knife", "Knife"),
			types.GenItemOnlyEntry("KnifeMRG"),
			types.GenOneOutput("KnifeMRG"),
			2
	case Rcp2BlockDelayedKnifeBuyer:
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenItemOnlyEntry("Knife"),
			types.GenOneOutput("Knife"),
			2
	default: // 5 x woodcoin -> 1 x chair coin recipe, no delay
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenEntries("chair", "Raichu"),
			types.GenOneOutput("chair", "Raichu"),
			0
	}
}

// GenAccount is a function to generate an account
func GenAccount() (secp256k1.PrivKeySecp256k1, sdk.AccAddress, error) {
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return secp256k1.PrivKeySecp256k1{}, nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return secp256k1.PrivKeySecp256k1{}, nil, err
	}

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return secp256k1.PrivKeySecp256k1{}, nil, err
	}

	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, "44'/118'/0'/0/0")
	if err != nil {
		return secp256k1.PrivKeySecp256k1{}, nil, err
	}

	priv := secp256k1.PrivKeySecp256k1(derivedPriv)
	cosmosAddr := sdk.AccAddress(priv.PubKey().Address().Bytes())
	return priv, cosmosAddr, nil
}
