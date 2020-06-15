package handlers

// FiatItemResponse is a struct to control fiat item response
type FiatItemResponse struct {
	ItemID  string `json:"ItemID"`
	Message string
	Status  string
}
