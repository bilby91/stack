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

type FetchNextOthers struct {
	Name        string          `json:"name"`
	FromPayload json.RawMessage `json:"fromPayload"`
	PageSize    int             `json:"pageSize"`
}

func (s FetchNextOthers) GetWorkflow() any {
	return Workflow{}.runFetchNextOthers
}

func (w Workflow) runFetchNextOthers(ctx workflow.Context, fetchNextOthers FetchNextOthers, nextTasks []*models.TaskTree) (err error) {
	var state json.RawMessage
	// TODO(polo): fetch state from database
	_ = state

	hasMore := true
	for hasMore {
		othersResponse, err := activities.FetchNextOthers(
			workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 60 * time.Second,
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:        time.Second,
					BackoffCoefficient:     2,
					MaximumInterval:        100 * time.Second,
					NonRetryableErrorTypes: []string{},
				},
			}),
			fetchNextOthers.Name,
			fetchNextOthers.FromPayload,
			state,
			fetchNextOthers.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next others")
		}

		// TODO(polo): store others and new state

		hasMore = othersResponse.HasMore
		state = othersResponse.NewState

		for _, other := range othersResponse.Others {
			payload, err := json.Marshal(other)
			if err != nil {
				return errors.Wrap(err, "marshalling other")
			}

			if err := w.runNextWorkflow(ctx, payload, fetchNextOthers.PageSize, nextTasks); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}
	}

	return nil
}
