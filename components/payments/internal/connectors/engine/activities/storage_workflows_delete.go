package activities

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageWorkflowsDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.WorkflowsDeleteFromConnectorID(ctx, connectorID)
}

var StorageWorkflowsDeleteActivity = Activities{}.StorageWorkflowsDelete

func StorageWorkflowsDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageWorkflowsDeleteActivity, nil, connectorID)
}
