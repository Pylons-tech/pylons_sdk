package handlers

import (
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
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
			types.GenOneOutput(1),
			0
	case Rcp5BlockDelayed5xWoodcoinTo1xChaircoin: // 5 x woodcoin -> 1 x chair coin recipe, 5 block delayed
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenCoinOnlyEntry("chair"),
			types.GenOneOutput(1),
			5
	case Rcp5xWoodcoinTo1xRaichuItemBuy:
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenItemOnlyEntry("Raichu"),
			types.GenOneOutput(1),
			0
	case RcpRaichuNameUpgrade:
		return types.CoinInputList{},
			types.GenItemInputList("Raichu"),
			types.GenEntriesFirstItemNameUpgrade("RaichuV2"),
			types.GenOneOutput(1),
			0
	case RcpRaichuNameUpgradeWithCatalyst:
		return types.CoinInputList{},
			types.GenItemInputList("RaichuTC", "catalyst"),
			types.GenEntriesFirstItemNameUpgrade("RaichuTCV2"),
			types.GenOneOutput(1),
			0
	case Rcp2BlockDelayedKnifeUpgrade:
		return types.CoinInputList{},
			types.GenItemInputList("Knife"),
			types.GenEntriesFirstItemNameUpgrade("KnifeV2"),
			types.GenOneOutput(1),
			2
	case Rcp2BlockDelayedKnifeMerge:
		return types.CoinInputList{},
			types.GenItemInputList("Knife", "Knife"),
			types.GenItemOnlyEntry("KnifeMRG"),
			types.GenOneOutput(1),
			2
	case Rcp2BlockDelayedKnifeBuyer:
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenItemOnlyEntry("Knife"),
			types.GenOneOutput(1),
			2
	default: // 5 x woodcoin -> 1 x chair coin recipe, no delay
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenEntries("chair", "Raichu"),
			types.GenOneOutput(1),
			0
	}
}
