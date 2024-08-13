package workflow

import (
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/connectors/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type FetchNextAccounts struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	FromPayload json.RawMessage    `json:"fromPayload"`
}

func (w Workflow) runFetchNextAccounts(
	ctx workflow.Context,
	plugin models.Plugin,
	fetchNextAccount FetchNextAccounts,
	nextTasks []models.TaskTree,
) error {
	stateID := models.StateID{
		Reference:   workflow.GetInfo(ctx).WorkflowExecution.ID,
		ConnectorID: fetchNextAccount.ConnectorID,
	}
	state, err := activities.StorageFetchState(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return errors.Wrapf(err, "retrieving state: %s", stateID.String)
	}

	hasMore := true
	for hasMore {
		accountsResponse, err := activities.PluginFetchNextAccounts(
			infiniteRetryContext(ctx),
			plugin,
			fetchNextAccount.FromPayload,
			state.State,
			fetchNextAccount.Config.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next accounts")
		}

		err = activities.StorageStoreAccounts(
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
		err = activities.StorageStoreState(
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

			if err := workflow.ExecuteChildWorkflow(
				workflow.WithChildOptions(
					ctx,
					workflow.ChildWorkflowOptions{
						TaskQueue:         fetchNextAccount.ConnectorID.Reference,
						ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
					},
				),
				Run,
				plugin,
				fetchNextAccount.Config,
				fetchNextAccount.ConnectorID,
				payload,
				nextTasks,
			).Get(ctx, nil); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = accountsResponse.HasMore
	}

	return nil
}

var RunFetchNextAccounts any

func init() {
	RunFetchNextAccounts = Workflow{}.runFetchNextAccounts
}
