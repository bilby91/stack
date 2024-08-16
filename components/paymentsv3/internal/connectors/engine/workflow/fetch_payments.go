package workflow

import (
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/connectors/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
)

type FetchNextPayments struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload json.RawMessage    `json:"fromPayload"`
}

func (w Workflow) runFetchNextPayments(
	ctx workflow.Context,
	plugin models.Plugin,
	fetchNextPayments FetchNextPayments,
	nextTasks []models.TaskTree,
) error {
	stateID := models.StateID{
		Reference:   workflow.GetInfo(ctx).WorkflowExecution.ID,
		ConnectorID: fetchNextPayments.ConnectorID,
	}
	state, err := activities.StorageStatesGet(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return errors.Wrapf(err, "retrieving state: %s", stateID.String)
	}

	hasMore := true
	for hasMore {
		paymentsResponse, err := activities.PluginFetchNextPayments(
			infiniteRetryContext(ctx),
			plugin,
			fetchNextPayments.FromPayload,
			state.State,
			fetchNextPayments.Config.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next payments")
		}

		err = activities.StoragePaymentsStore(
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
		err = activities.StorageStatesStore(
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

			if err := w.run(
				ctx,
				plugin,
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

var RunFetchNextPayments any

func init() {
	RunFetchNextPayments = Workflow{}.runFetchNextPayments
}
