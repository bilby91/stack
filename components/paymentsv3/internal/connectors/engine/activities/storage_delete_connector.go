package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageDeleteConnector(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.UninstallConnector(ctx, connectorID)
}

var StorageDeleteConnectorActivity = Activities{}.StorageDeleteConnector

func StorageDeleteConnector(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageDeleteConnectorActivity, nil, connectorID)
}
