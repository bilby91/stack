package workflow

import (
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/connectors/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type FetchNextExternalAccounts struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload json.RawMessage    `json:"fromPayload"`
}

func (w Workflow) runFetchNextExternalAccounts(
	ctx workflow.Context,
	plugin models.Plugin,
	fetchNextExternalAccount FetchNextExternalAccounts,
	nextTasks []models.TaskTree,
) error {
	stateID := models.StateID{
		Reference:   workflow.GetInfo(ctx).WorkflowExecution.ID,
		ConnectorID: fetchNextExternalAccount.ConnectorID,
	}
	state, err := activities.StorageStatesGet(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return errors.Wrapf(err, "retrieving state: %s", stateID.String)
	}

	hasMore := true
	for hasMore {
		externalAccountsResponse, err := activities.PluginFetchNextExternalAccounts(
			infiniteRetryContext(ctx),
			plugin,
			fetchNextExternalAccount.FromPayload,
			state.State,
			fetchNextExternalAccount.Config.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next accounts")
		}

		err = activities.StorageAccountsStore(
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
		err = activities.StorageStatesStore(
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

			if err := workflow.ExecuteChildWorkflow(
				workflow.WithChildOptions(
					ctx,
					workflow.ChildWorkflowOptions{
						TaskQueue:         fetchNextExternalAccount.ConnectorID.Reference,
						ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
					},
				),
				Run,
				plugin,
				fetchNextExternalAccount.Config,
				fetchNextExternalAccount.ConnectorID,
				payload,
				nextTasks,
			).Get(ctx, nil); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = externalAccountsResponse.HasMore
	}

	return nil
}

var RunFetchNextExternalAccounts any

func init() {
	RunFetchNextExternalAccounts = Workflow{}.runFetchNextExternalAccounts
}
