package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageDeletePayments(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.DeletePaymentsFromConnectorID(ctx, connectorID)
}

var StorageDeletePaymentsActivity = Activities{}.StorageDeletePayments

func StorageDeletePayments(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageDeletePaymentsActivity, nil, connectorID)
}
