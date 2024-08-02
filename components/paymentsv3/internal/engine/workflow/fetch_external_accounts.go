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

type FetchNextExternalAccounts struct {
	FromPayload json.RawMessage `json:"fromPayload"`
	PageSize    int             `json:"pageSize"`
}

func (s FetchNextExternalAccounts) GetWorkflow() any {
	return Workflow{}.runFetchNextExternalAccounts
}

func (w Workflow) runFetchNextExternalAccounts(ctx workflow.Context, fetchNextExternalAccount FetchNextExternalAccounts, nextTasks []*models.TaskTree) (err error) {
	var state json.RawMessage
	// TODO(polo): fetch state from database
	_ = state

	hasMore := true
	for hasMore {
		externalAccountsResponse, err := activities.FetchNextExternalAccounts(
			workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 60 * time.Second,
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:        time.Second,
					BackoffCoefficient:     2,
					MaximumInterval:        100 * time.Second,
					NonRetryableErrorTypes: []string{},
				},
			}),
			fetchNextExternalAccount.FromPayload,
			state,
			fetchNextExternalAccount.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next accounts")
		}

		// TODO(polo): store accounts and new state

		hasMore = externalAccountsResponse.HasMore
		state = externalAccountsResponse.NewState

		for _, externalAccount := range externalAccountsResponse.ExternalAccounts {
			payload, err := json.Marshal(externalAccount)
			if err != nil {
				return errors.Wrap(err, "marshalling external account")
			}

			if err := w.runNextWorkflow(ctx, payload, fetchNextExternalAccount.PageSize, nextTasks); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}
	}

	return nil
}
