package handlers

// UpdateRecipeResponse is a struct to control update recipe response
type UpdateRecipeResponse struct {
	RecipeID string `json:"RecipeID"`
	Message  string
	Status   string
}
