package handlers

// CheckExecutionResp is a struct of check execution response
type CheckExecutionResp struct {
	Message string
	Status  string
	Output  []byte
}
