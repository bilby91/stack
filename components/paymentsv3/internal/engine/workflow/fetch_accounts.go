package workflow

import (
	"encoding/json"

	"go.temporal.io/sdk/workflow"
)

type FetchNextAccounts struct {
	FromPayload json.RawMessage `json:"fromPayload"`
	PageSize    int             `json:"pageSize"`
}

func (s FetchNextAccounts) GetWorkflow() any {
	return RunFetchNextAccounts
}

func RunFetchNextAccounts(ctx workflow.Context, fetchNextAccount FetchNextAccounts) (err error) {
	return nil
}
