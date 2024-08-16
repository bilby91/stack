package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageInstancesDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.InstancesDeleteFromConnectorID(ctx, connectorID)
}

var StorageInstancesDeleteActivity = Activities{}.StorageInstancesDelete

func StorageInstancesDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageInstancesDeleteActivity, nil, connectorID)
}
