package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageStoreAccounts(ctx context.Context, accounts []models.Account) error {
	return a.storage.UpsertAccounts(ctx, accounts)
}

var StorageStoreAccountsActivity = Activities{}.StorageStoreAccounts

func StorageStoreAccounts(ctx workflow.Context, accounts []models.Account) error {
	return executeActivity(ctx, StorageStoreAccountsActivity, nil, accounts)
}
