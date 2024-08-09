package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorePayments(ctx context.Context, payments []models.Payment) error {
	return a.storage.UpsertPayment(ctx, payments)
}

var StorePaymentsActivity = Activities{}.StorePayments

func StorePayments(ctx workflow.Context, payments []models.Payment) error {
	return executeActivity(ctx, StorePaymentsActivity, nil, payments)
}
