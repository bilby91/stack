package workflow

import (
	"encoding/json"
	"time"

	"github.com/formancehq/paymentsv3/internal/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type FetchNextAccounts struct {
	FromPayload json.RawMessage `json:"fromPayload"`
	PageSize    int             `json:"pageSize"`
}

func (s FetchNextAccounts) GetWorkflow() any {
	return Workflow{}.runFetchNextAccounts
}

func (w Workflow) runFetchNextAccounts(ctx workflow.Context, fetchNextAccount FetchNextAccounts, nextTasks []*models.TaskTree) (err error) {
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

		// TODO(polo): store accounts and new state

		hasMore = accountsResponse.HasMore
		state = accountsResponse.NewState

		for _, account := range accountsResponse.Accounts {
			payload, err := json.Marshal(account)
			if err != nil {
				return errors.Wrap(err, "marshalling account")
			}

			if err := w.runNextWorkflow(ctx, payload, fetchNextAccount.PageSize, nextTasks); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}
	}

	return nil
}
