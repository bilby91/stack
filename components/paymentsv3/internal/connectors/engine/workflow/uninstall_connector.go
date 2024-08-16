package workflow

import (
	"github.com/formancehq/paymentsv3/internal/connectors/engine/activities"
	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type UninstallConnector struct {
	ConnectorID models.ConnectorID
}

func (w Workflow) runUninstallConnector(
	ctx workflow.Context,
	uninstallConnector UninstallConnector,
) error {
	if err := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(
			ctx,
			workflow.ChildWorkflowOptions{
				TaskQueue:         uninstallConnector.ConnectorID.Reference,
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			},
		),
		RunTerminateSchedules,
		TerminateSchedules{
			ConnectorID: uninstallConnector.ConnectorID,
		},
	).Get(ctx, nil); err != nil {
		return errors.Wrap(err, "terminate schedules")
	}

	// TODO(polo): workflow.Go
	err := activities.StorageSchedulesDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	err = activities.StorageInstancesDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	err = activities.StorageTasksTreeDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	err = activities.StorageWorkflowsDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	err = activities.StorageBankAccountsDeleteRelatedAccounts(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	err = activities.StorageAccountsDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	err = activities.StoragePaymentsDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	err = activities.StorageStatesDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	err = activities.StorageConnectorsDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	return nil
}

var RunUninstallConnector any

func init() {
	RunUninstallConnector = Workflow{}.runUninstallConnector
}
