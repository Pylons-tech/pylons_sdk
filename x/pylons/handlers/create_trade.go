package handlers

// CreateTradeResponse is struct of create trade response
type CreateTradeResponse struct {
	TradeID string `json:"TradeID"`
	Message string
	Status  string
}
