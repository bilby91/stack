package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageDeleteInstances(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.DeleteInstancesFromConnectorID(ctx, connectorID)
}

var StorageDeleteInstancesActivity = Activities{}.StorageDeleteInstances

func StorageDeleteInstances(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageDeleteInstancesActivity, nil, connectorID)
}
