package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageDeleteSchedules(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.DeleteSchedulesFromConnectorID(ctx, connectorID)
}

var StorageDeleteSchedulesActivity = Activities{}.StorageDeleteSchedules

func StorageDeleteSchedules(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageDeleteSchedulesActivity, nil, connectorID)
}
