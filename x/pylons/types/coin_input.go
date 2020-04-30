package types

// CoinInput is the game elements that are needs as inputs to a recipe
type CoinInput struct {
	Coin  string
	Count int64 `json:",string"` // TODO: This is added since we are using json.Marshal and amino Marshal together, when we convert all the Marshal into amino json:",string" can be removed
}

// CoinInputList is a list of Coin inputs
type CoinInputList []CoinInput
