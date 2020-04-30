package types

// CoinOutput is the game elements that are needs as output to a recipe
type CoinOutput struct {
	Coin string
	// coin output count is parsed by cel program
	Count string
}
