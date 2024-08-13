package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStoreConnector(ctx context.Context, connector models.Connector) error {
	return a.storage.InstallConnector(ctx, connector)
}

var StorageStoreConnectorActivity = Activities{}.StorageStoreConnector

func StorageStoreConnector(ctx workflow.Context, connector models.Connector) error {
	return executeActivity(ctx, StorageStoreConnectorActivity, nil, connector)
}
