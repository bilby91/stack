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
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload json.RawMessage    `json:"fromPayload"`
}

func (s FetchNextPayments) GetWorkflow() any {
	return Workflow{}.runFetchNextPayments
}

func (w Workflow) runFetchNextPayments(ctx workflow.Context, fetchNextPayments FetchNextPayments, nextTasks []models.TaskTree) (err error) {
	stateID := models.StateID{
		Reference:   workflow.GetInfo(ctx).WorkflowExecution.ID,
		ConnectorID: fetchNextPayments.ConnectorID,
	}
	state, err := activities.FetchState(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return errors.Wrapf(err, "retrieving state: %s", stateID.String)
	}

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
			state.State,
			fetchNextPayments.Config.PageSize(),
		)
		if err != nil {
			return errors.Wrap(err, "fetching next payments")
		}

		err = activities.StorePayments(
			infiniteRetryContext(ctx),
			models.FromPSPPayments(
				paymentsResponse.Payments,
				fetchNextPayments.ConnectorID,
			),
		)
		if err != nil {
			return errors.Wrap(err, "storing next accounts")
		}

		state.State = paymentsResponse.NewState
		err = activities.StoreState(
			infiniteRetryContext(ctx),
			state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		// TODO(polo): send events

		for _, payment := range paymentsResponse.Payments {
			payload, err := json.Marshal(payment)
			if err != nil {
				return errors.Wrap(err, "marshalling payment")
			}

			if err := w.runNextWorkflow(
				ctx,
				fetchNextPayments.Config,
				fetchNextPayments.ConnectorID,
				payload,
				nextTasks,
			); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = paymentsResponse.HasMore
	}

	return nil
}
