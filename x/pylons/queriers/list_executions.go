package queriers

import (
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
)

// query endpoints supported by the nameservice Querier
const (
	KeyListExecutions = "list_executions"
)

// ExecResp is the response for ListExecutions
type ExecResp struct {
	Executions []types.Execution
}
