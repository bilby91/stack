package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StoreAccounts(ctx context.Context, accounts []models.Account) error {
	return a.storage.UpsertAccounts(ctx, accounts)
}

var StoreAccountsActivity = Activities{}.StoreAccounts

func StoreAccounts(ctx workflow.Context, accounts []models.Account) error {
	return executeActivity(ctx, StoreAccountsActivity, nil, accounts)
}
