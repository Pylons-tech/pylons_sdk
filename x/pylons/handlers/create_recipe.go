package handlers

// CreateRecipeResponse is struct of create recipe response
type CreateRecipeResponse struct {
	RecipeID string `json:"RecipeID"`
	Message  string
	Status   string
}
