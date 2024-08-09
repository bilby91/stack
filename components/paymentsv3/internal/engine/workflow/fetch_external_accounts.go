package workflow

import (
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/workflow"
)

type FetchNextExternalAccounts struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload json.RawMessage    `json:"fromPayload"`
}

func (s FetchNextExternalAccounts) GetWorkflow() any {
	return Workflow{}.runFetchNextExternalAccounts
}

func (w Workflow) runFetchNextExternalAccounts(ctx workflow.Context, fetchNextExternalAccount FetchNextExternalAccounts, nextTasks []models.TaskTree) (err error) {
	stateID := models.StateID{
		Reference:   workflow.GetInfo(ctx).WorkflowExecution.ID,
		ConnectorID: fetchNextExternalAccount.ConnectorID,
	}
	state, err := activities.FetchState(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return errors.Wrapf(err, "retrieving state: %s", stateID.String)
	}

	hasMore := true
	for hasMore {
		externalAccountsResponse, err := activities.FetchNextExternalAccounts(
			infiniteRetryContext(ctx),
			fetchNextExternalAccount.FromPayload,
			state.State,
			fetchNextExternalAccount.Config.PageSize(),
		)
		if err != nil {
			return errors.Wrap(err, "fetching next accounts")
		}

		err = activities.StoreAccounts(
			infiniteRetryContext(ctx),
			models.FromPSPAccounts(
				externalAccountsResponse.ExternalAccounts,
				models.ACCOUNT_TYPE_EXTERNAL,
				fetchNextExternalAccount.ConnectorID,
			),
		)
		if err != nil {
			return errors.Wrap(err, "storing next accounts")
		}

		state.State = externalAccountsResponse.NewState
		err = activities.StoreState(
			infiniteRetryContext(ctx),
			state,
		)
		if err != nil {
			return errors.Wrap(err, "storing state")
		}

		// TODO(polo): send event

		for _, externalAccount := range externalAccountsResponse.ExternalAccounts {
			payload, err := json.Marshal(externalAccount)
			if err != nil {
				return errors.Wrap(err, "marshalling external account")
			}

			if err := w.runNextWorkflow(
				ctx,
				fetchNextExternalAccount.Config,
				fetchNextExternalAccount.ConnectorID,
				payload,
				nextTasks,
			); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = externalAccountsResponse.HasMore
	}

	return nil
}
