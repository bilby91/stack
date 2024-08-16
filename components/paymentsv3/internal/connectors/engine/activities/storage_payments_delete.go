package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StoragePaymentsDelete(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.PaymentsDeleteFromConnectorID(ctx, connectorID)
}

var StoragePaymentsDeleteActivity = Activities{}.StoragePaymentsDelete

func StoragePaymentsDelete(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StoragePaymentsDeleteActivity, nil, connectorID)
}
