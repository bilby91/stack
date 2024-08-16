package activities

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) StorageBankAccountsDeleteRelatedAccounts(ctx context.Context, connectorID models.ConnectorID) error {
	return a.storage.BankAccountsDeleteRelatedAccountFromConnectorID(ctx, connectorID)
}

var StorageBankAccountsDeleteRelatedAccountsActivity = Activities{}.StorageBankAccountsDeleteRelatedAccounts

func StorageBankAccountsDeleteRelatedAccounts(ctx workflow.Context, connectorID models.ConnectorID) error {
	return executeActivity(ctx, StorageBankAccountsDeleteRelatedAccountsActivity, nil, connectorID)
}
