package handlers

// ExecuteRecipeResp is the response for executeRecipe
type ExecuteRecipeResp struct {
	Message string
	Status  string
	Output  []byte
}

// ExecuteRecipeSerialize is a struct for execute recipe result serialization
type ExecuteRecipeSerialize struct {
	Type   string `json:"type"`   // COIN or ITEM
	Coin   string `json:"coin"`   // used when type is ITEM
	Amount int64  `json:"amount"` // used when type is COIN
	ItemID string `json:"itemID"` // used when type is ITEM
}

// ExecuteRecipeScheduleOutput is a struct that shows how execute recipe schedule output works
type ExecuteRecipeScheduleOutput struct {
	ExecID string
}
