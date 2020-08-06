package handlers

// UpdateCookbookResponse is a struct of update cookbook response
type UpdateCookbookResponse struct {
	CookbookID string `json:"CookbookID"`
	Message    string
	Status     string
}
