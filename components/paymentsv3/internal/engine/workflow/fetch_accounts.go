package workflow

import (
	"encoding/json"
	"log"
	"time"

	"github.com/formancehq/paymentsv3/internal/engine/activities"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
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
	var state json.RawMessage
	// TODO(polo): fetch state from database
	_ = state

	hasMore := true
	for hasMore {
		accountsResponse, err := activities.FetchNextAccounts(
			workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 60 * time.Second,
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:        time.Second,
					BackoffCoefficient:     2,
					MaximumInterval:        100 * time.Second,
					NonRetryableErrorTypes: []string{},
				},
			}),
			fetchNextAccount.FromPayload,
			state,
			fetchNextAccount.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next accounts")
		}

		log.Println(accountsResponse)
		hasMore = accountsResponse.HasMore
	}

	// TODO(polo): store accounts and new state

	return nil
}
