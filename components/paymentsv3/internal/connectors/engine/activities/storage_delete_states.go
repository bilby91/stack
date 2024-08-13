package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageDeleteStates(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.DeleteStatesFromConnectorID(ctx, connectorID)
}

var StorageDeleteStatesActivity = Activities{}.StorageDeleteStates

func StorageDeleteStates(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageDeleteStatesActivity, nil, connectorID)
}
