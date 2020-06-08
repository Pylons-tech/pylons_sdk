package handlers

// CheckExecutionResp is the response for checkExecution
type CheckExecutionResp struct {
	Message string
	Status  string
	Output  []byte
}
