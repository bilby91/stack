package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageDeleteTasksTree(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.DeleteTasksFromConnectorID(ctx, connectorID)
}

var StorageDeleteTasksTreeActivity = Activities{}.StorageDeleteTasksTree

func StorageDeleteTasksTree(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageDeleteTasksTreeActivity, nil, connectorID)
}
