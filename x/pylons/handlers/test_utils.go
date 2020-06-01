package handlers

import (
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
)

type PopularRecipeType int

const (
	RCP_DEFAULT                                   PopularRecipeType = 0
	RCP_5xWOODCOIN_TO_1xCHAIRCOIN                 PopularRecipeType = 1
	RCP_5_BLOCK_DELAYED_5xWOODCOIN_TO_1xCHAIRCOIN PopularRecipeType = 2
	RCP_5xWOODCOIN_1xRAICHU_BUY                   PopularRecipeType = 3
	RCP_RAICHU_NAME_UPGRADE                       PopularRecipeType = 4
	RCP_RAICHU_NAME_UPGRADE_WITH_CATALYST         PopularRecipeType = 5
	RCP_2_BLOCK_DELAYED_KNIFE_UPGRADE             PopularRecipeType = 6
	RCP_2_BLOCK_DELAYED_KNIFE_MERGE               PopularRecipeType = 7
	RCP_2_BLOCK_DELAYED_KNIFE_BUYER               PopularRecipeType = 8
)

func GetParamsForPopularRecipe(hfrt PopularRecipeType) (types.CoinInputList, types.ItemInputList, types.EntriesList, types.WeightedOutputsList, int64) {
	switch hfrt {
	case RCP_5xWOODCOIN_TO_1xCHAIRCOIN: // 5 x woodcoin -> 1 x chair coin recipe
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenCoinOnlyEntry("chair"),
			types.GenOneOutput(1),
			0
	case RCP_5_BLOCK_DELAYED_5xWOODCOIN_TO_1xCHAIRCOIN: // 5 x woodcoin -> 1 x chair coin recipe, 5 block delayed
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenCoinOnlyEntry("chair"),
			types.GenOneOutput(1),
			5
	case RCP_5xWOODCOIN_1xRAICHU_BUY:
		return types.GenCoinInputList("wood", 5),
			types.ItemInputList{},
			types.GenItemOnlyEntry("Raichu"),
			types.GenOneOutput(1),
			0
	case RCP_RAICHU_NAME_UPGRADE:
		return types.CoinInputList{},
			types.GenItemInputList("Raichu"),
			types.GenEntriesFirstItemNameUpgrade("RaichuV2"),
			types.GenOneOutput(1),
			0
	case RCP_RAICHU_NAME_UPGRADE_WITH_CATALYST:
		return types.CoinInputList{},
			types.GenItemInputList("RaichuTC", "catalyst"),
			types.GenEntriesFirstItemNameUpgrade("RaichuTCV2"),
			types.GenOneOutput(1),
			0
	case RCP_2_BLOCK_DELAYED_KNIFE_UPGRADE:
		return types.CoinInputList{},
			types.GenItemInputList("Knife"),
			types.GenEntriesFirstItemNameUpgrade("KnifeV2"),
			types.GenOneOutput(1),
			2
	case RCP_2_BLOCK_DELAYED_KNIFE_MERGE:
		return types.CoinInputList{},
			types.GenItemInputList("Knife", "Knife"),
			types.GenItemOnlyEntry("KnifeMRG"),
			types.GenOneOutput(1),
			2
	case RCP_2_BLOCK_DELAYED_KNIFE_BUYER:
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