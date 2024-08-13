package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageDeleteAccounts(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.DeleteAccountsFromConnectorID(ctx, connectorID)
}

var StorageDeleteAccountsActivity = Activities{}.StorageDeleteAccounts

func StorageDeleteAccounts(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageDeleteAccountsActivity, nil, connectorID)
}
