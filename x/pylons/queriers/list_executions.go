package queriers

import (
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
)

// query endpoints supported by the nameservice Querier
const (
	KeyListExecutions = "list_executions"
)

// ExecResponse is the response for ListExecutions
type ExecResponse struct {
	Executions []types.Execution
}
