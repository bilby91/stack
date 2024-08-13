package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStorePayments(ctx context.Context, payments []models.Payment) error {
	return a.storage.UpsertPayment(ctx, payments)
}

var StorageStorePaymentsActivity = Activities{}.StorageStorePayments

func StorageStorePayments(ctx workflow.Context, payments []models.Payment) error {
	return executeActivity(ctx, StorageStorePaymentsActivity, nil, payments)
}
