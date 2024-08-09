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
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	Name        string             `json:"name"`
	FromPayload json.RawMessage    `json:"fromPayload"`
}

func (s FetchNextOthers) GetWorkflow() any {
	return Workflow{}.runFetchNextOthers
}

func (w Workflow) runFetchNextOthers(ctx workflow.Context, fetchNextOthers FetchNextOthers, nextTasks []models.TaskTree) (err error) {
	stateID := models.StateID{
		Reference:   workflow.GetInfo(ctx).WorkflowExecution.ID,
		ConnectorID: fetchNextOthers.ConnectorID,
	}
	state, err := activities.FetchState(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return errors.Wrapf(err, "retrieving state: %s", stateID.String)
	}

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
			state.State,
			fetchNextOthers.Config.PageSize(),
		)
		if err != nil {
			return errors.Wrap(err, "fetching next others")
		}

		state.State = othersResponse.NewState
		err = activities.StoreState(
			infiniteRetryContext(ctx),
			state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		// TODO(polo): send event for others ? store others ?

		for _, other := range othersResponse.Others {
			payload, err := json.Marshal(other)
			if err != nil {
				return errors.Wrap(err, "marshalling other")
			}

			if err := w.runNextWorkflow(
				ctx,
				fetchNextOthers.Config,
				fetchNextOthers.ConnectorID,
				payload,
				nextTasks,
			); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = othersResponse.HasMore
	}

	return nil
}
