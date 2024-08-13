package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageDeleteWorkflow(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.DeleteWorkflowsFromConnectorID(ctx, connectorID)
}

var StorageDeleteWorkflowActivity = Activities{}.StorageDeleteWorkflow

func StorageDeleteWorkflows(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageDeleteWorkflowActivity, nil, connectorID)
}
