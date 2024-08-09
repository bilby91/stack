package workflow

import (
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
)

type FetchNextAccounts struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload json.RawMessage    `json:"fromPayload"`
}

func (s FetchNextAccounts) GetWorkflow() any {
	return Workflow{}.runFetchNextAccounts
}

func (w Workflow) runFetchNextAccounts(ctx workflow.Context, fetchNextAccount FetchNextAccounts, nextTasks []models.TaskTree) (err error) {
	stateID := models.StateID{
		Reference:   workflow.GetInfo(ctx).WorkflowExecution.ID,
		ConnectorID: fetchNextAccount.ConnectorID,
	}
	state, err := activities.FetchState(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return errors.Wrapf(err, "retrieving state: %s", stateID.String)
	}

	hasMore := true
	for hasMore {
		accountsResponse, err := activities.FetchNextAccounts(
			infiniteRetryContext(ctx),
			fetchNextAccount.FromPayload,
			state.State,
			fetchNextAccount.Config.PageSize(),
		)
		if err != nil {
			return errors.Wrap(err, "fetching next accounts")
		}

		err = activities.StoreAccounts(
			infiniteRetryContext(ctx),
			models.FromPSPAccounts(
				accountsResponse.Accounts,
				models.ACCOUNT_TYPE_INTERNAL,
				fetchNextAccount.ConnectorID,
			),
		)
		if err != nil {
			return errors.Wrap(err, "storing next accounts")
		}

		state.State = accountsResponse.NewState
		err = activities.StoreState(
			infiniteRetryContext(ctx),
			state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		// TODO(polo): send event

		for _, account := range accountsResponse.Accounts {
			payload, err := json.Marshal(account)
			if err != nil {
				return errors.Wrap(err, "marshalling account")
			}

			if err := w.runNextWorkflow(
				ctx,
				fetchNextAccount.Config,
				fetchNextAccount.ConnectorID,
				payload,
				nextTasks,
			); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = accountsResponse.HasMore
	}

	return nil
}
