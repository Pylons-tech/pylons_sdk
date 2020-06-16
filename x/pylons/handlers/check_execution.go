package handlers

// CheckExecutionResponse is the response for checkExecution
type CheckExecutionResponse struct {
	Message string
	Status  string
	Output  []byte
}
