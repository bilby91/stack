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

type FetchNextPayments struct {
	FromPayload json.RawMessage `json:"fromPayload"`
	PageSize    int             `json:"pageSize"`
}

func (s FetchNextPayments) GetWorkflow() any {
	return Workflow{}.runFetchNextPayments
}

func (w Workflow) runFetchNextPayments(ctx workflow.Context, fetchNextPayments FetchNextPayments, nextTasks []*models.TaskTree) (err error) {
	var state json.RawMessage
	// TODO(polo): fetch state from database
	_ = state

	hasMore := true
	for hasMore {
		paymentsResponse, err := activities.FetchNextPayments(
			workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 60 * time.Second,
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:        time.Second,
					BackoffCoefficient:     2,
					MaximumInterval:        100 * time.Second,
					NonRetryableErrorTypes: []string{},
				},
			}),
			fetchNextPayments.FromPayload,
			state,
			fetchNextPayments.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next payments")
		}

		// TODO(polo): store payments and new state

		hasMore = paymentsResponse.HasMore
		state = paymentsResponse.NewState

		for _, payment := range paymentsResponse.Payments {
			payload, err := json.Marshal(payment)
			if err != nil {
				return errors.Wrap(err, "marshalling payment")
			}

			if err := w.runNextWorkflow(ctx, payload, fetchNextPayments.PageSize, nextTasks); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}
	}

	return nil
}
