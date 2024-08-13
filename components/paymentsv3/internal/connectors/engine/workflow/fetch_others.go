package workflow

import (
	"encoding/json"

	"github.com/formancehq/paymentsv3/internal/connectors/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type FetchNextOthers struct {
	Config      models.Config      `json:"config"`
	ConnectorID models.ConnectorID `json:"connectorID"`
	Name        string             `json:"name"`
	FromPayload json.RawMessage    `json:"fromPayload"`
}

func (w Workflow) runFetchNextOthers(
	ctx workflow.Context,
	plugin models.Plugin,
	fetchNextOthers FetchNextOthers,
	nextTasks []models.TaskTree,
) error {
	stateID := models.StateID{
		Reference:   workflow.GetInfo(ctx).WorkflowExecution.ID,
		ConnectorID: fetchNextOthers.ConnectorID,
	}
	state, err := activities.StorageFetchState(infiniteRetryContext(ctx), stateID)
	if err != nil {
		return errors.Wrapf(err, "retrieving state: %s", stateID.String)
	}

	hasMore := true
	for hasMore {
		othersResponse, err := activities.PluginFetchNextOthers(
			infiniteRetryContext(ctx),
			plugin,
			fetchNextOthers.Name,
			fetchNextOthers.FromPayload,
			state.State,
			fetchNextOthers.Config.PageSize,
		)
		if err != nil {
			return errors.Wrap(err, "fetching next others")
		}

		state.State = othersResponse.NewState
		err = activities.StorageStoreState(
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

			if err := workflow.ExecuteChildWorkflow(
				workflow.WithChildOptions(
					ctx,
					workflow.ChildWorkflowOptions{
						TaskQueue:         fetchNextOthers.ConnectorID.Reference,
						ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
					},
				),
				Run,
				plugin,
				fetchNextOthers.Config,
				fetchNextOthers.ConnectorID,
				payload,
				nextTasks,
			).Get(ctx, nil); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}

		hasMore = othersResponse.HasMore
	}

	return nil
}

var RunFetchNextOthers any

func init() {
	RunFetchNextOthers = Workflow{}.runFetchNextOthers
}
