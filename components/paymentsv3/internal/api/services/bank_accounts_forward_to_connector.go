package services

import (
	"context"

	"github.com/formancehq/paymentsv3/internal/models"
	"github.com/google/uuid"
)

func (s *Service) BankAccountsForwardToConnector(ctx context.Context, bankAccountID uuid.UUID, connectorID models.ConnectorID) (*models.BankAccount, error) {
	return s.engine.CreateBankAccount(ctx, bankAccountID, connectorID)
}
